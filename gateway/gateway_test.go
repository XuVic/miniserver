package gateway

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimitGateway(t *testing.T) {
	mockGateway := LimitGatewayFactory.Build(Mock, 1, time.Second*2)
	_, ok := mockGateway.Get("products")
	assert.True(t, ok)
	mockGateway = LimitGatewayFactory.Build(Mock, 1, time.Second*2)
	_, ok = mockGateway.Get("products")
	assert.False(t, ok)
	_, ok = mockGateway.Get("products")
	assert.False(t, ok)
	time.Sleep(time.Second * 2)

	_, ok = mockGateway.Get("products")
	assert.True(t, ok)
}
