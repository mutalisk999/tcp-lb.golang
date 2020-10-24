package main

import (
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"io"
	"sync"
	"time"
)

type LBNode struct {
	mutex          *sync.RWMutex
	endPointListen string
	maxConnCount   uint32
	timeout        uint32
	accept         chan int
}

func (l *LBNode) Initialise(endPoint string, maxConn uint32, timeout uint32) {
	l.mutex = new(sync.RWMutex)
	l.endPointListen = endPoint
	l.maxConnCount = maxConn
	l.timeout = timeout
	l.accept = make(chan int, l.maxConnCount)
}

func (l *LBNode) GetMaxConnCount() uint32 {
	var maxCount uint32
	l.mutex.RLock()
	maxCount = l.maxConnCount
	l.mutex.RUnlock()
	return maxCount
}

func (l *LBNode) ProductNewConn() {
	l.accept <- 0
}

func (l *LBNode) ComsumeNewConn() {
	<-l.accept
}

func (l *LBNode) Destroy() {
	l.mutex = nil
	l.endPointListen = ""
	l.maxConnCount = 0
	l.timeout = 0
	l.accept = nil
}

func (l *LBNode) GetConnCount() uint32 {
	return uint32(LBConnectionPairMgrP.GetNode2TargetPairCount())
}

type LBTarget struct {
	mutex        *sync.RWMutex
	endPointConn string
	status       uint8
	maxConnCount uint32
	timeout      uint32
}

func (l *LBTarget) Initialise(endPoint string, maxConn uint32, timeout uint32) {
	l.mutex = new(sync.RWMutex)
	l.endPointConn = endPoint
	l.status = 0
	l.maxConnCount = maxConn
	l.timeout = timeout
}

func (l *LBTarget) DumpToLBTargetCopy() LBTargetCopy {
	var targetCopy LBTargetCopy
	l.mutex.RLock()
	targetCopy.EndPointConn = l.endPointConn
	targetCopy.Status = l.status
	targetCopy.MaxConnCount = l.maxConnCount
	targetId := CaclTargetId(targetCopy.EndPointConn)
	targetCopy.ConnCount = LBConnectionPairMgrP.GetTargetConnCountByTargetId(targetId)
	targetCopy.Timeout = l.timeout
	l.mutex.RUnlock()
	return targetCopy
}

func (l *LBTarget) Destroy() {
	l.mutex = nil
	l.endPointConn = ""
	l.status = 0
	l.maxConnCount = 0
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

	connTarget := ct.GetConnection()
	timeoutTarget := c.GetTimeOut()

	for {
		var buf [4096]byte

		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		n, err := conn.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				Info.Printf("Closed by remote: %s", conn.RemoteAddr().String())
			} else {
				Error.Printf("Read from %s error: %s", conn.RemoteAddr().String(), err.Error())
			}
			break
		}
		c.IncReadBytes(uint64(n))

		_ = connTarget.SetWriteDeadline(time.Now().Add(time.Duration(timeoutTarget) * time.Second))
		s, err := connTarget.Write(buf[0:n])
		if err != nil {
			Error.Printf("Write to %s error: %s", connTarget.RemoteAddr().String(), err.Error())
			break
		}
		ct.IncWriteBytes(uint64(s))
	}

	LBConnectionPairMgrP.RemoveByNodeConn(c)

	if conn != nil {
		_ = conn.Close()
	}

	if connTarget != nil {
		_ = connTarget.Close()
	}
}

func handleTargetData(g goroutine_mgr.Goroutine, a interface{}, b interface{}) {
	defer g.OnQuit()

	c := a.(*NodeConnection)
	ct := b.(*TargetConnection)

	conn := c.GetConnection()
	timeout := c.GetTimeOut()

	connTarget := ct.GetConnection()
	timeoutTarget := c.GetTimeOut()

	for {
		var buf [4096]byte

		_ = connTarget.SetReadDeadline(time.Now().Add(time.Duration(timeoutTarget) * time.Second))
		n, err := connTarget.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				Info.Printf("Closed by remote: %s", connTarget.RemoteAddr().String())
			} else {
				Error.Printf("Read from %s error: %s", connTarget.RemoteAddr().String(), err.Error())
			}
			break
		}
		ct.IncReadBytes(uint64(n))

		_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		s, err := conn.Write(buf[0:n])
		if err != nil {
			Error.Printf("Write to %s error: %s", conn.RemoteAddr().String(), err.Error())
			break
		}
		c.IncWriteBytes(uint64(s))
	}

	LBConnectionPairMgrP.RemoveByTargetConn(ct)

	if conn != nil {
		_ = conn.Close()
	}

	if connTarget != nil {
		_ = connTarget.Close()
	}
}
