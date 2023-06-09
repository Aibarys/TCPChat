package server

import "net"

type client struct {
	Name string
	Conn net.Conn
}

type message struct {
	name    string
	text    string
	address string
}
