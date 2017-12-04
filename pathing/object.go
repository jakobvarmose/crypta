package pathing

import (
	"errors"
)

type transaction2 struct {
	value interface{}
}

type Object struct {
	tx     *transaction2
	parent *Object
	key    interface{}
	cache  interface{}
}

func NewObject(value interface{}) *Object {
	value = convert(value)
	return &Object{
		tx: &transaction2{
			value: value,
		},
	}
}

func (o *Object) Value() interface{} {
	return o.getValue()
}

func (o *Object) Type() Type {
	switch o.getValue().(type) {
	case map[interface{}]interface{}:
		return Type_Map
	case []interface{}:
		return Type_Array
	case string:
		return Type_String
	case []byte:
		return Type_Bytes
	case uint64:
		return Type_Uint64
	case int64:
		return Type_Int64
	case bool:
		return Type_Bool
	case nil:
		return Type_Nil
	default:
		return Type_Unknown
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
		o.tx.value = value
		return nil
	}
	parent := o.parent.getValue()
	switch key := o.key.(type) {
	case string:
		switch parent2 := parent.(type) {
		case map[interface{}]interface{}:
			parent2[o.key] = value
		default:
			parent = map[interface{}]interface{}{
				o.key: value,
			}
		}
	case int:
		switch parent2 := parent.(type) {
		case []interface{}:
			if key >= 0 && key < len(parent2) {
				parent2[key] = value
			} else if key == -1 {
				o.key = len(parent2)
				parent = append(parent2, value)
			} else {
				return errors.New("invalid index")
			}
		default:
			if key == -1 {
				o.key = 0
				parent = []interface{}{value}
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

func convert(obj interface{}) interface{} {
	switch obj := obj.(type) {
	case map[string]interface{}:
		newObj := make(map[interface{}]interface{})
		for key, val := range obj {
			newObj[key] = convert(val)
		}
		return newObj
	case map[interface{}]interface{}:
		for key, val := range obj {
			obj[key] = convert(val)
		}
	case []interface{}:
		for key, val := range obj {
			obj[key] = convert(val)
		}
	}
	return obj
}

func (o *Object) getValue() interface{} {
	if o.cache != nil {
		return o.cache
	}
	if o.parent == nil {
		o.cache = o.tx.value
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
	return value
}
