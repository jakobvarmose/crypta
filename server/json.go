package server

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"

	cid "github.com/ipsn/go-ipfs/gxlibs/github.com/ipfs/go-cid"
)

func convert(val interface{}) interface{} {
	switch val := val.(type) {
	case map[interface{}]interface{}:
		res := make(map[string]interface{}, len(val))
		for k := range val {
			if k, ok := k.(string); ok {
				res[k] = convert(val[k])
			}
		}
		return res
	case map[string]interface{}:
		res := make(map[string]interface{}, len(val))
		for k := range val {
			res[k] = convert(val[k])
		}
		return res
	case []interface{}:
		res := make([]interface{}, len(val))
		for i := range val {
			res[i] = convert(val[i])
		}
		return res
	case *cid.Cid:
		return "\x00*" + val.String()
	case []byte:
		return "\x00" + base64.RawURLEncoding.EncodeToString(val)
	case string:
		if strings.HasPrefix(val, "\x00") {
			return "\x00" + val
		}
		return val
	default:
		return val
	}
}

func unconvert(val interface{}) interface{} {
	switch val := val.(type) {
	case map[interface{}]interface{}:
		res := make(map[string]interface{}, len(val))
		for k := range val {
			if k, ok := k.(string); ok {
				res[k] = unconvert(val[k])
			}
		}
		return res
	case map[string]interface{}:
		res := make(map[string]interface{}, len(val))
		for k := range val {
			res[k] = unconvert(val[k])
		}
		return res
	case []interface{}:
		res := make([]interface{}, len(val))
		for i := range val {
			res[i] = unconvert(val[i])
		}
		return res
	case string:
		if strings.HasPrefix(val, "\x00") {
			if strings.HasPrefix(val, "\x00\x00") {
				return val[1:]
			}
			if strings.HasPrefix(val, "\x00*") {
				res, err := cid.Decode(val[2:])
				if err != nil {
					return nil
				}
				return res
			}
			res, err := base64.RawURLEncoding.DecodeString(val[1:])
			if err != nil {
				return nil
			}
			return res
		}
		return val
	default:
		return val
	}
}

func WriteJSON(w io.Writer, val interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(convert(val))
}

func ReadJSON(r io.Reader) (interface{}, error) {
	var val interface{}
	err := json.NewDecoder(r).Decode(&val)
	if err != nil {
		return nil, err
	}
	return unconvert(val), nil
}
