package usecase

import (
	"errors"
	"testing"

	"finance-tracker/internal/domain"
)

type fakeCategoryResolver struct {
	ids map[string]string
	err error
}

func (f fakeCategoryResolver) Resolve(name string) (string, bool, error) {
	if f.err != nil {
		return "", false, f.err
	}
	id, ok := f.ids[name]
	return id, ok, nil
}

type fakeExpenseWriter struct {
	got domain.Transaction
	id  string
	err error
}

func (f *fakeExpenseWriter) CreateExpense(tx domain.Transaction) (string, error) {
	f.got = tx
	return f.id, f.err
}

func TestRegistrar_Register_ResolvesCategory(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{"comida": "cat-123"}}
	writer := &fakeExpenseWriter{id: "page-1"}
	r := NewRegistrar(resolver, writer)

	pageID, err := r.Register("egreso mecato 13000 efectivo comida")
	if err != nil {
		t.Fatalf("Register error inesperado: %v", err)
	}
	if pageID != "page-1" {
		t.Errorf("pageID = %q, want %q", pageID, "page-1")
	}
	if writer.got.CategoryID != "cat-123" {
		t.Errorf("CategoryID = %q, want %q", writer.got.CategoryID, "cat-123")
	}
}

func TestRegistrar_Register_SinCategoria(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{id: "page-2"}
	r := NewRegistrar(resolver, writer)

	_, err := r.Register("ingreso 3000000 banco")
	if err != nil {
		t.Fatalf("Register error inesperado: %v", err)
	}
	if writer.got.CategoryID != "" {
		t.Errorf("CategoryID = %q, want vacío", writer.got.CategoryID)
	}
}

func TestRegistrar_Register_CategoriaDesconocida(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{}}
	writer := &fakeExpenseWriter{}
	r := NewRegistrar(resolver, writer)

	_, err := r.Register("egreso mecato 13000 efectivo inventada")
	var target ErrUnrecognizedCategory
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want ErrUnrecognizedCategory", err)
	}
}

func TestRegistrar_Register_FallaEscritura(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{err: errors.New("boom")}
	r := NewRegistrar(resolver, writer)

	_, err := r.Register("egreso mecato 13000 efectivo")
	var target ErrWriteFailed
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want ErrWriteFailed", err)
	}
}
