package commands

import (
	"fmt"

	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/transaction"
)

func CreatePage(p Provider, t, name string) (string, error) {
	si, err := transaction.NewSigner(p.CTX(), p.DB(), "")
	if err != nil {
		return "", err
	}
	addr, err := p.DB().KeyGen()
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
	err = p.US().UpdateUser(addr, func(obj *pathing.Object) error {
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
