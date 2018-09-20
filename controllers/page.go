package controllers

import (
	"github.com/jakobvarmose/crypta/commands"
	"github.com/jakobvarmose/crypta/pathing"
)

func Page(p Provider, args *pathing.Object) (interface{}, error) {
	addr := args.Get("address").String()
	user := args.Get("myAddress").String()
	return commands.PostList(p, addr, user)
}

func PageSet(p Provider, args *pathing.Object) (interface{}, error) {
	myAddr := args.Get("myAddress").String()
	key := args.Get("key").String()
	val := args.Get("val").String()
	user, err := p.US().GetUser(myAddr)
	if err != nil {
		return nil, err
	}
	err = user.Get("info").Get("key").Set(val)
	if err != nil {
		return nil, err
	}
	err = commands.SetInfo(p, myAddr, key, val)
	if err != nil {
		return nil, err
	}
	return true, nil
}

func PageSetWriters(p Provider, args *pathing.Object) (interface{}, error) {
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
	err := commands.SetWriters(p, id, writers)
	if err != nil {
		return nil, err
	}
	err = p.US().UpdateUser(myAddr, func(obj *pathing.Object) error {
		obj.Get("writers").Set(writers2)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return true, nil
}

func PagePost(p Provider, args *pathing.Object) (interface{}, error) {
	addr := args.Get("address").String()
	myAddr := args.Get("myAddress").String()
	text := args.Get("text").String()
	attachments := args.Get("attachments")
	result, err := commands.CreatePost(p, addr, myAddr, text, attachments)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"result": result,
	}, nil
}

func PageComment(p Provider, args *pathing.Object) (interface{}, error) {
	addr := args.Get("address").String()
	myAddr := args.Get("myAddress").String()
	postHash := args.Get("postHash").String()
	text := args.Get("text").String()
	result, err := commands.CreateTextComment(p, addr, myAddr, postHash, text)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"result": result,
	}, nil
}
