package transaction

import (
	"context"
	"errors"
	"strings"
)

type Signer struct {
	ctx  context.Context
	db   Database
	addr string
	tx   *Transaction
}

func NewSigner(ctx context.Context,
	db Database, addr string) (*Signer, error) {
	s := &Signer{
		ctx:  ctx,
		db:   db,
		addr: addr,
	}
	if addr == "" {
		s.tx = New(s.ctx, s.db, "")
	} else {
		hash, err := s.db.Resolve(s.addr)
		if err != nil {
			return nil, err
		}
		if !strings.HasPrefix(hash, "/ipfs/") {
			return nil, errors.New("Invalid type")
		}
		s.tx = New(s.ctx, s.db, hash[6:])
	}
	return s, nil
}

func (s *Signer) Root() *Object {
	return s.tx.Root()
}

func (s *Signer) Hash() string {
	return s.tx.Hash()
}

func (s *Signer) Commit(addr string) error {
	return s.db.Publish(addr, s.tx.Hash())
}

func (s *Signer) Address() string {
	return s.addr
}
