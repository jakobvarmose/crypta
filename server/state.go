package server

import (
	"github.com/jakobvarmose/crypta/server/rpc"
	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
)

type State struct {
	MyAddress string
	RPC       *rpc.RPC

	DB transaction.Database
	US *userstore.Userstore
}
