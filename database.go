package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"errors"
	"io"

	blocks "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-block-format"
	cid "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-cid"
	ipldcbor "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipld-cbor"
	format "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-ipld-format"
	ci "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-crypto"
	peer "github.com/ipsn/go-ipfs/gxlibs/github.com/libp2p/go-libp2p-peer"
	mh "github.com/ipsn/go-ipfs/gxlibs/github.com/multiformats/go-multihash"
	cbor "github.com/ipsn/go-ipfs/gxlibs/github.com/whyrusleeping/cbor/go"

	"github.com/ipsn/go-ipfs/core"
	"github.com/ipsn/go-ipfs/core/coreunix"
	"github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-merkledag"
	"github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-path"

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

func (db *database) GetData(hash string) (uint64, []byte, error) {
	c, err := cid.Decode(hash)
	if err != nil {
		return 0, nil, err
	}
	block, err := db.node.Blocks.GetBlock(context.TODO(), c)
	if err != nil {
		return 0, nil, err
	}
	return c.Type(), block.RawData(), nil
}

func (db *database) PutData(typ uint64, data []byte) (string, error) {
	hash, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		return "", err
	}
	c := cid.NewCidV1(typ, hash)
	b, err := blocks.NewBlockWithCid(data, c)
	if err != nil {
		return "", err
	}
	err = db.node.Blocks.AddBlock(b)
	if err != nil {
		return "", err
	}
	return c.String(), nil
}

// Get fetches an object from IPFS and decodes it.
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
		dec := cbor.NewDecoder(bytes.NewReader(buf))
		dec.TagDecoders[ipldcbor.CBORTagLink] = new(ipldcbor.IpldLinkDecoder)
		var res interface{}
		err = dec.Decode(&res)
		if err != nil {
			return nil, err
		}
		return res, nil
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

// Put encodes an object and stores it in IPFS.
func (db *database) Put(val interface{}) (string, error) {
	var nd format.Node
	var err error
	switch val := val.(type) {
	case io.Reader:
		return coreunix.AddWithContext(context.TODO(), db.node, val)
	case []byte:
		nd = merkledag.NewRawNode(val)
	default:
		nd, err = ipldcbor.WrapObject(val, mh.SHA2_256, mh.DefaultLengths[mh.SHA2_256])
		if err != nil {
			return "", err
		}
	}
	err = db.node.DAG.Add(context.TODO(), nd)
	if err != nil {
		return "", err
	}
	return nd.Cid().String(), nil
}

// Resolve resolves an address using IPNS.
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

// Publish publishes an address to IPNS.
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

// KeyGen generates a new key for use with IPNS.
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

// KeyList returns all the generated keys.
func (db *database) KeyList() ([]string, error) {
	return db.node.Repo.Keystore().List()
}
