package gateway

import (
	"io/ioutil"
	"net/http"
)

const CatEndpoint = "https://cat-fact.herokuapp.com"

type CatFactGateway struct {
	Endpoint string
}

func (g *CatFactGateway) Get(path string) string {
	res, err := http.Get(g.Endpoint + path)
	if err != nil {
		return ""
	}

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}
