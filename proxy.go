package main

import (
	"log"
	"net"
)

func StartTcpProxy(c *Config) {
	addr, err := net.ResolveTCPAddr("tcp", c.Node.ListenEndPoint)
	if err != nil {
		Error.Fatalf("Error: %v", err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		Error.Fatalf("Error: %v", err)
	}
	defer server.Close()

	Info.Printf("Node listening on %s", c.Node.ListenEndPoint)

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			continue
		}

		LBNodeP.IncConnCount()
		log.Printf("Node connection count: [%d/%d]", LBNodeP.GetConnCount(), LBNodeP.GetMaxConnCount())
		if LBNodeP.GetConnCount() > LBNodeP.GetMaxConnCount() {
			Debug.Printf("Close Node connection, connection count: [%d/%d]", LBNodeP.GetConnCount(), LBNodeP.GetMaxConnCount())
			_ = conn.Close()
			LBNodeP.DecConnCount()
			continue
		}
		_ = conn.SetKeepAlive(true)

		// TODO
		//ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		// banned and need to ban

	}
}
