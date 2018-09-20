package server

import (
	"context"
	"net/http"

	"github.com/ipsn/go-ipfs/core"

	"github.com/jakobvarmose/crypta/controllers"
	"github.com/jakobvarmose/crypta/pathing"
	"github.com/jakobvarmose/crypta/transaction"
	"github.com/jakobvarmose/crypta/userstore"
)

func returner3(
	p controllers.Provider,
	callback func(p controllers.Provider, args *pathing.Object) (interface{}, error),
) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		args, err := ReadJSON(req.Body)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusBadRequest)
			return
		}
		val, err := callback(p, pathing.NewObject(args))
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		resp.Header().Set("Content-Type", "application/json")
		WriteJSON(resp, val)
	}
}

type provider struct {
	n   *core.IpfsNode
	us  *userstore.Userstore
	db  transaction.Database
	ctx context.Context
}

func (p *provider) N() *core.IpfsNode {
	return p.n
}
func (p *provider) US() *userstore.Userstore {
	return p.us
}
func (p *provider) DB() transaction.Database {
	return p.db
}
func (p *provider) CTX() context.Context {
	return p.ctx
}

func NewApiServer(n *core.IpfsNode, us *userstore.Userstore, db transaction.Database) (http.Handler, error) {
	apis := map[string]func(p controllers.Provider, args *pathing.Object) (interface{}, error){
		"v0/name/resolve": controllers.NameResolve,

		"v0/block/get": controllers.BlockGet,
		"v0/block/put": controllers.BlockPut,

		"v0/user/list":   controllers.UserList,
		"v0/user/create": controllers.UserCreate,

		"v0/canPost":       controllers.CanPost,
		"v0/home":          controllers.Home,
		"v0/notifications": controllers.Notifications,
		"v0/subscribe":     controllers.Subscribe,

		"v0/page":            controllers.Page,
		"v0/page/set":        controllers.PageSet,
		"v0/page/setwriters": controllers.PageSetWriters,
		"v0/page/post":       controllers.PagePost,
		"v0/page/comment":    controllers.PageComment,
	}
	app := http.NewServeMux()
	p := &provider{n, us, db, context.TODO()}
	for key, val := range apis {
		app.Handle("/api/"+key, securityCheck(http.HandlerFunc(returner3(p, val))))
	}
	app.Handle("/api/v0/upload", securityCheck(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		controllers.Upload(p, w, req)
	})))
	return app, nil
}
