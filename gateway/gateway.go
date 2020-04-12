package gateway

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Gateway interface {
	Get(string) string
}

const (
	Mock = iota
)

var Factory *GatewayFactory = NewFactory()

func NewFactory() *GatewayFactory {
	gateway := make(map[int]*LimitGateway)
	chs := make(map[int]chan time.Duration)

	chs[0] = make(chan time.Duration)

	f := &GatewayFactory{Gateways: gateway, chs: chs}

	go func() {
		for {
			select {
			case <-f.chs[0]:
				go Wait(f.Gateways[0])

			}
		}
	}()

	return f
}

func Wait(lGateway *LimitGateway) {
	time.Sleep(lGateway.Wait)
	lGateway.m.Lock()
	lGateway.Num = 0
	lGateway.isWait = false
	lGateway.m.Unlock()
}

type GatewayFactory struct {
	Gateways map[int]*LimitGateway
	chs      map[int]chan time.Duration
}

func (f *GatewayFactory) Build(g, limit int, wait time.Duration) *LimitGateway {
	_, ok := f.Gateways[g]

	if ok == false {
		switch g {
		case 0:
			mock := DefaultMockGateway()
			f.Gateways[g] = NewLimitGateway(limit, wait, mock, f.chs[g])
		default:
			panic("Error")
		}
	}

	return f.Gateways[g]
}

func NewLimitGateway(limit int, wait time.Duration, g Gateway, ch chan<- time.Duration) *LimitGateway {
	lGateway := &LimitGateway{Limit: limit, Wait: wait, gateway: g, m: &sync.Mutex{}, ch: ch}
	return lGateway
}

type LimitGateway struct {
	Num, Limit int
	Wait       time.Duration
	gateway    Gateway
	m          *sync.Mutex
	ch         chan<- time.Duration
	isWait     bool
}

func (g *LimitGateway) Get(path string) (string, bool) {
	g.m.Lock()
	if g.Num >= g.Limit {
		if g.isWait == false {
			g.ch <- g.Wait
			g.isWait = true
		}
		g.m.Unlock()
		return "", false
	}
	g.Num++
	g.m.Unlock()
	res := g.gateway.Get(path)
	if res == "" {
		return "", false
	}
	return res, true
}

func DefaultMockGateway() *MockGateway {
	data := make(map[string][]map[string]string)

	productsName := []string{"iPhoneX", "GooglePixel", "SamSumgGalaxy"}

	products := make([]map[string]string, len(productsName))

	for i := range products {
		products[i] = map[string]string{
			strconv.Itoa(i): productsName[i],
		}
	}
	data["products"] = products

	return &MockGateway{Data: data}
}

type MockGateway struct {
	Data map[string][]map[string]string
}

func (g *MockGateway) parsePath(path string) string {
	tokens := strings.Split(path, "/")

	return tokens[len(tokens)-1]
}

func (g *MockGateway) Get(path string) string {
	path = g.parsePath(path)
	data := map[string]interface{}{path: g.Data[path]}
	json, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(json)
}
