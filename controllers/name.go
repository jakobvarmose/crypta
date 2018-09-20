package controllers

import (
	"github.com/jakobvarmose/crypta/pathing"
)

func NameResolve(p Provider, args *pathing.Object) (interface{}, error) {
	hash, err := p.DB().Resolve(args.Get("address").String())
	if err != nil {
		return nil, err
	}
	if hash[:6] == "/ipfs/" {
		hash = hash[6:]
	}
	return map[string]interface{}{
		"hash": hash,
	}, nil
}
