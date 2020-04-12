package handler

import (
	"encoding/json"
	"io"
	"time"

	"github.com/XuVic/miniserver/gateway"
	"github.com/XuVic/miniserver/server"
)

var mockGateway *gateway.LimitGateway = gateway.Factory.Build(gateway.Mock, 30, time.Second)

var catGateway *gateway.LimitGateway = gateway.Factory.Build(gateway.CatFact, 3, time.Second)

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

func FactsHandler(r *server.Request, res server.ResponseWriter) {
	catfacts, ok := catGateway.Get("/facts")
	if ok {
		var data map[string]interface{}
		json.Unmarshal([]byte(catfacts), &data)

		facts := map[string][]interface{}{
			"facts": data["all"].([]interface{})[:5],
		}
		res.SetType("json")
		factsJSON, _ := json.Marshal(facts)
		io.WriteString(res, string(factsJSON))
	} else {
		res.SetStatus("503")
		io.WriteString(res, "Resource are unavailable!")
	}
}
