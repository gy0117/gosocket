package internal

type Pair struct {
	Key   string
	Value string
}

var (
	ConnectionPair          = Pair{"Connection", "Upgrade"}
	UpgradePair             = Pair{"Upgrade", "websocket"}
	SecWebSocketVersionPair = Pair{"Sec-WebSocket-Version", "13"}
	SecWebSocketKeyPair     = Pair{"Sec-WebSocket-Key", ""}
)
