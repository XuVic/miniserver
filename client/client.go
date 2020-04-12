package client

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/XuVic/miniserver/util"
)

func NewConn(addr, port string) *Client {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	if err != nil {
		return nil
	}
	return &Client{conn: conn, Addr: conn.LocalAddr().String(), Closed: false}
}

type Client struct {
	conn   net.Conn
	Addr   string
	Closed bool
}

func (c *Client) Send(r *Request) *Response {
	if c.Closed == true {
		return nil
	}

	io.WriteString(c.conn, r.String())
	res := parseResponse(r, c.conn)

	if res.Status == 408 && c.Closed == false {
		c.Closed = true
		c.conn.Close()
	}
	return res
}

func (c *Client) Get(uri string) *Response {
	r := NewRequest("Get", uri, c.Addr)

	return c.Send(r)
}

func (c *Client) Close() error {
	io.WriteString(c.conn, "QUIT")
	c.Closed = true
	return c.conn.Close()
}

func parseResponse(r *Request, conn net.Conn) *Response {
	var str strings.Builder
	char := make([]byte, 1)
	for {
		conn.Read(char)
		str.WriteString(string(char))
		char = make([]byte, 1)

		if strings.Contains(str.String(), "RES_END") {
			break
		}
	}
	lines := strings.Split(str.String(), "\n")

	status, _ := strconv.Atoi(strings.Split(lines[0], " ")[3])
	header := make(map[string]string)
	header["content_type"] = strings.Split(lines[1], " ")[1]
	header["time"] = strings.Split(lines[2], " ")[1]
	body := lines[3]

	return &Response{Request: r, Header: header, Body: body, Status: status}
}

type Request struct {
	Method    string
	URI       string
	Body      string
	Timestamp string
	Id        string
}

func (r *Request) String() string {
	return fmt.Sprintf("%s %s\nID %s\nTIME %s\n%s\nREQ_END", r.Method, r.URI, r.Id, r.Body, r.Timestamp)
}

func NewRequest(method, uri string, localaddr string) *Request {
	now := util.TimeNow()

	id := base64.StdEncoding.EncodeToString([]byte(localaddr + now))

	return &Request{Method: method, URI: uri, Timestamp: now, Id: id}
}

type Response struct {
	Request *Request
	Header  map[string]string
	Body    string
	Status  int
}

func (r *Response) Timestamp() *time.Time {
	v, ok := r.Header["time"]

	if ok == false {
		return nil
	}

	t, _ := time.Parse("2006-01-02T15:04:05", v)
	return &t
}

func (r *Response) ContentType() string {
	v, ok := r.Header["content_type"]

	if ok == false {
		return ""
	}
	return v
}
