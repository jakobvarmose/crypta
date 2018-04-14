package server

import (
	"context"
	"fmt"

	"github.com/jakobvarmose/crypta/commands"
	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/server/rpc"
	"github.com/jakobvarmose/crypta/transaction"
)

var methods = rpc.Methods{
	"v0/user/list": func(state *State, args *pathing.Object) (interface{}, error) {
		keys, err := state.DB.KeyList()
		if err != nil {
			return nil, err
		}
		val := make([]interface{}, 0)
		for _, key := range keys {
			si, err := transaction.NewSigner(context.TODO(), state.DB, key)
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
	},
	"v0/user/create": func(state *State, args *pathing.Object) (interface{}, error) {
		name := args.Get("name").String()
		address, err := commands.CreatePage(state.US, state.DB, "USER", name)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"address": address,
			"name":    name,
		}, nil
	},
	"v0/home": func(state *State, args *pathing.Object) (interface{}, error) {
		myAddr := args.Get("myAddress").String()
		return commands.Home(state.US, state.DB, myAddr)
	},
	"v0/subscribe": func(state *State, args *pathing.Object) (interface{}, error) {
		myAddr := args.Get("myAddress").String()
		addr := args.Get("address").String()
		value := args.Get("value").Bool()
		err := state.US.UpdateUser(myAddr, func(obj *pathing.Object) error {
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
	},
	"v0/canPost": func(state *State, args *pathing.Object) (interface{}, error) {
		myAddr := args.Get("myAddress").String()
		addr := args.Get("address").String()
		value := args.Get("value").Bool()
		si, err := transaction.NewSigner(context.TODO(), state.DB, myAddr)
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
	},
	"v0/notifications": func(state *State, args *pathing.Object) (interface{}, error) {
		myAddr := args.Get("myAddress").String()
		user, err := state.US.GetUser(myAddr)
		if err != nil {
			return nil, err
		}
		res := map[interface{}]interface{}{
			"notifications": user.Get("notifications").Value(),
		}
		return res, nil
	},
	"v0/page": func(state *State, args *pathing.Object) (interface{}, error) {
		addr := args.Get("address").String()
		user := args.Get("myAddress").String()
		return commands.PostList(state.US, state.DB, addr, user)
	},
	"v0/page/set": func(state *State, args *pathing.Object) (interface{}, error) {
		myAddr := args.Get("myAddress").String()
		key := args.Get("key").String()
		val := args.Get("val").String()
		user, err := state.US.GetUser(myAddr)
		if err != nil {
			return nil, err
		}
		err = user.Get("info").Get("key").Set(val)
		if err != nil {
			return nil, err
		}
		err = commands.SetInfo(state.DB, myAddr, key, val)
		if err != nil {
			return nil, err
		}
		return true, nil
	},
	"v0/page/setwriters": func(state *State, args *pathing.Object) (interface{}, error) {
		id := args.Get("address").String()
		myAddr := args.Get("myAddress").String()
		var writers []string
		writers2 := make(map[interface{}]interface{})
		args.Get("writers").EachSimple(func(_ *pathing.Object, writer *pathing.Object) error {
			writerString := writer.String()
			writers = append(writers, writerString)
			writers2[writerString] = true
			return nil
		})
		err := commands.SetWriters(state.DB, id, writers)
		if err != nil {
			return nil, err
		}
		err = state.US.UpdateUser(myAddr, func(obj *pathing.Object) error {
			obj.Get("writers").Set(writers2)
			return nil
		})
		if err != nil {
			return nil, err
		}
		return true, nil
	},
	"v0/page/post": func(state *State, args *pathing.Object) (interface{}, error) {
		addr := args.Get("address").String()
		myAddr := args.Get("myAddress").String()
		text := args.Get("text").String()
		attachments := args.Get("attachments")
		result, err := commands.CreatePost(state.DB, addr, myAddr, text, attachments)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"result": result,
		}, nil
	},
	"v0/page/comment": func(state *State, args *pathing.Object) (interface{}, error) {
		addr := args.Get("address").String()
		myAddr := args.Get("myAddress").String()
		postHash := args.Get("postHash").String()
		text := args.Get("text").String()
		result, err := commands.CreateTextComment(state.DB, addr, myAddr, postHash, text)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"result": result,
		}, nil
	},
}
