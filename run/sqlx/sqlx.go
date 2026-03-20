package sqlx

import (
	"context"
	"fmt"

	jsqlx "github.com/jmoiron/sqlx"
	"github.com/routis819/sqb"
)

func NamedQueryStruct[T any](db *jsqlx.DB, s sqb.Statement[T], arg any) ([]*T, error) {
	rows, err := db.NamedQuery(s.String(), arg)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	return scanRows[T](rows)
}

func NamedQueryStructContext[T any](ctx context.Context, db *jsqlx.DB, s sqb.Statement[T], arg any) ([]*T, error) {
	rows, err := db.NamedQueryContext(ctx, s.String(), arg)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	return scanRows[T](rows)
}

func scanRows[T any](rows *jsqlx.Rows) ([]*T, error) {
	defer rows.Close()

	r := make([]*T, 0, 32)
	for rows.Next() {
		var row T
		if err := rows.StructScan(&row); err != nil {
			return nil, fmt.Errorf("failed to StructScan: %w", err)
		}
		r = append(r, &row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return r, nil
}
