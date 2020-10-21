package main

import (
	"sync"
)

type LBNode struct {
	mutex          *sync.RWMutex
	endPointListen string
	maxConnCount   uint32
	connCount      uint32
	timeoutRead    uint32
}

func (l *LBNode) Initialise(endPoint string, maxConn uint32, timeoutRead uint32) {
	l.mutex = new(sync.RWMutex)
	l.endPointListen = endPoint
	l.maxConnCount = maxConn
	l.connCount = 0
	l.timeoutRead = timeoutRead
}

func (l *LBNode) GetConnCount() uint32 {
	var count uint32
	l.mutex.Lock()
	count = l.connCount
	l.mutex.Unlock()
	return count
}

func (l *LBNode) GetMaxConnCount() uint32 {
	var maxCount uint32
	l.mutex.Lock()
	maxCount = l.maxConnCount
	l.mutex.Unlock()
	return maxCount
}

func (l *LBNode) IncConnCount() {
	l.mutex.Lock()
	l.connCount++
	l.mutex.Unlock()
}

func (l *LBNode) DecConnCount() {
	l.mutex.Lock()
	l.connCount--
	l.mutex.Unlock()
}

func (l *LBNode) Destroy() {
	l.mutex = nil
	l.endPointListen = ""
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

func (l *LBTarget) GetConnCount() uint32 {
	var count uint32
	l.mutex.Lock()
	count = l.connCount
	l.mutex.Unlock()
	return count
}

func (l *LBTarget) IncConnCount() {
	l.mutex.Lock()
	l.connCount++
	l.mutex.Unlock()
}

func (l *LBTarget) DecConnCount() {
	l.mutex.Lock()
	l.connCount--
	l.mutex.Unlock()
}

func (l *LBTarget) DumpToLBTargetCopy() LBTargetCopy {
	var targetCopy LBTargetCopy
	l.mutex.Lock()
	targetCopy.EndPointConn = l.endPointConn
	targetCopy.MaxConnCount = l.maxConnCount
	targetCopy.ConnCount = l.connCount
	targetCopy.TimeoutConn = l.timeoutConn
	targetCopy.TimeoutRead = l.timeoutRead
	l.mutex.Unlock()
	return targetCopy
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

type LBTargetCopy struct {
	EndPointConn string
	Status       uint8
	MaxConnCount uint32
	ConnCount    uint32
	TimeoutConn  uint32
	TimeoutRead  uint32
}

type LBTargetsMgr struct {
	mutex      *sync.RWMutex
	targetsMap map[string]*LBTarget
}

func (l *LBTargetsMgr) Initialise() {
	l.mutex = new(sync.RWMutex)
	l.targetsMap = make(map[string]*LBTarget)
}

func (l *LBTargetsMgr) Get(targetId string) *LBTarget {
	l.mutex.Lock()
	v, ok := l.targetsMap[targetId]
	l.mutex.Unlock()
	if !ok {
		return nil
	}
	return v
}

func (l *LBTargetsMgr) Delete(targetId string) {
	l.mutex.Lock()
	delete(l.targetsMap, targetId)
	l.mutex.Unlock()
}

func (l *LBTargetsMgr) Set(targetId string, target *LBTarget) {
	l.mutex.Lock()
	if target == nil {
		delete(l.targetsMap, targetId)
	} else {
		l.targetsMap[targetId] = target
	}
	l.mutex.Unlock()
}

func (l *LBTargetsMgr) DumpTargetsCopy() []LBTargetCopy {
	var lbTargetsCopy []LBTargetCopy
	l.mutex.Lock()
	for _, v := range l.targetsMap {
		lbTargetsCopy = append(lbTargetsCopy, v.DumpToLBTargetCopy())
	}
	l.mutex.Unlock()
	return lbTargetsCopy
}

func (l *LBTargetsMgr) GetTargetsCount() int {
	var count int
	l.mutex.Lock()
	count = len(l.targetsMap)
	l.mutex.Unlock()
	return count
}

func (l *LBTargetsMgr) Destroy() {
	l.mutex = nil
	l.targetsMap = nil
}
