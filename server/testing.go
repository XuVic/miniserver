package server

import (
	"encoding/base64"

	"github.com/XuVic/miniserver/util"
)

type ResponseRecord struct {
	Content string
	Status  string
	Body    string
	Request *Request
}

func (res *ResponseRecord) Write(b []byte) (int, error) {
	if res.Status == "" {
		res.SetStatus("200")
	}
	if res.Content == "" {
		res.SetType("text")
	}
	res.Body = string(b)
	return 0, nil
}

func (res *ResponseRecord) SetStatus(code string) {
	res.Status = code
}

func (res *ResponseRecord) SetType(content string) {
	res.Content = content
}

func (res *ResponseRecord) Refresh() {
	res.Body = ""
	res.Content = ""
	res.Status = ""
}

func MockRequest(method, uri string) (*Request, *ResponseRecord) {
	now := util.TimeNow()
	localaddr := "127.0.0.1:1234"
	id := base64.StdEncoding.EncodeToString([]byte(localaddr + now))

	req := &Request{Method: method, URI: uri, Timestamp: now, Id: id}
	res := &ResponseRecord{Request: req}
	return req, res
}
