package sqb

import (
	"fmt"
	"strings"
)

type Opener interface {
	StmtOpen()
}

type Closer interface {
	StmtClose()
}

type Clause interface {
	AcceptNext(Clause) bool
	String() string
}

type Statement[T any] struct{ stmtString string }

func (s Statement[T]) String() string {
	return s.stmtString
}

type StmtBuilder[T any] struct {
	building []Clause
}

func (sb *StmtBuilder[T]) AcceptNext(c Clause) bool {
	if len(sb.building) != 0 {
		return sb.building[len(sb.building)-1].AcceptNext(c)
	}
	if _, ok := c.(Opener); !ok {
		return false
	}

	return true

}

func (sb *StmtBuilder[T]) Append(c Clause) error {
	if !sb.AcceptNext(c) {
		return fmt.Errorf("failed to append: %v", c)
	}
	sb.building = append(sb.building, c)

	return nil
}

func (sb *StmtBuilder[T]) MustAppend(c Clause) {
	if err := sb.Append(c); err != nil {
		panic(err)
	}
}

func (sb *StmtBuilder[T]) Stmt() (Statement[T], error) {
	if len(sb.building) == 0 {
		return Statement[T]{}, fmt.Errorf("empty builder")
	}
	sbuilder := strings.Builder{}
	for _, clause := range sb.building {
		sbuilder.WriteString(clause.String())
	}
	if _, ok := sb.building[len(sb.building)-1].(Closer); !ok {
		return Statement[T]{}, fmt.Errorf("statement is not finished: %s", sbuilder.String())
	}
	sbuilder.WriteString(";")

	return Statement[T]{stmtString: sbuilder.String()}, nil
}

func (sb *StmtBuilder[T]) MustStmt() Statement[T] {
	stmt, err := sb.Stmt()
	if err != nil {
		panic(err)
	}

	return stmt
}
