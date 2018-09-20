package controllers

import (
	"github.com/jakobvarmose/crypta/commands"
	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/transaction"
)

func UserList(p Provider, args *pathing.Object) (interface{}, error) {
	keys, err := p.DB().KeyList()
	if err != nil {
		return nil, err
	}
	val := make([]interface{}, 0)
	for _, key := range keys {
		si, err := transaction.NewSigner(p.CTX(), p.DB(), key)
		if err != nil {
			return nil, err
		}
		name := si.Root().Get("info").Get("name").String()
		val = append(val, map[string]interface{}{
			"address": key,
			"name":    name,
		})
	}
	return val, nil
}

func UserCreate(p Provider, args *pathing.Object) (interface{}, error) {
	name := args.Get("name").String()
	address, err := commands.CreatePage(p, "USER", name)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"address": address,
		"name":    name,
	}, nil
}
