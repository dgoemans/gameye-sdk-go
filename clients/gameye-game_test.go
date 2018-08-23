package clients

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/Gameye/gameye-sdk-go/models"

	"github.com/Gameye/gameye-sdk-go/test"
	"github.com/stretchr/testify/assert"
)

func TestGameyeClient_SubscribeGame(t *testing.T) {
	var err error
	defer func() {
		assert.NoError(t, err)
	}()

	patchChannel := make(chan string, 1)
	mux := test.CreateAPITestServerMux(
		`{}`, patchChannel,
	)

	var listener net.Listener
	listener, err = net.Listen("tcp", ":8084")
	if err != nil {
		return
	}
	defer listener.Close()
	go http.Serve(listener, mux)

	client := NewGameyeClient(GameyeClientConfig{
		Endpoint: "http://localhost:8084",
		Token:    "",
	})

	var sub *GameQuerySubscription
	sub, err = client.SubscribeGame()
	if err != nil {
		return
	}
	defer sub.Cancel()

	{
		patchChannel <- fmt.Sprintf(`[{"path":[],"value":%s}]`, strings.Replace(models.GameStateJSONMock, "\n", "", -1))
		var state *models.GameQueryState
		state, err = sub.NextState()
		if err != nil {
			return
		}

		assert.Equal(t, &models.GameStateMock, state)
	}
}

func TestGameyeClient_QueryGame(t *testing.T) {
	var err error
	defer func() {
		assert.NoError(t, err)
	}()

	mux := test.CreateAPITestServerMux(
		models.GameStateJSONMock, nil,
	)

	var listener net.Listener
	listener, err = net.Listen("tcp", ":8085")
	if err != nil {
		return
	}
	defer listener.Close()
	go http.Serve(listener, mux)

	client := NewGameyeClient(GameyeClientConfig{
		Endpoint: "http://localhost:8085",
		Token:    "",
	})

	var state *models.GameQueryState
	state, err = client.QueryGame()
	if err != nil {
		return
	}
	assert.Equal(t, &models.GameStateMock, state)
}