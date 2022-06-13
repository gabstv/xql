package xql

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

type Optional[T comparable] struct {
	Val   T
	Valid bool
}

// Scan implements the Scanner interface.
func (n *Optional[T]) Scan(value any) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	switch value.(type) {
	case T:
		n.Valid = true
		n.Val = value.(T)
	}
	if vv, ok := any(n.Val).(sql.Scanner); ok {
		if err := vv.Scan(value); err != nil {
			return err
		}
		n.Valid = true
		return nil
	}
	if vv, ok := any(&n.Val).(sql.Scanner); ok {
		if err := vv.Scan(value); err != nil {
			return err
		}
		n.Valid = true
		return nil
	}
	return fmt.Errorf("Optional[T]: unsupported value %v (type %T) converting to %T", value, value, n.Val)
}

// Value implements the driver Valuer interface.
func (v Optional[T]) Value() (driver.Value, error) {
	if v.Valid {
		return v.Val, nil
	}
	return nil, nil
}

func Opt[T comparable](v T) Optional[T] {
	return Optional[T]{
		Val:   v,
		Valid: true,
	}
}

// Optz checks if v is the zero value equivalent of T. If true, the Valid
// field of the Optional is set to false.
func Optz[T comparable](v T) Optional[T] {
	var zv T
	if zv == v {
		return Optional[T]{
			Val:   v,
			Valid: false,
		}
	}
	return Optional[T]{
		Val:   v,
		Valid: true,
	}
}
