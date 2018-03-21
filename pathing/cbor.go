package pathing

import (
	"reflect"

	cid "github.com/ipfs/go-cid"

	multibase "github.com/multiformats/go-multibase"
	"github.com/ugorji/go/codec"
)

var handle = new(codec.CborHandle)

type extLink struct{}

func (extLink) ConvertExt(v interface{}) interface{} {
	var c cid.Cid
	switch v := v.(type) {
	case cid.Cid:
		c = v
	case *cid.Cid:
		c = *v
	default:
		panic("cannot happen")
	}
	buf, err := c.StringOfBase(multibase.Identity)
	if err != nil {
		panic(err)
	}
	return []byte(buf)
}

func (extLink) UpdateExt(dst interface{}, src interface{}) {
	c, _ := cid.Decode(string(src.([]byte)))
	*dst.(*cid.Cid) = *c
}

func init() {
	handle.Canonical = true
	handle.SetInterfaceExt(reflect.TypeOf(cid.Cid{}), 42, new(extLink))
}

func Unmarshal(data []byte) (interface{}, error) {
	d := codec.NewDecoderBytes(data, handle)
	var val interface{}
	if err := d.Decode(&val); err != nil {
		return nil, err
	}
	return val, nil
}

func MustMarshal(val interface{}) []byte {
	data, err := Marshal(val)
	if err != nil {
		panic(err)
	}
	return data
}

func Marshal(val interface{}) ([]byte, error) {
	var data []byte
	e := codec.NewEncoderBytes(&data, handle)
	err := e.Encode(val)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func UnmarshalInto(dst interface{}, src []byte) error {
	d := codec.NewDecoderBytes(src, handle)
	return d.Decode(dst)
}
