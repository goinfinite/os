package valueObject

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

var NewNetworkProtocol = tkValueObject.NewNetworkProtocol

var ValidNetworkProtocols = []string{
	"http", "https", "ws", "wss", "grpc", "grpcs", "tcp", "udp",
}
