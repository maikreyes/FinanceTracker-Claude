package services

import (
	"errors"
	"testing"

	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/receipt"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

type fakeCategoryResolver struct {
	ids   map[string]string
	names []string
	err   error
}

func (f fakeCategoryResolver) Resolve(name string) (string, bool, error) {
	if f.err != nil {
		return "", false, f.err
	}
	id, ok := f.ids[name]
	return id, ok, nil
}

func (f fakeCategoryResolver) ListNames() ([]string, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.names, nil
}

func (f fakeCategoryResolver) Create(name string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if id, ok := f.ids[name]; ok {
		return id, nil
	}
	return "new-category-id", nil
}

type fakeExpenseWriter struct {
	got transaction.Transaction
	id  string
	err error
}

func (f *fakeExpenseWriter) CreateExpense(tx transaction.Transaction) (string, error) {
	f.got = tx
	return f.id, f.err
}

type fakeReceiptUploader struct {
	gotData        []byte
	gotFilename    string
	gotContentType string
	fileUploadID   string
	err            error
}

func (f *fakeReceiptUploader) Upload(data []byte, filename, contentType string) (string, error) {
	f.gotData = data
	f.gotFilename = filename
	f.gotContentType = contentType
	return f.fileUploadID, f.err
}

type fakeReceiptInterpreter struct {
	gotCategoryNames []string
	result           receipt.Interpreted
	err              error
}

func (f *fakeReceiptInterpreter) Interpret(imageData []byte, categoryNames []string) (receipt.Interpreted, error) {
	f.gotCategoryNames = categoryNames
	return f.result, f.err
}

func amountPtr(v float64) *float64 { return &v }

func TestRegistrar_Register_ResolvesCategory(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{"comida": "cat-123"}}
	writer := &fakeExpenseWriter{id: "page-1"}
	r := NewServices(resolver, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})

	result, err := r.Register("egreso mecato 13000 efectivo comida")
	if err != nil {
		t.Fatalf("Register error inesperado: %v", err)
	}
	if result.PageID != "page-1" {
		t.Errorf("PageID = %q, want %q", result.PageID, "page-1")
	}
	if result.Transaction.CategoryID != "cat-123" {
		t.Errorf("CategoryID = %q, want %q", result.Transaction.CategoryID, "cat-123")
	}
	if result.CategoryName != "comida" {
		t.Errorf("CategoryName = %q, want %q", result.CategoryName, "comida")
	}
}

func TestRegistrar_Register_SinCategoria(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{id: "page-2"}
	r := NewServices(resolver, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})

	result, err := r.Register("ingreso 3000000 banco")
	if err != nil {
		t.Fatalf("Register error inesperado: %v", err)
	}
	if result.Transaction.CategoryID != "" {
		t.Errorf("CategoryID = %q, want vacío", result.Transaction.CategoryID)
	}
	if result.CategoryName != "" {
		t.Errorf("CategoryName = %q, want vacío", result.CategoryName)
	}
}

func TestRegistrar_Register_CategoriaDesconocida(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{}, names: []string{"Comida", "Ahorro"}}
	writer := &fakeExpenseWriter{}
	r := NewServices(resolver, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})

	_, err := r.Register("egreso mecato 13000 efectivo inventada")
	var target failure.ErrUnrecognizedCategory
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want failure.ErrUnrecognizedCategory", err)
	}
	if target.Got != "inventada" {
		t.Errorf("Got = %q, want %q", target.Got, "inventada")
	}
	if len(target.Valid) != 2 {
		t.Errorf("Valid = %v, want 2 elementos", target.Valid)
	}
}

func TestRegistrar_Register_FallaEscritura(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{err: errors.New("boom")}
	r := NewServices(resolver, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})

	_, err := r.Register("egreso mecato 13000 efectivo")
	var target failure.ErrWriteFailed
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want failure.ErrWriteFailed", err)
	}
}

func TestRegistrar_RegisterWithPhoto_SubeYAdjuntaFoto(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{"comida": "cat-123"}}
	writer := &fakeExpenseWriter{id: "page-3"}
	uploader := &fakeReceiptUploader{fileUploadID: "file-upload-1"}
	r := NewServices(resolver, writer, uploader, &fakeReceiptInterpreter{})

	photo := []byte("fake-jpeg-bytes")
	result, err := r.RegisterWithPhoto("egreso mecato 13000 efectivo comida", photo, "recibo.jpg", "image/jpeg")
	if err != nil {
		t.Fatalf("RegisterWithPhoto error inesperado: %v", err)
	}
	if result.PageID != "page-3" {
		t.Errorf("PageID = %q, want %q", result.PageID, "page-3")
	}
	if len(result.Transaction.ReceiptFileIDs) != 1 || result.Transaction.ReceiptFileIDs[0] != "file-upload-1" {
		t.Errorf("ReceiptFileIDs = %v, want [file-upload-1]", result.Transaction.ReceiptFileIDs)
	}
	if uploader.gotFilename != "recibo.jpg" || uploader.gotContentType != "image/jpeg" {
		t.Errorf("Upload llamado con filename=%q contentType=%q", uploader.gotFilename, uploader.gotContentType)
	}
	if writer.got.CategoryID != "cat-123" {
		t.Errorf("CategoryID = %q, want %q", writer.got.CategoryID, "cat-123")
	}
}

func TestRegistrar_RegisterWithPhoto_FallaSubida(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{}
	uploader := &fakeReceiptUploader{err: errors.New("boom")}
	r := NewServices(resolver, writer, uploader, &fakeReceiptInterpreter{})

	_, err := r.RegisterWithPhoto("egreso mecato 13000 efectivo", []byte("x"), "recibo.jpg", "image/jpeg")
	var target failure.ErrWriteFailed
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want failure.ErrWriteFailed", err)
	}
}

func TestRegistrar_RegisterWithPhoto_FormatoInvalido(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{}
	uploader := &fakeReceiptUploader{}
	r := NewServices(resolver, writer, uploader, &fakeReceiptInterpreter{})

	_, err := r.RegisterWithPhoto("bitcoin", []byte("x"), "recibo.jpg", "image/jpeg")
	var target failure.ErrInvalidFormat
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want failure.ErrInvalidFormat", err)
	}
	if uploader.gotData != nil {
		t.Error("no debería haber intentado subir la foto si el caption es inválido")
	}
}

func TestRegistrar_RegisterFromPhotoAI_Completo(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{"Comida": "cat-123"}, names: []string{"Comida", "Ahorro"}}
	writer := &fakeExpenseWriter{id: "page-ai-1"}
	uploader := &fakeReceiptUploader{fileUploadID: "file-upload-ai"}
	interpreter := &fakeReceiptInterpreter{result: receipt.Interpreted{
		Amount:        amountPtr(45000),
		Concept:       "Supermercado XYZ",
		CategoryName:  "Comida",
		PaymentMethod: "Credit Card",
	}}
	r := NewServices(resolver, writer, uploader, interpreter)

	result, err := r.RegisterFromPhotoAI([]byte("fake-jpeg"), "recibo.jpg", "image/jpeg", transaction.TypeExpense)
	if err != nil {
		t.Fatalf("RegisterFromPhotoAI error inesperado: %v", err)
	}
	if !result.Interpreted {
		t.Error("Interpreted = false, want true")
	}
	if result.Transaction.Type != transaction.TypeExpense {
		t.Errorf("Type = %q, want Egreso", result.Transaction.Type)
	}
	if result.Transaction.Amount != 45000 {
		t.Errorf("Amount = %v, want 45000", result.Transaction.Amount)
	}
	if result.Transaction.Name != "Supermercado XYZ" {
		t.Errorf("Name = %q, want %q", result.Transaction.Name, "Supermercado XYZ")
	}
	if result.Transaction.PaymentMethod != transaction.PaymentMethodCreditCard {
		t.Errorf("PaymentMethod = %q, want Credit Card", result.Transaction.PaymentMethod)
	}
	if result.Transaction.CategoryID != "cat-123" {
		t.Errorf("CategoryID = %q, want cat-123", result.Transaction.CategoryID)
	}
	if len(interpreter.gotCategoryNames) != 2 {
		t.Errorf("categoryNames pasados a Interpret = %v, want 2 elementos", interpreter.gotCategoryNames)
	}
}

func TestRegistrar_RegisterFromPhotoAI_SinMonto(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{}
	uploader := &fakeReceiptUploader{}
	interpreter := &fakeReceiptInterpreter{result: receipt.Interpreted{Amount: nil}}
	r := NewServices(resolver, writer, uploader, interpreter)

	_, err := r.RegisterFromPhotoAI([]byte("x"), "recibo.jpg", "image/jpeg", transaction.TypeExpense)
	var target failure.ErrReceiptAmountNotFound
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want failure.ErrReceiptAmountNotFound", err)
	}
	if uploader.gotData != nil {
		t.Error("no debería haber intentado subir la foto sin monto")
	}
}

func TestRegistrar_RegisterFromPhotoAI_FallaInterpretacion(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{}
	uploader := &fakeReceiptUploader{}
	interpreter := &fakeReceiptInterpreter{err: errors.New("groq caído")}
	r := NewServices(resolver, writer, uploader, interpreter)

	_, err := r.RegisterFromPhotoAI([]byte("x"), "recibo.jpg", "image/jpeg", transaction.TypeExpense)
	var target failure.ErrReceiptInterpretationFailed
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want failure.ErrReceiptInterpretationFailed", err)
	}
}

func TestRegistrar_RegisterFromPhotoAI_CategoriaYMedioDePagoNoDeterminados(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{}}
	writer := &fakeExpenseWriter{id: "page-ai-2"}
	uploader := &fakeReceiptUploader{}
	interpreter := &fakeReceiptInterpreter{result: receipt.Interpreted{
		Amount: amountPtr(9000),
		// Concept, CategoryName y PaymentMethod quedan vacíos.
	}}
	r := NewServices(resolver, writer, uploader, interpreter)

	result, err := r.RegisterFromPhotoAI([]byte("x"), "recibo.jpg", "image/jpeg", transaction.TypeExpense)
	if err != nil {
		t.Fatalf("RegisterFromPhotoAI error inesperado (categoría/medio de pago vacíos no deberían bloquear): %v", err)
	}
	if result.Transaction.CategoryID != "" {
		t.Errorf("CategoryID = %q, want vacío", result.Transaction.CategoryID)
	}
	if result.Transaction.PaymentMethod != "" {
		t.Errorf("PaymentMethod = %q, want vacío", result.Transaction.PaymentMethod)
	}
	if result.Transaction.Name != string(transaction.TypeExpense) {
		t.Errorf("Name = %q, want %q (default)", result.Transaction.Name, transaction.TypeExpense)
	}
}

func TestRegistrar_RegisterFromPhotoAI_CategoriaSugeridaNoMatchea(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{}} // "Otra" no existe
	writer := &fakeExpenseWriter{id: "page-ai-3"}
	uploader := &fakeReceiptUploader{}
	interpreter := &fakeReceiptInterpreter{result: receipt.Interpreted{
		Amount:       amountPtr(9000),
		CategoryName: "Otra",
	}}
	r := NewServices(resolver, writer, uploader, interpreter)

	result, err := r.RegisterFromPhotoAI([]byte("x"), "recibo.jpg", "image/jpeg", transaction.TypeExpense)
	if err != nil {
		t.Fatalf("una categoría de IA que no matchea no debería bloquear el registro: %v", err)
	}
	if result.Transaction.CategoryID != "" || result.CategoryName != "" {
		t.Errorf("categoría debería quedar vacía, no bloquear: CategoryID=%q CategoryName=%q", result.Transaction.CategoryID, result.CategoryName)
	}
}

func TestRegistrar_RegisterFromFields_Completo(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{"comida": "cat-123"}}
	writer := &fakeExpenseWriter{id: "page-fields-1"}
	r := NewServices(resolver, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})

	result, err := r.RegisterFromFields(register.ConversationInput{
		Type:          transaction.TypeExpense,
		Concept:       "mecato",
		Amount:        13000,
		PaymentMethod: transaction.PaymentMethodCash,
		CategoryName:  "comida",
	})
	if err != nil {
		t.Fatalf("RegisterFromFields error inesperado: %v", err)
	}
	if result.PageID != "page-fields-1" {
		t.Errorf("PageID = %q, want %q", result.PageID, "page-fields-1")
	}
	if result.Transaction.Name != "mecato" {
		t.Errorf("Name = %q, want %q", result.Transaction.Name, "mecato")
	}
	if result.Transaction.CategoryID != "cat-123" {
		t.Errorf("CategoryID = %q, want %q", result.Transaction.CategoryID, "cat-123")
	}
}

func TestRegistrar_RegisterFromFields_SinConceptoUsaDefault(t *testing.T) {
	resolver := fakeCategoryResolver{}
	writer := &fakeExpenseWriter{id: "page-fields-2"}
	r := NewServices(resolver, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})

	result, err := r.RegisterFromFields(register.ConversationInput{
		Type:          transaction.TypeIncome,
		Amount:        3000000,
		PaymentMethod: transaction.PaymentMethodBank,
	})
	if err != nil {
		t.Fatalf("RegisterFromFields error inesperado: %v", err)
	}
	if result.Transaction.Name != string(transaction.TypeIncome) {
		t.Errorf("Name = %q, want %q (default)", result.Transaction.Name, transaction.TypeIncome)
	}
	if result.Transaction.CategoryID != "" {
		t.Errorf("CategoryID = %q, want vacío", result.Transaction.CategoryID)
	}
}

func TestRegistrar_RegisterFromFields_CategoriaDesconocida(t *testing.T) {
	resolver := fakeCategoryResolver{ids: map[string]string{}, names: []string{"Comida", "Ahorro"}}
	writer := &fakeExpenseWriter{}
	r := NewServices(resolver, writer, &fakeReceiptUploader{}, &fakeReceiptInterpreter{})

	_, err := r.RegisterFromFields(register.ConversationInput{
		Type:          transaction.TypeExpense,
		Amount:        13000,
		PaymentMethod: transaction.PaymentMethodCash,
		CategoryName:  "inventada",
	})
	var target failure.ErrUnrecognizedCategory
	if !errors.As(err, &target) {
		t.Fatalf("error = %v, want failure.ErrUnrecognizedCategory", err)
	}
}
