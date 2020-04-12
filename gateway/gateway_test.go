package gateway

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimitGateway(t *testing.T) {
	mockGateway := Factory.Build(Mock, 1, time.Millisecond)
	_, ok := mockGateway.Get("/products")
	assert.True(t, ok)

	mockGateway2 := Factory.Build(Mock, 1, time.Millisecond)
	assert.Equal(t, mockGateway, mockGateway2)

	_, ok = mockGateway2.Get("/products")
	assert.False(t, ok)
	_, ok = mockGateway2.Get("/products")
	assert.False(t, ok)
	time.Sleep(time.Millisecond * 1)

	_, ok = mockGateway2.Get("/products")
	assert.True(t, ok)
}

func TestCatFactGateway(t *testing.T) {
	gateway := Factory.Build(CatFact, 1, time.Millisecond)

	gateway2 := Factory.Build(CatFact, 1, time.Millisecond)

	assert.True(t, reflect.DeepEqual(gateway, gateway2))

	res, ok := gateway.Get("/facts")
	assert.True(t, ok)
	assert.NotNil(t, res)

	res, ok = gateway.Get("/facts")
	assert.False(t, ok)
	assert.Equal(t, "", res)
	time.Sleep(time.Millisecond * 1)
	res, ok = gateway.Get("/facts")
	assert.True(t, ok)
	assert.NotNil(t, res)

}
