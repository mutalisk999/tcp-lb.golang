package main

import (
	"net"
	"sync"
	"time"
)

type NodeConnection struct {
	mutex              *sync.RWMutex
	conn               *net.TCPConn
	timeout            uint32
	periodStartTime1m  time.Time
	periodStartTime5m  time.Time
	periodStartTime30m time.Time
	readBytes1m        uint64
	readBytes5m        uint64
	readBytes30m       uint64
	sendBytes1m        uint64
	sendBytes5m        uint64
	sendBytes30m       uint64
}

func (c *NodeConnection) Initialise(conn *net.TCPConn, timeout uint32) {
	c.mutex = new(sync.RWMutex)
	c.conn = conn
	c.timeout = timeout
	c.periodStartTime1m = time.Now()
	c.periodStartTime5m = time.Now()
	c.periodStartTime30m = time.Now()
	c.readBytes1m = 0
	c.readBytes5m = 0
	c.readBytes30m = 0
	c.sendBytes1m = 0
	c.sendBytes5m = 0
	c.sendBytes30m = 0
}

func (c *NodeConnection) SetKeepAlive() {
	c.mutex.Lock()
	_ = c.conn.SetKeepAlive(true)
	c.mutex.Unlock()
}

func (c *NodeConnection) GetConnection() *net.TCPConn {
	var conn *net.TCPConn
	c.mutex.Lock()
	conn = c.conn
	c.mutex.Unlock()
	return conn
}

func (c *NodeConnection) GetTimeOut() uint32 {
	var timeout uint32
	c.mutex.Lock()
	timeout = c.timeout
	c.mutex.Unlock()
	return timeout
}

func (c *NodeConnection) Destroy() {
	c.mutex = nil
	c.conn = nil
	c.timeout = 0
	c.periodStartTime1m = time.Now()
	c.periodStartTime5m = time.Now()
	c.periodStartTime30m = time.Now()
	c.readBytes1m = 0
	c.readBytes5m = 0
	c.readBytes30m = 0
	c.sendBytes1m = 0
	c.sendBytes5m = 0
	c.sendBytes30m = 0
}

type TargetConnection struct {
	mutex              *sync.RWMutex
	conn               *net.TCPConn
	timeout            uint32
	periodStartTime1m  time.Time
	periodStartTime5m  time.Time
	periodStartTime30m time.Time
	readBytes1m        uint64
	readBytes5m        uint64
	readBytes30m       uint64
	sendBytes1m        uint64
	sendBytes5m        uint64
	sendBytes30m       uint64
	targetId           string
}

func (c *TargetConnection) Initialise(conn *net.TCPConn, timeout uint32, targetId string) {
	c.mutex = new(sync.RWMutex)
	c.conn = conn
	c.timeout = timeout
	c.periodStartTime1m = time.Now()
	c.periodStartTime5m = time.Now()
	c.periodStartTime30m = time.Now()
	c.readBytes1m = 0
	c.readBytes5m = 0
	c.readBytes30m = 0
	c.sendBytes1m = 0
	c.sendBytes5m = 0
	c.sendBytes30m = 0
	c.targetId = targetId
}

func (c *TargetConnection) SetKeepAlive() {
	c.mutex.Lock()
	_ = c.conn.SetKeepAlive(true)
	c.mutex.Unlock()
}

func (c *TargetConnection) Destroy() {
	c.mutex = nil
	c.conn = nil
	c.timeout = 0
	c.periodStartTime1m = time.Now()
	c.periodStartTime5m = time.Now()
	c.periodStartTime30m = time.Now()
	c.readBytes1m = 0
	c.readBytes5m = 0
	c.readBytes30m = 0
	c.sendBytes1m = 0
	c.sendBytes5m = 0
	c.sendBytes30m = 0
	c.targetId = ""
}

type LBConnectionPairMgr struct {
	mutex                   *sync.RWMutex
	nodeConnToTargetConnMap map[*NodeConnection]*TargetConnection
	targetConnToNodeConnMap map[*TargetConnection]*NodeConnection
}

func (l *LBConnectionPairMgr) Initialise() {
	l.mutex = new(sync.RWMutex)
	l.nodeConnToTargetConnMap = make(map[*NodeConnection]*TargetConnection)
	l.targetConnToNodeConnMap = make(map[*TargetConnection]*NodeConnection)
}

func (l *LBConnectionPairMgr) GetNode2TargetPairCount() int {
	var count int
	l.mutex.RLock()
	count = len(l.nodeConnToTargetConnMap)
	l.mutex.RUnlock()
	return count
}

func (l *LBConnectionPairMgr) GetTarget2NodePairCount() int {
	var count int
	l.mutex.RLock()
	count = len(l.targetConnToNodeConnMap)
	l.mutex.RUnlock()
	return count
}

func (l *LBConnectionPairMgr) AddConnectionPair(nodeConn *NodeConnection, targetConn *TargetConnection) {
	if nodeConn == nil || targetConn == nil {
		return
	}

	l.mutex.Lock()
	delete(l.nodeConnToTargetConnMap, nodeConn)
	delete(l.targetConnToNodeConnMap, targetConn)
	l.nodeConnToTargetConnMap[nodeConn] = targetConn
	l.targetConnToNodeConnMap[targetConn] = nodeConn
	l.mutex.Unlock()
}

func (l *LBConnectionPairMgr) RemoveByNodeConn(nodeConn *NodeConnection) {
	if nodeConn == nil {
		return
	}

	l.mutex.Lock()
	targetConn, ok := l.nodeConnToTargetConnMap[nodeConn]
	if ok {
		delete(l.nodeConnToTargetConnMap, nodeConn)
		delete(l.targetConnToNodeConnMap, targetConn)
	}
	l.mutex.Unlock()
}

func (l *LBConnectionPairMgr) RemoveByTargetConn(targetConn *TargetConnection) {
	if targetConn == nil {
		return
	}

	l.mutex.Lock()
	nodeConn, ok := l.targetConnToNodeConnMap[targetConn]
	if ok {
		delete(l.nodeConnToTargetConnMap, nodeConn)
		delete(l.targetConnToNodeConnMap, targetConn)
	}
	l.mutex.Unlock()
}

func (l *LBConnectionPairMgr) Destroy() {
	l.mutex = nil
	l.nodeConnToTargetConnMap = nil
	l.targetConnToNodeConnMap = nil
}
