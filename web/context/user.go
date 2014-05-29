package context

import (
	"encoding/json"
	"github.com/gocraft/web"
)

type UserContext struct {
	HelloCount int `json:"count"`
}

func (c *UserContext) SetHelloCount(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	c.HelloCount = 3
	next(rw, req)
}

func (c *UserContext) SayHello(rw web.ResponseWriter, req *web.Request) {
	json, _ := json.Marshal(c)
	headers := rw.Header()
	headers.Set("Content-Type", "application/json;charset-utf8")
	rw.WriteHeader(200)
	rw.Write(json)
}
