package rocketmq

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
)

func TestResolveNameServer(t *testing.T) {
	res, err := ResolveNameServer("localhost:9876")
	require.NoError(t, err)
	t.Log(neatjsons.S(res))
}
