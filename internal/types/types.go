package types

type Pair struct {
	Key   string
	Value string
}

var (
	ConnectionPair             = Pair{"Connection", "Upgrade"}
	UpgradePair                = Pair{"Upgrade", "websocket"}
	SecWebSocketVersionPair    = Pair{"Sec-WebSocket-Version", "13"}
	SecWebSocketKeyPair        = Pair{"Sec-WebSocket-Key", ""}
	SecWebSocketAcceptPair     = Pair{"Sec-WebSocket-Accept", ""}
	SecWebSocketExtensionsPair = Pair{"Sec-WebSocket-Extensions", "permessage-deflate"}
)

const PermessageDeflate = "permessage-deflate"
