package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	cid "github.com/ipfs/go-cid"

	"github.com/ipfs/go-ipfs/core"

	"github.com/jakobvarmose/crypta/commands"
	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
)

func securityCheck(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		originString := req.Header.Get("Origin")
		referrerString := req.Header.Get("Referer")
		if originString == "" && referrerString == "" {
			http.Error(resp, "403 forbidden - no origin or referrer", 403)
			return
		}
		if originString != "" {
			origin, err := url.ParseRequestURI(originString)
			if err != nil {
				http.Error(resp, "403 forbidden - invalid origin", 403)
				return
			}
			if origin.Host != req.Host {
				http.Error(resp, "403 forbidden - wrong origin", 403)
				return
			}
		}
		if referrerString != "" {
			referrer, err := url.ParseRequestURI(referrerString)
			if err != nil {
				http.Error(resp, "403 forbidden - invalid referrer", 403)
				return
			}
			if referrer.Host != req.Host {
				http.Error(resp, "403 forbidden - wrong referrer", 403)
				return
			}
		}
		handler.ServeHTTP(resp, req)
	})
}

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

func returner3(
	callback func(args *pathing.Object) (interface{}, error),
) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		var args interface{}
		err := json.NewDecoder(req.Body).Decode(&args)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		val, err := callback(pathing.NewObject(unconvert(args)))
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(resp)
		enc.Encode(convert(val))
	}
}

func NewApiServer(n *core.IpfsNode, us *userstore.Userstore, db transaction.Database) (http.Handler, error) {
	apis := map[string]func(args *pathing.Object) (interface{}, error){
		"v0/name/resolve": func(args *pathing.Object) (interface{}, error) {
			hash, err := db.Resolve(args.Get("address").String())
			if err != nil {
				return nil, err
			}
			if hash[:6] == "/ipfs/" {
				hash = hash[6:]
			}
			return map[string]interface{}{
				"hash": hash,
			}, nil
		},
		"v0/block/get": func(args *pathing.Object) (interface{}, error) {
			typ, data, err := db.GetData(args.Get("hash").String())
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"type": typ,
				"data": data,
			}, nil
		},
		"v0/block/put": func(args *pathing.Object) (interface{}, error) {
			typ := args.Get("type").Uint64()
			fmt.Printf("%T\n", args.Get("type").Value())
			data := args.Get("data").Bytes()
			hash, err := db.PutData(typ, data)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"hash": hash,
			}, nil
		},
		"v0/user/list": func(args *pathing.Object) (interface{}, error) {
			keys, err := db.KeyList()
			if err != nil {
				return nil, err
			}
			val := make([]interface{}, 0)
			for _, key := range keys {
				si, err := transaction.NewSigner(context.TODO(), db, key)
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
		"v0/user/create": func(args *pathing.Object) (interface{}, error) {
			name := args.Get("name").String()
			address, err := commands.CreatePage(us, db, "USER", name)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"address": address,
				"name":    name,
			}, nil
		},
		"v0/home": func(args *pathing.Object) (interface{}, error) {
			myAddr := args.Get("myAddress").String()
			return commands.Home(us, db, myAddr)
		},
		"v0/subscribe": func(args *pathing.Object) (interface{}, error) {
			myAddr := args.Get("myAddress").String()
			addr := args.Get("address").String()
			value := args.Get("value").Bool()
			err := us.UpdateUser(myAddr, func(obj *pathing.Object) error {
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
		"v0/canPost": func(args *pathing.Object) (interface{}, error) {
			myAddr := args.Get("myAddress").String()
			addr := args.Get("address").String()
			value := args.Get("value").Bool()
			si, err := transaction.NewSigner(context.TODO(), db, myAddr)
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
		"v0/notifications": func(args *pathing.Object) (interface{}, error) {
			myAddr := args.Get("myAddress").String()
			user, err := us.GetUser(myAddr)
			if err != nil {
				return nil, err
			}
			res := map[interface{}]interface{}{
				"notifications": user.Get("notifications").Value(),
			}
			return res, nil
		},
		"v0/page": func(args *pathing.Object) (interface{}, error) {
			addr := args.Get("address").String()
			user := args.Get("myAddress").String()
			return commands.PostList(us, db, addr, user)
		},
		"v0/page/set": func(args *pathing.Object) (interface{}, error) {
			myAddr := args.Get("myAddress").String()
			key := args.Get("key").String()
			val := args.Get("val").String()
			user, err := us.GetUser(myAddr)
			if err != nil {
				return nil, err
			}
			err = user.Get("info").Get("key").Set(val)
			if err != nil {
				return nil, err
			}
			err = commands.SetInfo(db, myAddr, key, val)
			if err != nil {
				return nil, err
			}
			return true, nil
		},
		"v0/page/setwriters": func(args *pathing.Object) (interface{}, error) {
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
			err := commands.SetWriters(db, id, writers)
			if err != nil {
				return nil, err
			}
			err = us.UpdateUser(myAddr, func(obj *pathing.Object) error {
				obj.Get("writers").Set(writers2)
				return nil
			})
			if err != nil {
				return nil, err
			}
			return true, nil
		},
		"v0/page/post": func(args *pathing.Object) (interface{}, error) {
			addr := args.Get("address").String()
			myAddr := args.Get("myAddress").String()
			text := args.Get("text").String()
			attachments := args.Get("attachments")
			result, err := commands.CreatePost(db, addr, myAddr, text, attachments)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"result": result,
			}, nil
		},
		"v0/page/comment": func(args *pathing.Object) (interface{}, error) {
			addr := args.Get("address").String()
			myAddr := args.Get("myAddress").String()
			postHash := args.Get("postHash").String()
			text := args.Get("text").String()
			result, err := commands.CreateTextComment(db, addr, myAddr, postHash, text)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{
				"result": result,
			}, nil
		},
	}
	app := http.NewServeMux()
	for key, val := range apis {
		app.Handle("/api/"+key, securityCheck(http.HandlerFunc(returner3(val))))
	}
	app.Handle("/api/v0/upload", securityCheck(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		header := make([]byte, 512)
		m, err := io.ReadFull(req.Body, header)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hash, err := db.Put(io.MultiReader(bytes.NewReader(header[:m]), req.Body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c, err := cid.Parse(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		mime := http.DetectContentType(header[:m])
		if strings.HasPrefix(mime, "image/") {
			mime = "IMAGE"
		} else if strings.HasPrefix(mime, "video/") {
			mime = "VIDEO"
		} else {
			http.Error(w, "Unknown filetype", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)

		enc.Encode(map[string]interface{}{
			"t":    mime,
			"hash": c.String(),
		})
	})))
	return app, nil
}
