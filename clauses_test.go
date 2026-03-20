package sqb

import (
	"testing"
)

type User struct {
	ID   int
	Name string
}

func TestSelectFromWhere(t *testing.T) {
	builder := &StmtBuilder[User]{}
	builder.MustAppend(Select())
	builder.MustAppend(From("users"))
	builder.MustAppend(Where("id = :id"))

	stmt, err := builder.Stmt()
	if err != nil {
		t.Fatalf("failed to build statement: %v", err)
	}

	expected := "SELECT * FROM users WHERE id = :id;"
	if stmt.String() != expected {
		t.Errorf("expected %q, got %q", expected, stmt.String())
	}
}

func TestFullQuery(t *testing.T) {
	builder := &StmtBuilder[User]{}
	builder.MustAppend(Select("id", "name"))
	builder.MustAppend(From("users"))
	builder.MustAppend(Where("active = 1"))
	builder.MustAppend(OrderBy("name DESC"))
	builder.MustAppend(Limit(10))

	stmt, err := builder.Stmt()
	if err != nil {
		t.Fatalf("failed to build statement: %v", err)
	}

	expected := "SELECT id, name FROM users WHERE active = 1 ORDER BY name DESC LIMIT 10;"
	if stmt.String() != expected {
		t.Errorf("expected %q, got %q", expected, stmt.String())
	}
}

func TestFromWithLimit(t *testing.T) {
	builder := &StmtBuilder[User]{}
	builder.MustAppend(Select())
	builder.MustAppend(From("users"))
	builder.MustAppend(Limit(5))

	stmt, err := builder.Stmt()
	if err != nil {
		t.Fatalf("failed to build statement: %v", err)
	}

	expected := "SELECT * FROM users LIMIT 5;"
	if stmt.String() != expected {
		t.Errorf("expected %q, got %q", expected, stmt.String())
	}
}
