package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/XuVic/miniserver/util"
)

func handleTimeOut(c *connContext) {
	r := &Request{Method: "None", URI: "/"}
	buf := make([]byte, 0, 10)
	res := &response{context: c, buffer: bytes.NewBuffer(buf), request: r}
	writer := res.Writer()
	writer.SetStatus("408")
	io.WriteString(writer, "Time out! Connection is closed.")
}

// RunMiniServer runs a TCP server on the given address.
func RunMiniServer(addr, port string, handler Handler) error {
	serveraddr := fmt.Sprintf("%s:%s", addr, port)
	engine := &Engine{Handler: handler, Stat: &stat{}, ConnLog: NewConnLog(serveraddr),
		TimeOut: time.Second * 60, HandleTimeOut: handleTimeOut}
	return engine.RunTCP()
}

// Handler interface regulates handler type.
type Handler interface {
	Serve(*Request, ResponseWriter)
}

// NewEngine to instantiate an Engine.
func NewEngine(addr, port string) *Engine {
	serveraddr := fmt.Sprintf("%s:%s", addr, port)
	return &Engine{Stat: &stat{}, HandleTimeOut: handleTimeOut, ConnLog: NewConnLog(serveraddr), Addr: addr, Port: port}
}

// Engine is the core component that serves as configuration and management connections.
type Engine struct {
	Addr string
	Port string

	HandleTimeOut func(c *connContext)
	TimeOut       time.Duration
	Handler       Handler
	ReadTimeout   time.Duration
	*ConnLog
	Stat *stat
}

// HandleStat a handler function is used to respond to statistical data.
func (e *Engine) HandleStat(req *Request, res ResponseWriter) {
	io.WriteString(res, e.Stat.String())
}

// RunTCP run a TCP server
func (e *Engine) RunTCP() error {
	if e.Port == "" {
		e.Port = "3030"
	}
	ln, err := net.Listen("tcp", e.Addr+":"+e.Port)
	if err != nil {
		return err
	}
	e.Log().Printf("Miniserver is running in %s:%s", e.Addr, e.Port)
	return e.Serve(ln)
}

// Serve method bind the engine to a specific listener.
func (e *Engine) Serve(ln net.Listener) error {
	ln = &onceCloseListener{Listener: ln}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		e.Stat.AddConn()
		c := newConn(e, conn)
		e.Log().Printf("Accept a connect from %s", c.remoteAddr)
		e.Log().Print(e.Stat.String())
		go c.server()
	}

}

// NewConnLog instantiate a ConnLog.
func NewConnLog(addr string) *ConnLog {
	connLog := &ConnLog{serveraddr: addr}
	log := log.New(connLog, "", 0)
	connLog.Logger = log
	return connLog
}

// ConnLog is used to manage log.
type ConnLog struct {
	remoteaddr string
	serveraddr string
	*log.Logger
}

// Log method print log message on server level.
func (l *ConnLog) Log() *ConnLog {
	l.remoteaddr = ""
	return l
}

// CLog method print log message on connection level.
func (l *ConnLog) CLog(addr string) *ConnLog {
	l.remoteaddr = addr
	return l
}

func (l *ConnLog) Write(bytes []byte) (int, error) {
	if l.remoteaddr == "" {
		return fmt.Printf("%s(s) [%s] %s", l.serveraddr, util.TimeNow(), string(bytes))
	}
	return fmt.Printf("%s(c) [%s] %s", l.remoteaddr, util.TimeNow(), string(bytes))
}

type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (l *onceCloseListener) OnceClose() error {
	l.once.Do(l.close)
	return l.closeErr
}

func (l *onceCloseListener) close() {
	l.closeErr = l.Listener.Close()
}

type stat struct {
	AliveConn int32
	BuiltConn int32
}

func (s *stat) AddConn() int32 {
	atomic.AddInt32(&s.AliveConn, 1)
	return atomic.AddInt32(&s.BuiltConn, 1)
}

func (s *stat) RemoveConn() int32 {
	return atomic.AddInt32(&s.AliveConn, -1)
}

func (s *stat) String() string {
	return fmt.Sprintf("Alive Conn: %d; Built Conn: %d", s.AliveConn, s.BuiltConn)
}
