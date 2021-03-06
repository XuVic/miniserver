package handler

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/XuVic/miniserver/server"
)

var mux *server.Mux

func init() {
	mux = Routes()
}

func TestProductsHandler(t *testing.T) {
	req, res := server.MockRequest("Get", "/products")
	mux.Serve(req, res)
	assert.Equal(t, "200", res.Status)
	assert.NotNil(t, res.Body)

	for i := 1; i <= 30; i++ {
		res.Refresh()
		mux.Serve(req, res)
	}

	assert.Equal(t, "503", res.Status)
	assert.NotNil(t, res.Body)

	time.Sleep(time.Second * 1)

	res.Refresh()
	mux.Serve(req, res)
	assert.Equal(t, "200", res.Status)
	assert.NotNil(t, res.Body)
}

func TestFactHandler(t *testing.T) {
	req, res := server.MockRequest("Get", "/facts")
	mux.Serve(req, res)
	assert.Equal(t, "200", res.Status)
	assert.NotNil(t, res.Body)

	wg := sync.WaitGroup{}

	wg.Add(30)
	for i := 1; i <= 30; i++ {
		go func() {
			res.Refresh()
			mux.Serve(req, res)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, "503", res.Status)
	assert.NotNil(t, res.Body)

	time.Sleep(time.Second * 1)

	res.Refresh()
	mux.Serve(req, res)
	assert.Equal(t, "200", res.Status)
	assert.NotNil(t, res.Body)
}
