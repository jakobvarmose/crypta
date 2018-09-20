package controllers

import (
	"fmt"

	"github.com/jakobvarmose/crypta/commands"
	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/transaction"
)

func CanPost(p Provider, args *pathing.Object) (interface{}, error) {
	myAddr := args.Get("myAddress").String()
	addr := args.Get("address").String()
	value := args.Get("value").Bool()
	si, err := transaction.NewSigner(p.CTX(), p.DB(), myAddr)
	if err != nil {
		return nil, err
	}
	if value {
		err = si.Root().Get("writers").Get(addr).Set(true)
	} else {
		err = si.Root().Get("writers").Get(addr).Delete()
	}
	if err != nil {
		return nil, err
	}
	go func() {
		err := si.Commit(myAddr)
		if err != nil {
			fmt.Println(err)
		}
	}()
	return true, nil
}

func Home(p Provider, args *pathing.Object) (interface{}, error) {
	myAddr := args.Get("myAddress").String()
	return commands.Home(p, myAddr)
}

func Notifications(p Provider, args *pathing.Object) (interface{}, error) {
	myAddr := args.Get("myAddress").String()
	user, err := p.US().GetUser(myAddr)
	if err != nil {
		return nil, err
	}
	res := map[interface{}]interface{}{
		"notifications": user.Get("notifications").Value(),
	}
	return res, nil
}

func Subscribe(p Provider, args *pathing.Object) (interface{}, error) {
	myAddr := args.Get("myAddress").String()
	addr := args.Get("address").String()
	value := args.Get("value").Bool()
	err := p.US().UpdateUser(myAddr, func(obj *pathing.Object) error {
		if value {
			obj.Get("subscriptions").Get(addr).Set(true)
		} else {
			obj.Get("subscriptions").Get(addr).Delete()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return true, err
}
