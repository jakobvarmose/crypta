package main

import (
	"context"
	"crypto/rand"
	"errors"
	"io"

	cid "github.com/ipfs/go-cid"
	ipldcbor "github.com/ipfs/go-ipld-cbor"
	format "github.com/ipfs/go-ipld-format"
	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	mh "github.com/multiformats/go-multihash"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreunix"
	"github.com/ipfs/go-ipfs/merkledag"
	"github.com/ipfs/go-ipfs/path"

	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/userstore"
)

type database struct {
	us   *userstore.Userstore
	node *core.IpfsNode
}

func NewDatabase(us *userstore.Userstore, n *core.IpfsNode) *database {
	return &database{
		us:   us,
		node: n,
	}
}

func (db *database) Get(hash string) (interface{}, error) {
	c, err := cid.Decode(hash)
	if err != nil {
		return nil, err
	}
	block, err := db.node.Blocks.GetBlock(context.TODO(), c)
	if err != nil {
		return nil, err
	}
	buf := block.RawData()
	switch c.Type() {
	case cid.DagCBOR:
		return pathing.Unmarshal(buf)
	case cid.Raw:
		return buf, nil
	case cid.DagProtobuf:
		r, err := coreunix.Cat(context.TODO(), db.node, "/ipfs/"+hash)
		if err != nil {
			return nil, err
		}
		return io.Reader(r), nil
	}
	return nil, errors.New("unsupported codec")

}

func (db *database) Put(val interface{}) (string, error) {
	var nd format.Node
	switch val := val.(type) {
	case io.Reader:
		return coreunix.AddWithContext(context.TODO(), db.node, val)
	case []byte:
		nd = merkledag.NewRawNode(val)
	default:
		buf, err := pathing.Marshal(val)
		if err != nil {
			return "", err
		}
		nd, err = ipldcbor.Decode(buf, mh.SHA2_256, mh.DefaultLengths[mh.SHA2_256])
		//nd, err := ipldcbor.WrapObject(val, mh.SHA2_256, mh.DefaultLengths[mh.SHA2_256])
		if err != nil {
			return "", err
		}
	}
	err := db.node.DAG.Add(context.TODO(), nd)
	if err != nil {
		return "", err
	}
	return nd.Cid().String(), nil
}

func (db *database) Resolve(addr string) (string, error) {
	obj, err := db.us.GetUser(addr)
	if err != nil {
		return "", err
	}
	hash := obj.Get("hash").String()
	if hash != "" {
		return "/ipfs/" + hash, nil
	}
	path, err := db.node.Namesys.Resolve(context.TODO(), "/ipns/"+addr)
	if err != nil {
		return "", err
	}
	return path.String(), nil
}

func (db *database) Publish(addr, pathString string) error {
	err := db.us.UpdateUser(addr, func(obj *pathing.Object) error {
		obj.Get("hash").Set(pathString)
		return nil
	})
	if err != nil {
		return err
	}
	priv, err := db.node.GetKey(addr)
	if err != nil {
		return err
	}
	return db.node.Namesys.Publish(context.TODO(), priv, path.FromString("/ipfs/"+pathString))
}

func (db *database) KeyGen() (string, error) {
	priv, pub, err := ci.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return "", err
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return "", err
	}
	idString := id.Pretty()
	err = db.node.Repo.Keystore().Put(idString, priv)
	if err != nil {
		return "", err
	}
	return idString, nil
}

func (db *database) KeyList() ([]string, error) {
	return db.node.Repo.Keystore().List()
}
