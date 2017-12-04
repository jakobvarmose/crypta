package transaction

import (
	"errors"

	cid "gx/ipfs/QmNp85zy9RLrQ5oQD4hPyS39ezrrXpcaa7R4Y9kxdWQLLQ/go-cid"

	"github.com/jakobvarmose/crypta/pathing"
)

type Object struct {
	tx     *Transaction
	parent *Object
	key    interface{}
	cache  interface{}
	cid    *cid.Cid
}

func (o *Object) Type() pathing.Type {
	switch o.getValue().(type) {
	case map[interface{}]interface{}:
		return pathing.Type_Map
	case []interface{}:
		return pathing.Type_Array
	case string:
		return pathing.Type_String
	case []byte:
		return pathing.Type_Bytes
	case uint64:
		return pathing.Type_Uint64
	case int64:
		return pathing.Type_Int64
	case bool:
		return pathing.Type_Bool
	case nil:
		return pathing.Type_Nil
	default:
		return pathing.Type_Unknown
	}
}

func (o *Object) Get(key interface{}) *Object {
	child := &Object{
		tx:     o.tx,
		parent: o,
		key:    key,
	}
	return child
}

func (o *Object) Append() *Object {
	return o.Get(-1)
}

func (o *Object) Class() *Object {
	return o.Get("")
}

func (o *Object) Each(callback func(interface{}, *Object) error) error {
	switch value := o.getValue().(type) {
	case map[interface{}]interface{}:
		for key := range value {
			val := &Object{
				tx:     o.tx,
				parent: o,
				key:    key,
			}
			err := callback(key, val)
			if err != nil {
				return err
			}
		}
	case []interface{}:
		for key := range value {
			val := &Object{
				tx:     o.tx,
				parent: o,
				key:    key,
			}
			err := callback(key, val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Object) EachSimple(callback func(*Object, *Object) error) error {
	switch value := o.getValue().(type) {
	case map[interface{}]interface{}:
		for key := range value {
			val := &Object{
				tx:     o.tx,
				parent: o,
				key:    key,
			}
			err := callback(&Object{cache: key}, val)
			if err != nil {
				return err
			}
		}
	case []interface{}:
		for key := range value {
			val := &Object{
				tx:     o.tx,
				parent: o,
				key:    key,
			}
			err := callback(&Object{cache: int64(key)}, val)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *Object) Set(value interface{}) error {
	o.cache = nil
	if o.parent == nil {
		hash, err := o.tx.db.Put(value)
		if err != nil {
			return err
		}
		o.tx.hash = hash
		return nil
	}
	var ins interface{}
	data, err := pathing.Marshal(value)
	if err != nil {
		return err
	}
	if len(data) > 256 {
		hash, err := o.tx.db.Put(value)
		if err != nil {
			return err
		}
		c, err := cid.Decode(hash)
		if err != nil {
			return err
		}
		ins = *c
	} else {
		ins = value
	}
	parent := o.parent.getValue()
	switch key := o.key.(type) {
	case string:
		switch parent2 := parent.(type) {
		case map[interface{}]interface{}:
			parent2[o.key] = ins
		default:
			parent = map[interface{}]interface{}{
				o.key: ins,
			}
		}
	case int:
		switch parent2 := parent.(type) {
		case []interface{}:
			if key >= 0 && key < len(parent2) {
				parent2[key] = ins
			} else if key == -1 {
				o.key = len(parent2)
				parent = append(parent2, ins)
			} else {
				return errors.New("invalid index")
			}
		default:
			if key == -1 {
				o.key = 0
				parent = []interface{}{ins}
			} else {
				return errors.New("invalid index")
			}
		}
	default:
		return errors.New("invalid key")
	}
	return o.parent.Set(parent)
}

func (o *Object) Delete() error {
	o.cache = nil
	if o.parent == nil {
		return o.Set(nil)
	}
	parent := o.parent.getValue()
	switch parent2 := parent.(type) {
	case map[interface{}]interface{}:
		delete(parent2, o.key)
	case []interface{}:
		index, ok := o.key.(int)
		if ok && index >= 0 && index < len(parent2) {
			parent = append(parent2[:index], parent2[index+1:]...)
		}
	default:
		return nil
	}
	return o.parent.Set(parent)
}

func (o *Object) String() string {
	switch value := o.getValue().(type) {
	case string:
		return value
	case []byte:
		return string(value)
	default:
		return ""
	}
}

func (o *Object) Bytes() []byte {
	switch value := o.getValue().(type) {
	case []byte:
		return value
	case string:
		return []byte(value)
	default:
		return nil
	}
}

func (o *Object) Int() int {
	switch value := o.getValue().(type) {
	case uint64:
		return int(value)
	case int64:
		return int(value)
	default:
		return 0
	}
}

func (o *Object) Uint64() uint64 {
	switch value := o.getValue().(type) {
	case uint64:
		return value
	case int64:
		return uint64(value)
	default:
		return 0
	}
}

func (o *Object) Int64() int64 {
	switch value := o.getValue().(type) {
	case uint64:
		return int64(value)
	case int64:
		return value
	default:
		return 0
	}
}

func (o *Object) Bool() bool {
	switch value := o.getValue().(type) {
	case bool:
		return value
	default:
		return false
	}
}

func (o *Object) Cid() *cid.Cid {
	o.getValue()
	return o.cid
}

func (o *Object) Len() int {
	switch value := o.getValue().(type) {
	case []interface{}:
		return len(value)
	case map[interface{}]interface{}:
		return len(value)
	default:
		return 0
	}
}

func (o *Object) getValue() interface{} {
	if o.cache != nil {
		return o.cache
	}
	if o.parent == nil {
		item, err := o.tx.db.Get(o.tx.hash)
		if err != nil {
			return nil
		}
		o.cache = item
		return o.cache
	}
	var value interface{}
	payload := o.parent.getValue()
	switch payload := payload.(type) {
	case map[interface{}]interface{}:
		value = payload[o.key]
	case []interface{}:
		index, ok := o.key.(int)
		if ok && index >= 0 && index < len(payload) {
			value = payload[index]
		}
	}
	if value, ok := value.(cid.Cid); ok {
		item, err := o.tx.db.Get(value.String())
		if err != nil {
			return nil
		}
		o.cid = &value
		o.cache = item
		return item
	}
	return value
}
