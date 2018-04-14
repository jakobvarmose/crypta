package rpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jakobvarmose/crypta/pathing"
)

type message struct {
	ID     int                    `json:"id,omitempty"`
	Method string                 `json:"method,omitempty"`
	Params map[string]interface{} `json:"params,omitempty"`
	Result interface{}            `json:"result,omitempty"`
	Error  *Error                 `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error %v: %v", e.Code, e.Message)
}

type Methods = map[string]interface{}
type Value = map[string]interface{}

type RPC struct {
	conn    *websocket.Conn
	state   interface{}
	methods map[string]interface{}
	calls   map[int]chan message
	index   int
	mu      sync.Mutex
}

func New(
	conn *websocket.Conn,
	state interface{},
	methods map[string]interface{},
) *RPC {
	r := &RPC{
		conn:    conn,
		state:   state,
		methods: methods,
		calls:   make(map[int]chan message),
		index:   0,
	}
	go func() {
		defer conn.Close()
		for {
			var msg message
			err := r.conn.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				return
			}
			if msg.Method != "" {
				if msg.Params == nil {
					msg.Params = map[string]interface{}{}
				}
				method, ok := r.methods[msg.Method]
				if !ok {
					method = r.methods[""]
				}
				args := []reflect.Value{
					reflect.ValueOf(r.state),
					reflect.ValueOf(pathing.NewObject(msg.Params)),
				}
				res := reflect.ValueOf(method).Call(args)
				result := res[0].Interface()
				err, _ = res[1].Interface().(error)
				if msg.ID != 0 {
					if err != nil {
						r.conn.WriteJSON(message{
							ID: msg.ID,
							Error: &Error{
								Code:    500,
								Message: err.Error(),
							},
						})
					} else {
						r.conn.WriteJSON(message{
							ID:     msg.ID,
							Result: result,
						})
					}
				}
			} else {
				r.mu.Lock()
				if call, ok := r.calls[msg.ID]; ok {
					call <- msg
				}
				r.mu.Unlock()
			}
		}
	}()
	return r
}

func (r *RPC) Call(ctx context.Context, method string, params map[string]interface{}) (*pathing.Object, error) {
	ch := make(chan message)
	r.mu.Lock()
	r.index++
	r.calls[r.index] = ch
	r.mu.Unlock()
	err := r.conn.WriteJSON(message{
		ID:     r.index,
		Method: method,
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	select {
	case res := <-ch:
		r.mu.Lock()
		delete(r.calls, r.index)
		r.mu.Unlock()
		if res.Error != nil {
			return nil, res.Error
		}
		return pathing.NewObject(res.Result), nil
	case <-ctx.Done():
		r.mu.Lock()
		delete(r.calls, r.index)
		r.mu.Unlock()
		return nil, ctx.Err()
	}
}

func (r *RPC) Notify(method string, params map[string]interface{}) error {
	return r.conn.WriteJSON(message{
		Method: method,
		Params: params,
	})
}
