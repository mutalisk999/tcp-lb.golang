package main

import (
	"sync"
)

type LBNode struct {
	mutex          *sync.RWMutex
	endPointListen string
	status         uint8
	maxConnCount   uint32
	connCount      uint32
	timeoutRead    uint32
}

func (l *LBNode) Initialise(endPoint string, maxConn uint32, timeoutRead uint32) {
	l.mutex = new(sync.RWMutex)
	l.endPointListen = endPoint
	l.status = 0
	l.maxConnCount = maxConn
	l.connCount = 0
	l.timeoutRead = timeoutRead
}

func (l *LBNode) Destroy() {
	l.mutex = nil
	l.endPointListen = ""
	l.status = 0
	l.maxConnCount = 0
	l.connCount = 0
	l.timeoutRead = 0
}

type LBTarget struct {
	mutex        *sync.RWMutex
	endPointConn string
	status       uint8
	maxConnCount uint32
	connCount    uint32
	timeoutConn  uint32
	timeoutRead  uint32
}

func (l *LBTarget) Initialise(endPoint string, maxConn uint32, timeoutConn uint32, timeoutRead uint32) {
	l.mutex = new(sync.RWMutex)
	l.endPointConn = endPoint
	l.status = 0
	l.maxConnCount = maxConn
	l.connCount = 0
	l.timeoutConn = timeoutConn
	l.timeoutRead = timeoutRead
}

func (l *LBTarget) Destroy() {
	l.mutex = nil
	l.endPointConn = ""
	l.status = 0
	l.maxConnCount = 0
	l.connCount = 0
	l.timeoutConn = 0
	l.timeoutRead = 0
}

var LBNodeP LBNode
var LBTargetMap sync.Map
