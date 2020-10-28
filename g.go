package main

import (
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"sync"
)

var LBConfig Config
var LBConfigMutex sync.Mutex

var LBNodeP *LBNode
var LBTargetsMgrP *LBTargetsMgr
var LBConnectionPairMgrP *LBConnectionPairMgr

var LBGoroutineManagerP *goroutine_mgr.GoroutineManager
