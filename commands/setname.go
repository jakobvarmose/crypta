package commands

import (
	"fmt"

	"github.com/jakobvarmose/crypta/transaction"
)

func SetInfo(p Provider, addr, key, val string) error {
	si, err := transaction.NewSigner(p.CTX(), p.DB(), addr)
	if err != nil {
		return err
	}
	if val == "" {
		err = si.Root().Get("info").Get(key).Delete()
	} else {
		err = si.Root().Get("info").Get(key).Set(val)
	}
	if err != nil {
		return err
	}
	return si.Commit(addr)
}

func GetInfo(p Provider, addr, key string) (string, error) {
	si, err := transaction.NewSigner(p.CTX(), p.DB(), addr)
	if err != nil {
		return "", err
	}
	si.Root().EachSimple(func(arg1 *transaction.Object, arg2 *transaction.Object) error {
		fmt.Println(arg1, arg2)
		return nil
	})
	name := si.Root().Get("info").Get(key).String()
	return name, nil
}

func GetInfos(p Provider, addr string) (map[string]string, error) {
	si, err := transaction.NewSigner(p.CTX(), p.DB(), addr)
	if err != nil {
		return nil, err
	}
	infos := make(map[string]string)
	err = si.Root().Get("info").EachSimple(func(key *transaction.Object, val *transaction.Object) error {
		infos[key.String()] = val.String()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return infos, nil
}
