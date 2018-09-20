package controllers

import (
	"github.com/jakobvarmose/crypta/pathing"
)

func BlockGet(p Provider, args *pathing.Object) (interface{}, error) {
	typ, data, err := p.DB().GetData(args.Get("hash").String())
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"type": typ,
		"data": data,
	}, nil
}

func BlockPut(p Provider, args *pathing.Object) (interface{}, error) {
	typ := args.Get("type").Uint64()
	data := args.Get("data").Bytes()
	hash, err := p.DB().PutData(typ, data)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"hash": hash,
	}, nil
}
