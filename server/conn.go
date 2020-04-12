package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/XuVic/miniserver/util"
)

func newConn(e *Engine, conn net.Conn) *connContext {
	id := base64.StdEncoding.EncodeToString([]byte(conn.RemoteAddr().String()))
	return &connContext{engine: e, conn: conn, remoteAddr: conn.RemoteAddr().String(), id: id}
}

type connContext struct {
	conn       net.Conn
	engine     *Engine
	remoteAddr string
	mu         sync.Mutex
	id         string
}

func (c *connContext) server() {
	defer func() {
		c.engine.CLog(c.remoteAddr).Printf("Connection Close!")
		c.engine.Stat.RemoveConn()
		c.conn.Close()
	}()

	chRequest := make(chan *Request)

	go func() {
		for {
			res, leave := c.readRequest(chRequest)
			if leave {
				chRequest <- nil
				break
			} else {
				chRequest <- res
			}
		}
	}()

	connClose := false
	for !connClose {
		select {
		case request := <-chRequest:
			if request == nil {
				connClose = true
			} else {
				response := c.createResponse(request)
				c.engine.CLog(c.remoteAddr).Printf("Recieve %s", request.RawString())
				go c.engine.Handler.Serve(request, response.Writer())
			}
		case <-time.After(c.engine.TimeOut):
			c.engine.HandleTimeOut(c)
			connClose = true
		}

	}
}

func (c *connContext) readRequest(ch chan<- *Request) (*Request, bool) {

	var str strings.Builder
	char := make([]byte, 1)
	eoc := false
	for {
		c.conn.Read(char)
		str.WriteString(string(char))
		char = make([]byte, 1)

		if strings.Contains(str.String(), "REQ_END") {
			break
		}
		if strings.Contains(str.String(), "QUIT") {
			eoc = true
			break
		}
	}

	lines := strings.Split(str.String(), "\n")

	if eoc {
		return nil, eoc
	}

	return parseRequest(lines), false
}

func (c *connContext) createResponse(r *Request) *response {
	buf := make([]byte, 0, 10)
	return &response{context: c, buffer: bytes.NewBuffer(buf), request: r}
}

func parseRequest(str []string) *Request {
	lines := make([][]string, len(str))
	for i, s := range str {
		lines[i] = strings.Split(s, " ")
	}
	method := lines[0][0]
	uri := lines[0][1]
	id := lines[1][1]
	time := lines[2][1]

	var body string
	if len(lines) > 2 {
		body = lines[2][0]
	}

	return &Request{URI: uri, Method: method, Id: id, Timestamp: time, body: body}
}

type Request struct {
	URI       string
	Method    string
	body      string
	Id        string
	Timestamp string
}

func (r *Request) RawString() string {
	return fmt.Sprintf("%s %s", r.Method, r.URI)
}

type response struct {
	request *Request
	context *connContext
	buffer  *bytes.Buffer
	content string
	body    string
}

func (res *response) Writer() *responseWriter {
	return &responseWriter{response: res}
}

func (res *response) flush() (int, error) {

	responseStr := fmt.Sprintf("CONTENT_TYPE %s\nTIME %s\n%s\nRES_END", res.content, util.TimeNow(), res.body)
	n, err := io.WriteString(res.buffer, responseStr)
	if err != nil {
		return n, err
	}
	return res.context.conn.Write(res.buffer.Bytes())
}

type ResponseWriter interface {
	Write(b []byte) (int, error)
	SetStatus(string)
	SetType(content string)
}

type responseWriter struct {
	response *response
	status   string
	content  string
}

func (res *responseWriter) Write(b []byte) (int, error) {
	if res.status == "" {
		res.SetStatus("200")
	}
	if res.content == "" {
		res.SetType("text")
	}
	res.response.body = string(b)
	return res.response.flush()
}

func (res *responseWriter) SetStatus(code string) {
	res.status = code
	requset_str := res.response.request.RawString()
	firstLine := fmt.Sprintf("%s Status %s\n", requset_str, res.status)
	io.WriteString(res.response.buffer, firstLine)
}

func (res *responseWriter) SetType(content string) {
	res.content = content
	res.response.content = content
}
