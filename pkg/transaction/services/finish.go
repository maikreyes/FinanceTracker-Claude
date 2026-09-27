package services

import (
	"finance-tracker/pkg/transaction/model/failure"
	"finance-tracker/pkg/transaction/model/register"
	"finance-tracker/pkg/transaction/model/transaction"
)

type receiptFile struct {
	data        []byte
	filename    string
	contentType string
}

// finish sube el recibo (si hay) y crea la página: el tramo común a todos
// los caminos de registro, una vez que tx ya tiene todos sus campos.
func (s *Services) finish(tx transaction.Transaction, categoryName string, receipt *receiptFile, interpreted bool) (register.Result, error) {
	if receipt != nil {
		fileUploadID, err := s.uploader.Upload(receipt.data, receipt.filename, receipt.contentType)
		if err != nil {
			return register.Result{}, failure.ErrWriteFailed{Err: err}
		}
		tx.ReceiptFileIDs = []string{fileUploadID}
	}

	pageID, err := s.writer.CreateExpense(tx)
	if err != nil {
		return register.Result{}, failure.ErrWriteFailed{Err: err}
	}

	return register.Result{PageID: pageID, Transaction: tx, CategoryName: categoryName, Interpreted: interpreted}, nil
}

// resolveCategory resuelve name contra Notion. Si no hay match, arma
// ErrUnrecognizedCategory con la lista de categorías válidas (best effort:
// si listarlas también falla, el error se devuelve igual, sin la lista).
func (s *Services) resolveCategory(name string) (string, error) {
	categoryID, ok, err := s.categories.Resolve(name)
	if err != nil {
		return "", failure.ErrWriteFailed{Err: err}
	}
	if !ok {
		valid, _ := s.categories.ListNames()
		return "", failure.ErrUnrecognizedCategory{Got: name, Valid: valid}
	}
	return categoryID, nil
}
