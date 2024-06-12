package client

type Config struct {
	Server   string
	Username string
	Password string

	// Engine to use, avaliable: http, websocket, tcp, grpc
	Engine string
}
