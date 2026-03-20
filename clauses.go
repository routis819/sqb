package sqb

import (
	"fmt"
	"strings"
)

// SelectClause represents a SELECT statement part.
type selectClause struct {
	columns []string
}

func (s *selectClause) StmtOpen() {}

func (s *selectClause) AcceptNext(c Clause) bool {
	_, ok := c.(*fromClause)
	return ok
}

func (s *selectClause) String() string {
	if len(s.columns) == 0 {
		return "SELECT *"
	}
	return "SELECT " + strings.Join(s.columns, ", ")
}

func Select(columns ...string) Clause {
	return &selectClause{columns: columns}
}

// FromClause represents a FROM statement part.
type fromClause struct {
	table string
}

func (f *fromClause) AcceptNext(c Clause) bool {
	switch c.(type) {
	case *whereClause, *orderByClause, *limitClause:
		return true
	default:
		return false
	}
}

func (f *fromClause) String() string {
	return " FROM " + f.table
}

func (f *fromClause) StmtClose() {}

func From(table string) Clause {
	return &fromClause{table: table}
}

// WhereClause represents a WHERE statement part.
type whereClause struct {
	condition string
}

func (w *whereClause) StmtClose() {}

func (w *whereClause) AcceptNext(c Clause) bool {
	switch c.(type) {
	case *orderByClause, *limitClause:
		return true
	default:
		return false
	}
}

func (w *whereClause) String() string {
	return " WHERE " + w.condition
}

func Where(condition string) Clause {
	return &whereClause{condition: condition}
}

// OrderByClause represents an ORDER BY statement part.
type orderByClause struct {
	columns []string
}

func (o *orderByClause) StmtClose() {}

func (o *orderByClause) AcceptNext(c Clause) bool {
	_, ok := c.(*limitClause)
	return ok
}

func (o *orderByClause) String() string {
	return " ORDER BY " + strings.Join(o.columns, ", ")
}

func OrderBy(columns ...string) Clause {
	return &orderByClause{columns: columns}
}

// LimitClause represents a LIMIT statement part.
type limitClause struct {
	limit int
}

func (l *limitClause) StmtClose() {}

func (l *limitClause) AcceptNext(c Clause) bool {
	return false
}

func (l *limitClause) String() string {
	return " LIMIT " + fmt.Sprint(l.limit)
}

func Limit(limit int) Clause {
	return &limitClause{limit: limit}
}
