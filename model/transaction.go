package model

type TransactionStatus int

const (
	PENDING TransactionStatus = iota
	COMMIT
	ROLLBACK
)

type Transaction struct {
	model  *TransactionalModel
	status TransactionStatus
}

func NewTransaction(model Model) *Transaction {
	return &Transaction{
		model: NewTransactionalModel(model),
	}
}

func (transaction *Transaction) Model() Model {
	return transaction.model
}

func (transaction *Transaction) Status() TransactionStatus {
	return transaction.status
}

func (transaction *Transaction) Commit() error {
	if err := transaction.model.Commit(); err != nil {
		return err
	}

	transaction.status = COMMIT

	return nil
}

func (transaction *Transaction) Rollback() error {
	if err := transaction.model.Rollback(); err != nil {
		return err
	}

	transaction.status = ROLLBACK

	return nil
}

func (transaction *Transaction) Close() error {
	if transaction.status != COMMIT {
		return transaction.Rollback()
	}
	return nil
}

type TransactionProvider interface {
	NewTransaction() *Transaction
}
