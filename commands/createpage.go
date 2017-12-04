package commands

import (
	"context"
	"fmt"

	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
)

func CreatePage(us *userstore.Userstore, db transaction.Database, t, name string) (string, error) {
	si, err := transaction.NewSigner(context.TODO(), db, "")
	if err != nil {
		return "", err
	}
	addr, err := db.KeyGen()
	if err != nil {
		return "", err
	}
	err = si.Root().Set(map[interface{}]interface{}{
		"t": t,
		"info": map[string]interface{}{
			"name": name,
		},
		"writers": map[string]interface{}{
			addr: true,
		},
	})
	if err != nil {
		return "", err
	}
	err = us.UpdateUser(addr, func(obj *pathing.Object) error {
		obj.Get("hash").Set(si.Hash())
		obj.Get("subscriptions").Get(addr).Set(true)
		return nil
	})
	if err != nil {
		return "", err
	}
	go func() {
		err := si.Commit(addr)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return addr, nil
}
