package main

import (
	"fmt"
	"net"
)

type channelPool struct {
}

type PoolConn struct {
	net.Conn
	c *channelPool
}

func (c *channelPool) wrapConn(conn net.Conn) net.Conn {
	p := &PoolConn{c: c}
	p.Conn = conn
	return p
}

func (p *PoolConn) Close() error {
	fmt.Print("I am close...")
	return nil
}

func main() {
	c := &channelPool{}

	p := &PoolConn{
		c: c,
	}

	fmt.Println(p.c)
	fmt.Println(p.Close())
}
