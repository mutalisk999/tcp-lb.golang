package main

import (
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"net"
	"sort"
)

func handleTcpProxyConn(g goroutine_mgr.Goroutine, a interface{}) {
	defer g.OnQuit()

	var connToTarget *net.TCPConn = nil
	var targetCopy LBTargetCopy

	conn := a.(*net.TCPConn)
	targetsCopy := LBTargetsMgrP.DumpTargetsCopy()
	sort.Sort(LBTargetCopys(targetsCopy))

	for _, t := range targetsCopy {
		if t.Status != 0 {
			continue
		}
		if t.ConnCount >= t.MaxConnCount {
			continue
		}

		targetAddr, err := net.ResolveTCPAddr("tcp", t.EndPointConn)
		if err != nil {
			Error.Fatalf("Error: %v", err)
		}
		connToTarget, err = net.DialTCP("tcp", nil, targetAddr)
		if err != nil {
			Warn.Printf("Can not connect to target: %s", t.EndPointConn)
			continue
		} else {
			Info.Printf("Succeed connect to target: %s", t.EndPointConn)
			targetCopy = t
			break
		}
	}

	if connToTarget == nil {
		Warn.Printf("Can not connect to any target endpoint, Close node connection")
		_ = conn.Close()
		LBNodeP.ComsumeNewConn()

	} else {
		targetId := CaclTargetId(targetCopy.EndPointConn)

		var nodeConn NodeConnection
		nodeConn.Initialise(conn, LBNodeP.timeout)

		var targetConn TargetConnection
		targetConn.Initialise(connToTarget, targetCopy.Timeout, targetId)

		LBConnectionPairMgrP.AddConnectionPair(&nodeConn, &targetConn)

		LBGoroutineManagerP.GoroutineCreateP2("tcp_node_data", handleNodeData, &nodeConn, &targetConn)
		LBGoroutineManagerP.GoroutineCreateP2("tcp_target_data", handleTargetData, &nodeConn, &targetConn)
	}
}

func startTcpProxy(g goroutine_mgr.Goroutine, a interface{}) {
	defer g.OnQuit()

	cfg := a.(*Config)
	addr, err := net.ResolveTCPAddr("tcp", cfg.Node.ListenEndPoint)
	if err != nil {
		Error.Fatalf("Error: %v", err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		Error.Fatalf("Error: %v", err)
	}
	defer server.Close()
	Info.Printf("Node listening on %s", cfg.Node.ListenEndPoint)

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			continue
		}
		_ = conn.SetKeepAlive(true)

		// TODO
		//ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		// banned and need to ban

		LBNodeP.ProductNewConn()
		LBGoroutineManagerP.GoroutineCreateP1("tcp_proxy_conn", handleTcpProxyConn, conn)
	}
}
