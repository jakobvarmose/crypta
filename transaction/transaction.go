package transaction

import (
	"context"
)

type Transaction struct {
	ctx  context.Context
	db   Database
	hash string
}

func New(ctx context.Context, db Database, hash string) *Transaction {
	tx := &Transaction{
		ctx:  ctx,
		db:   db,
		hash: hash,
	}
	return tx
}

func (tx *Transaction) Hash() string {
	return tx.hash
}

func (tx *Transaction) Root() *Object {
	root := &Object{
		tx:     tx,
		parent: nil,
		key:    nil,
	}
	return root
}
