package clients

import (
	"net"
	"net/http"
	"testing"

	"github.com/Gameye/gameye-sdk-go/test"
	"github.com/stretchr/testify/assert"
)

func TestGameyeClient_subscribe(t *testing.T) {
	var err error
	defer func() {
		assert.NoError(t, err)
	}()

	patchChannel := make(chan string, 1)
	mux := test.CreateAPITestServerMux(
		`{}`, patchChannel,
	)

	var listener net.Listener
	listener, err = net.Listen("tcp", ":8083")
	if err != nil {
		return
	}
	defer listener.Close()
	go http.Serve(listener, mux)

	client := NewGameyeClient(GameyeClientConfig{
		Endpoint: "http://localhost:8083",
		Token:    "",
	})

	var qs *querySubscription
	qs, err = client.subscribe("noop", nil)
	if err != nil {
		return
	}
	defer qs.Cancel()

	{
		expected := map[string]interface{}{
			"a": map[string]interface{}{
				"b": "c",
			},
		}
		patchChannel <- `[{"path":[],"value":{"a":{"b":"c"}}}]`
		var actual map[string]interface{}
		actual, err = qs.NextState()
		if err != nil {
			return
		}

		assert.Equal(t, expected, actual)
	}

	{
		expected := map[string]interface{}{
			"a": map[string]interface{}{
				"b": "d",
			},
		}
		patchChannel <- `[{"path":["a","b"],"value":"d"}]`
		var actual map[string]interface{}
		actual, err = qs.NextState()
		if err != nil {
			return
		}

		assert.Equal(t, expected, actual)
	}

}