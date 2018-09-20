package commands

import (
	"context"

	"github.com/ipsn/go-ipfs/core"
	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
)

type Provider interface {
	N() *core.IpfsNode
	US() *userstore.Userstore
	DB() transaction.Database
	CTX() context.Context
}
