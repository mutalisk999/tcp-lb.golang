package main

import (
	"fmt"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"sync"
	"time"
)

type LBNode struct {
	mutex          *sync.RWMutex
	endPointListen string
	maxConnCount   uint32
	connCount      uint32
	timeout        uint32
}

func (l *LBNode) Initialise(endPoint string, maxConn uint32, timeout uint32) {
	l.mutex = new(sync.RWMutex)
	l.endPointListen = endPoint
	l.maxConnCount = maxConn
	l.connCount = 0
	l.timeout = timeout
}

func (l *LBNode) GetConnCount() uint32 {
	var count uint32
	l.mutex.RLock()
	count = l.connCount
	l.mutex.RUnlock()
	return count
}

func (l *LBNode) GetMaxConnCount() uint32 {
	var maxCount uint32
	l.mutex.RLock()
	maxCount = l.maxConnCount
	l.mutex.RUnlock()
	return maxCount
}

func (l *LBNode) GetConnInfoStr() string {
	var str string
	l.mutex.RLock()
	str = fmt.Sprintf("[%d/%d]", l.connCount, l.maxConnCount)
	l.mutex.RUnlock()
	return str
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
	l.timeout = 0
}

type LBTarget struct {
	mutex        *sync.RWMutex
	endPointConn string
	status       uint8
	maxConnCount uint32
	connCount    uint32
	timeout      uint32
}

func (l *LBTarget) Initialise(endPoint string, maxConn uint32, timeout uint32) {
	l.mutex = new(sync.RWMutex)
	l.endPointConn = endPoint
	l.status = 0
	l.maxConnCount = maxConn
	l.connCount = 0
	l.timeout = timeout
}

func (l *LBTarget) GetConnCount() uint32 {
	var count uint32
	l.mutex.RLock()
	count = l.connCount
	l.mutex.RUnlock()
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
	l.mutex.RLock()
	targetCopy.EndPointConn = l.endPointConn
	targetCopy.MaxConnCount = l.maxConnCount
	targetCopy.ConnCount = l.connCount
	targetCopy.Timeout = l.timeout
	l.mutex.RUnlock()
	return targetCopy
}

func (l *LBTarget) Destroy() {
	l.mutex = nil
	l.endPointConn = ""
	l.status = 0
	l.maxConnCount = 0
	l.connCount = 0
	l.timeout = 0
}

type LBTargetCopy struct {
	EndPointConn string
	Status       uint8
	MaxConnCount uint32
	ConnCount    uint32
	Timeout      uint32
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
	l.mutex.RLock()
	v, ok := l.targetsMap[targetId]
	l.mutex.RUnlock()
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
	l.mutex.RLock()
	for _, v := range l.targetsMap {
		lbTargetsCopy = append(lbTargetsCopy, v.DumpToLBTargetCopy())
	}
	l.mutex.RUnlock()
	return lbTargetsCopy
}

func (l *LBTargetsMgr) GetTargetsCount() int {
	var count int
	l.mutex.RLock()
	count = len(l.targetsMap)
	l.mutex.RUnlock()
	return count
}

func (l *LBTargetsMgr) Destroy() {
	l.mutex = nil
	l.targetsMap = nil
}

//for sort and sort by ConnCount
type LBTargetCopys []LBTargetCopy

func (s LBTargetCopys) Len() int           { return len(s) }
func (s LBTargetCopys) Less(i, j int) bool { return s[i].ConnCount < s[j].ConnCount }
func (s LBTargetCopys) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func handleNodeData(g goroutine_mgr.Goroutine, a interface{}, b interface{}) {
	defer g.OnQuit()

	c := a.(*NodeConnection)
	ct := b.(*TargetConnection)

	conn := c.GetConnection()
	timeout := c.GetTimeOut()

	for {
		_ = conn.SetDeadline(time.Now().Add(time.Duration(int64(timeout) * 1000 * 1000 * 1000)))

	}

}

func handleTargetData(g goroutine_mgr.Goroutine, a interface{}, b interface{}) {
	defer g.OnQuit()

}
