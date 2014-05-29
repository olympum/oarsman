package main

import (
	"github.com/gocraft/web"
	"github.com/olympum/gorower/web/context"
	"net/http"
)

func main() {
	router := web.New(context.UserContext{}). // Create your router
							Middleware(web.LoggerMiddleware).                 // Use some included middleware
							Middleware(web.ShowErrorsMiddleware).             // ...
							Middleware((*context.UserContext).SetHelloCount). // Your own middleware!
							Get("/", (*context.UserContext).SayHello)         // Add a route
	http.ListenAndServe("localhost:3333", router) // Start the server!
}
