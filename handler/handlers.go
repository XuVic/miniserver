package handler

import (
	"io"
	"time"

	"github.com/XuVic/miniserver/gateway"
	"github.com/XuVic/miniserver/server"
)

var mockGateway *gateway.LimitGateway = gateway.Factory.Build(gateway.Mock, 3, time.Second*2)

func IndexHandler(r *server.Request, res server.ResponseWriter) {
	res.Write([]byte("Welcome to use MiniServer!"))
}

func ProductsHandler(r *server.Request, res server.ResponseWriter) {
	products, ok := mockGateway.Get("/products")
	if ok {
		res.SetType("json")
		io.WriteString(res, products)
	} else {
		res.SetStatus("503")
		io.WriteString(res, "Resource are unavailable!")
	}
}
