package server

import (
	"testing"
	"time"

	"github.com/XuVic/miniserver/client"
	"github.com/stretchr/testify/assert"
)

func testHandler(r *Request, res ResponseWriter) {
	res.Write([]byte("Welcome to use MiniServer!"))
}

func init() {
	mux := NewMux()
	mux.HandleFunc("/", testHandler)
	engine := NewEngine("localhost", "3000")
	mux.HandleFunc("/stat", engine.HandleStat)
	engine.Handler = mux
	engine.TimeOut = time.Second * 1
	go engine.RunTCP()
	time.Sleep(time.Millisecond * 100)
}

func TestRunServer(t *testing.T) {

	client := client.NewConn("localhost", "3000")
	assert.NotNil(t, client)

	res := client.Get("/")
	assert.Equal(t, "Welcome to use MiniServer!", res.Body)
	assert.NotNil(t, res.Timestamp())
	assert.NotEqual(t, "", res.ContentType())
	assert.Equal(t, 200, res.Status)
	if client != nil {
		client.Close()
	}
}

func TestRequestTimeOut(t *testing.T) {
	client := client.NewConn("localhost", "3000")
	assert.NotNil(t, client)
	time.Sleep(time.Second * 3)
	res := client.Get("/")
	assert.Equal(t, 408, res.Status)
	if client != nil {
		client.Close()
	}
}

func TestStatRequest(t *testing.T) {
	client := client.NewConn("localhost", "3000")
	assert.NotNil(t, client)
	res := client.Get("/stat")
	assert.Equal(t, 200, res.Status)
	assert.NotNil(t, res.Body)
	if client != nil {
		client.Close()
	}
}
