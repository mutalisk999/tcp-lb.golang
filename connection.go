package main

import (
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
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
	writeBytes1m       uint64
	writeBytes5m       uint64
	writeBytes30m      uint64
}

type NodeConnectionCopy struct {
	Timeout            uint32
	PeriodStartTime1m  time.Time
	PeriodStartTime5m  time.Time
	PeriodStartTime30m time.Time
	ReadBytes1m        uint64
	ReadBytes5m        uint64
	ReadBytes30m       uint64
	WriteBytes1m       uint64
	WriteBytes5m       uint64
	WriteBytes30m      uint64
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
	c.writeBytes1m = 0
	c.writeBytes5m = 0
	c.writeBytes30m = 0
}

func (c *NodeConnection) SetKeepAlive() {
	c.mutex.Lock()
	_ = c.conn.SetKeepAlive(true)
	c.mutex.Unlock()
}

func (c *NodeConnection) GetConnection() *net.TCPConn {
	var conn *net.TCPConn
	c.mutex.RLock()
	conn = c.conn
	c.mutex.RUnlock()
	return conn
}

func (c *NodeConnection) GetTimeOut() uint32 {
	var timeout uint32
	c.mutex.RLock()
	timeout = c.timeout
	c.mutex.RUnlock()
	return timeout
}

func (c *NodeConnection) DumpToNodeConnectionCopy() NodeConnectionCopy {
	var connCopy NodeConnectionCopy
	c.mutex.RLock()
	connCopy.Timeout = c.timeout
	connCopy.PeriodStartTime1m = c.periodStartTime1m
	connCopy.PeriodStartTime5m = c.periodStartTime5m
	connCopy.PeriodStartTime30m = c.periodStartTime30m
	connCopy.ReadBytes1m = c.readBytes1m
	connCopy.ReadBytes5m = c.readBytes5m
	connCopy.ReadBytes30m = c.readBytes30m
	connCopy.WriteBytes1m = c.writeBytes1m
	connCopy.WriteBytes5m = c.writeBytes5m
	connCopy.WriteBytes30m = c.writeBytes30m
	c.mutex.RUnlock()
	return connCopy
}

func (c *NodeConnection) IncReadBytes(readn uint64) {
	c.mutex.Lock()
	c.readBytes1m += readn
	c.readBytes5m += readn
	c.readBytes30m += readn
	c.mutex.Unlock()
}

func (c *NodeConnection) IncWriteBytes(writen uint64) {
	c.mutex.Lock()
	c.writeBytes1m += writen
	c.writeBytes5m += writen
	c.writeBytes30m += writen
	c.mutex.Unlock()
}

func (c *NodeConnection) ResetReadWriteBytes1m() {
	c.mutex.Lock()
	c.periodStartTime1m = time.Now()
	c.readBytes1m = 0
	c.writeBytes1m = 0
	c.mutex.Unlock()
}

func (c *NodeConnection) ResetReadWriteBytes5m() {
	c.mutex.Lock()
	c.periodStartTime5m = time.Now()
	c.readBytes5m = 0
	c.writeBytes5m = 0
	c.mutex.Unlock()
}

func (c *NodeConnection) ResetReadWriteBytes30m() {
	c.mutex.Lock()
	c.periodStartTime30m = time.Now()
	c.readBytes30m = 0
	c.writeBytes30m = 0
	c.mutex.Unlock()
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
	c.writeBytes1m = 0
	c.writeBytes5m = 0
	c.writeBytes30m = 0
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
	writeBytes1m       uint64
	writeBytes5m       uint64
	writeBytes30m      uint64
	targetId           string
}

type TargetConnectionCopy struct {
	Timeout            uint32
	PeriodStartTime1m  time.Time
	PeriodStartTime5m  time.Time
	PeriodStartTime30m time.Time
	ReadBytes1m        uint64
	ReadBytes5m        uint64
	ReadBytes30m       uint64
	WriteBytes1m       uint64
	WriteBytes5m       uint64
	WriteBytes30m      uint64
	TargetId           string
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
	c.writeBytes1m = 0
	c.writeBytes5m = 0
	c.writeBytes30m = 0
	c.targetId = targetId
}

func (c *TargetConnection) SetKeepAlive() {
	c.mutex.Lock()
	_ = c.conn.SetKeepAlive(true)
	c.mutex.Unlock()
}

func (c *TargetConnection) GetConnection() *net.TCPConn {
	var conn *net.TCPConn
	c.mutex.Lock()
	conn = c.conn
	c.mutex.Unlock()
	return conn
}

func (c *TargetConnection) GetTimeOut() uint32 {
	var timeout uint32
	c.mutex.Lock()
	timeout = c.timeout
	c.mutex.Unlock()
	return timeout
}

func (c *TargetConnection) DumpToTargetConnectionCopy() TargetConnectionCopy {
	var connCopy TargetConnectionCopy
	c.mutex.Lock()
	connCopy.Timeout = c.timeout
	connCopy.PeriodStartTime1m = c.periodStartTime1m
	connCopy.PeriodStartTime5m = c.periodStartTime5m
	connCopy.PeriodStartTime30m = c.periodStartTime30m
	connCopy.ReadBytes1m = c.readBytes1m
	connCopy.ReadBytes5m = c.readBytes5m
	connCopy.ReadBytes30m = c.readBytes30m
	connCopy.WriteBytes1m = c.writeBytes1m
	connCopy.WriteBytes5m = c.writeBytes5m
	connCopy.WriteBytes30m = c.writeBytes30m
	connCopy.TargetId = c.targetId
	c.mutex.Unlock()
	return connCopy
}

func (c *TargetConnection) IncReadBytes(readn uint64) {
	c.mutex.Lock()
	c.readBytes1m += readn
	c.readBytes5m += readn
	c.readBytes30m += readn
	c.mutex.Unlock()
}

func (c *TargetConnection) IncWriteBytes(writen uint64) {
	c.mutex.Lock()
	c.writeBytes1m += writen
	c.writeBytes5m += writen
	c.writeBytes30m += writen
	c.mutex.Unlock()
}

func (c *TargetConnection) ResetReadWriteBytes1m() {
	c.mutex.Lock()
	c.periodStartTime1m = time.Now()
	c.readBytes1m = 0
	c.writeBytes1m = 0
	c.mutex.Unlock()
}

func (c *TargetConnection) ResetReadWriteBytes5m() {
	c.mutex.Lock()
	c.periodStartTime5m = time.Now()
	c.readBytes5m = 0
	c.writeBytes5m = 0
	c.mutex.Unlock()
}

func (c *TargetConnection) ResetReadWriteBytes30m() {
	c.mutex.Lock()
	c.periodStartTime30m = time.Now()
	c.readBytes30m = 0
	c.writeBytes30m = 0
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
	c.writeBytes1m = 0
	c.writeBytes5m = 0
	c.writeBytes30m = 0
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

func (l *LBConnectionPairMgr) GetTargetConnCountByTargetId(targetId string) uint32 {
	var count uint32 = 0
	l.mutex.RLock()
	for k, _ := range l.targetConnToNodeConnMap {
		if k.targetId == targetId {
			count++
		}
	}
	l.mutex.RUnlock()
	return count
}

func (l *LBConnectionPairMgr) GetTargetConnPairsByTargetId(targetId string) map[*TargetConnection]*NodeConnection {
	connPair := make(map[*TargetConnection]*NodeConnection)
	l.mutex.RLock()
	for k, v := range l.targetConnToNodeConnMap {
		if k.targetId == targetId {
			connPair[k] = v
		}
	}
	l.mutex.RUnlock()
	return connPair
}

func (l *LBConnectionPairMgr) GetAllTargetConnPairs() map[*TargetConnection]*NodeConnection {
	connPair := make(map[*TargetConnection]*NodeConnection)
	l.mutex.RLock()
	for k, v := range l.targetConnToNodeConnMap {
		connPair[k] = v
	}
	l.mutex.RUnlock()
	return connPair
}

func (l *LBConnectionPairMgr) ResetReadWriteBytes1m() {
	l.mutex.RLock()
	for k, v := range l.nodeConnToTargetConnMap {
		k.ResetReadWriteBytes1m()
		v.ResetReadWriteBytes1m()
	}
	l.mutex.RUnlock()
}

func (l *LBConnectionPairMgr) ResetReadWriteBytes5m() {
	l.mutex.RLock()
	for k, v := range l.nodeConnToTargetConnMap {
		k.ResetReadWriteBytes5m()
		v.ResetReadWriteBytes5m()
	}
	l.mutex.RUnlock()
}

func (l *LBConnectionPairMgr) ResetReadWriteBytes30m() {
	l.mutex.RLock()
	for k, v := range l.nodeConnToTargetConnMap {
		k.ResetReadWriteBytes30m()
		v.ResetReadWriteBytes30m()
	}
	l.mutex.RUnlock()
}

func (l *LBConnectionPairMgr) Destroy() {
	l.mutex = nil
	l.nodeConnToTargetConnMap = nil
	l.targetConnToNodeConnMap = nil
}

type ConnectionPairStatInfo struct {
	NodeConnFrom          string  `json:"nodeConnFrom"`
	NodeConnTo            string  `json:"nodeConnTo"`
	NodeReadSpeed1Min     float32 `json:"nodeReadSpeed1Min"`
	NodeWriteSpeed1Min    float32 `json:"nodeWriteSpeed1Min"`
	NodeReadSpeed5Min     float32 `json:"nodeReadSpeed5Min"`
	NodeWriteSpeed5Min    float32 `json:"nodeWriteSpeed5Min"`
	NodeReadSpeed30Min    float32 `json:"nodeReadSpeed30Min"`
	NodeWriteSpeed30Min   float32 `json:"nodeWriteSpeed30Min"`
	TargetConnFrom        string  `json:"targetConnFrom"`
	TargetConnTo          string  `json:"targetConnTo"`
	TargetReadSpeed1Min   float32 `json:"targetReadSpeed1Min"`
	TargetWriteSpeed1Min  float32 `json:"targetWriteSpeed1Min"`
	TargetReadSpeed5Min   float32 `json:"targetReadSpeed5Min"`
	TargetWriteSpeed5Min  float32 `json:"targetWriteSpeed5Min"`
	TargetReadSpeed30Min  float32 `json:"targetReadSpeed30Min"`
	TargetWriteSpeed30Min float32 `json:"targetWriteSpeed30Min"`
	TargetId              string  `json:"targetId"`
}

func getConnectionPairStatInfo(nodeConn *NodeConnection, targetConn *TargetConnection) *ConnectionPairStatInfo {
	statInfo := new(ConnectionPairStatInfo)

	nodeConnRaw := nodeConn.GetConnection()
	if nodeConnRaw == nil {
		return nil
	}
	targetConnRaw := targetConn.GetConnection()
	if targetConnRaw == nil {
		return nil
	}

	statInfo.NodeConnFrom = nodeConnRaw.RemoteAddr().String()
	statInfo.NodeConnTo = nodeConnRaw.LocalAddr().String()
	statInfo.NodeReadSpeed1Min = float32(nodeConn.readBytes1m) * 8 * 1e9 / float32(time.Now().Sub(nodeConn.periodStartTime1m).Nanoseconds())
	statInfo.NodeReadSpeed5Min = float32(nodeConn.readBytes5m) * 8 * 1e9 / float32(time.Now().Sub(nodeConn.periodStartTime1m).Nanoseconds())
	statInfo.NodeReadSpeed30Min = float32(nodeConn.readBytes30m) * 8 * 1e9 / float32(time.Now().Sub(nodeConn.periodStartTime1m).Nanoseconds())
	statInfo.NodeWriteSpeed1Min = float32(nodeConn.writeBytes1m) * 8 * 1e9 / float32(time.Now().Sub(nodeConn.periodStartTime1m).Nanoseconds())
	statInfo.NodeWriteSpeed5Min = float32(nodeConn.writeBytes5m) * 8 * 1e9 / float32(time.Now().Sub(nodeConn.periodStartTime1m).Nanoseconds())
	statInfo.NodeWriteSpeed30Min = float32(nodeConn.writeBytes30m) * 8 * 1e9 / float32(time.Now().Sub(nodeConn.periodStartTime1m).Nanoseconds())

	statInfo.TargetConnFrom = targetConnRaw.LocalAddr().String()
	statInfo.TargetConnTo = targetConnRaw.RemoteAddr().String()
	statInfo.TargetReadSpeed1Min = float32(targetConn.readBytes1m) * 8 * 1e9 / float32(time.Now().Sub(targetConn.periodStartTime1m).Nanoseconds())
	statInfo.TargetReadSpeed5Min = float32(targetConn.readBytes5m) * 8 * 1e9 / float32(time.Now().Sub(targetConn.periodStartTime1m).Nanoseconds())
	statInfo.TargetReadSpeed30Min = float32(targetConn.readBytes30m) * 8 * 1e9 / float32(time.Now().Sub(targetConn.periodStartTime1m).Nanoseconds())
	statInfo.TargetWriteSpeed1Min = float32(targetConn.writeBytes1m) * 8 * 1e9 / float32(time.Now().Sub(targetConn.periodStartTime1m).Nanoseconds())
	statInfo.TargetWriteSpeed5Min = float32(targetConn.writeBytes5m) * 8 * 1e9 / float32(time.Now().Sub(targetConn.periodStartTime1m).Nanoseconds())
	statInfo.TargetWriteSpeed30Min = float32(targetConn.writeBytes30m) * 8 * 1e9 / float32(time.Now().Sub(targetConn.periodStartTime1m).Nanoseconds())
	statInfo.TargetId = targetConn.targetId

	return statInfo
}

func startMaintainLoop(g goroutine_mgr.Goroutine) {
	defer g.OnQuit()

	maintainIntv := 60 * time.Second
	maintainTimer := time.NewTimer(maintainIntv)

	maintainIndex := uint64(0)

	for {
		select {
		case <-maintainTimer.C:
			if maintainIndex%1 == 0 {
				LBConnectionPairMgrP.ResetReadWriteBytes1m()
			}
			if maintainIndex%5 == 0 {
				LBConnectionPairMgrP.ResetReadWriteBytes5m()
			}
			if maintainIndex%30 == 0 {
				LBConnectionPairMgrP.ResetReadWriteBytes30m()
			}

			maintainIndex++
			maintainTimer.Reset(maintainIntv)
		}
	}
}
