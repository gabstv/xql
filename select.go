package xql

import (
	"context"
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type SelectArg func(b squirrel.SelectBuilder) squirrel.SelectBuilder

// type Table[T comparable]

type Row interface {
	comparable
	Table() string
}

type RowWithJoin interface {
	TableJoins()
}

type WithSelectColumns interface {
	SelectColumns() []string
}

func Where(pred any, args ...any) SelectArg {
	return func(b squirrel.SelectBuilder) squirrel.SelectBuilder {
		return b.Where(pred, args...)
	}
}

func Select[T Row](db sqlx.Queryer, dst *[]T, args ...SelectArg) error {
	var zv T
	// extract columns
	var cols []string
	if vv, ok := any(zv).(WithSelectColumns); ok {
		cols = vv.SelectColumns()
	} else {
		cols = Explode(",", ExtractStructTags(zv, "columns", "select_column", "column")...)
	}
	b := squirrel.Select(cols...).From(zv.Table())
	for _, arg := range args {
		b = arg(b)
	}
	q, qargs, err := b.ToSql()
	if err != nil {
		return err
	}
	return sqlx.Select(db, dst, q, qargs...)
}

func SelectContext[T Row](ctx context.Context, db sqlx.QueryerContext, dst *[]T, args ...SelectArg) error {
	var zv T
	// extract columns
	b := squirrel.Select(Explode(",", ExtractStructTags(zv, "columns", "select_column", "column")...)...).From(zv.Table())
	for _, arg := range args {
		b = arg(b)
	}
	q, qargs, err := b.ToSql()
	if err != nil {
		return err
	}
	return sqlx.SelectContext(ctx, db, dst, q, qargs...)
}

func Get[T Row](db sqlx.Queryer, dst *T, args ...SelectArg) error {
	var items []T
	if err := Select(db, &items, args...); err != nil {
		return err
	}
	if len(items) == 0 {
		return sql.ErrNoRows
	}
	*dst = items[0]
	return nil
}

func GetContext[T Row](ctx context.Context, db sqlx.QueryerContext, dst *T, args ...SelectArg) error {
	var items []T
	if err := SelectContext(ctx, db, &items, args...); err != nil {
		return err
	}
	if len(items) == 0 {
		return sql.ErrNoRows
	}
	*dst = items[0]
	return nil
}
