package main

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"net/http"
)

type Service struct {
}

func (s *Service) GetNodeInfo(r *http.Request, args *interface{}, reply *LBNodeCopy) error {
	*reply = LBNodeP.DumpToLBNodeCopy()
	return nil
}

func (s *Service) GetTargetsInfo(r *http.Request, args *interface{}, reply *[]LBTargetCopy) error {
	*reply = LBTargetsMgrP.DumpTargetsCopy()
	return nil
}

func (s *Service) SetTargetEnable(r *http.Request, args *string, reply *interface{}) error {
	targetId := calcTargetId(*args)
	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget == nil {
		return errors.New("invalid args, can not find this targetId")
	}
	active := lbTarget.GetActive()
	if active {
		return nil
	} else {
		lbTarget.SetActive(true)

		LBConfigMutex.Lock()
		for idx, _ := range LBConfig.Targets {
			if LBConfig.Targets[idx].ConnEndPoint == *args {
				LBConfig.Targets[idx].Active = true
			}
		}
		saveConfig(&LBConfig)
		LBConfigMutex.Unlock()

		return nil
	}
}

func (s *Service) SetTargetDisable(r *http.Request, args *string, reply *interface{}) error {
	targetId := calcTargetId(*args)
	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget == nil {
		return errors.New("invalid args, can not find this target endpoint")
	}
	active := lbTarget.GetActive()
	if !active {
		return nil
	} else {
		lbTarget.SetActive(false)
		connPair := LBConnectionPairMgrP.GetTargetConnPairsByTargetId(targetId)
		for k, v := range connPair {
			LBConnectionPairMgrP.RemoveByTargetConn(k)
			targetConn := k.GetConnection()
			if targetConn != nil {
				_ = targetConn.Close()
			}
			nodeConn := v.GetConnection()
			if nodeConn != nil {
				_ = nodeConn.Close()
			}
		}

		LBConfigMutex.Lock()
		for idx, _ := range LBConfig.Targets {
			if LBConfig.Targets[idx].ConnEndPoint == *args {
				LBConfig.Targets[idx].Active = false
			}
		}
		saveConfig(&LBConfig)
		LBConfigMutex.Unlock()

		return nil
	}
}

func (s *Service) SetTargetInfo(r *http.Request, args *LBTargetCopy, reply *interface{}) error {
	targetId := calcTargetId(args.EndPointConn)
	targetEndPoint := args.EndPointConn
	targetActive := args.Active
	targetMaxConn := args.MaxConnCount
	targetTimeOut := args.Timeout

	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget == nil {
		return errors.New("invalid args, can not find this target endpoint")
	}
	targetCopy := lbTarget.DumpToLBTargetCopy()

	if targetEndPoint == targetCopy.EndPointConn && targetActive == targetCopy.Active &&
		targetMaxConn == targetCopy.MaxConnCount && targetTimeOut == targetCopy.Timeout {
		return nil
	}

	lbTarget.SetActive(false)

	if targetCopy.Active {
		connPair := LBConnectionPairMgrP.GetTargetConnPairsByTargetId(targetId)
		for k, v := range connPair {
			LBConnectionPairMgrP.RemoveByTargetConn(k)
			targetConn := k.GetConnection()
			if targetConn != nil {
				_ = targetConn.Close()
			}
			nodeConn := v.GetConnection()
			if nodeConn != nil {
				_ = nodeConn.Close()
			}
		}
	}
	lbTarget.Update(targetEndPoint, targetActive, targetMaxConn, targetTimeOut)

	LBConfigMutex.Lock()
	for idx, _ := range LBConfig.Targets {
		if LBConfig.Targets[idx].ConnEndPoint == args.EndPointConn {
			LBConfig.Targets[idx].ConnEndPoint = targetEndPoint // not change
			LBConfig.Targets[idx].Active = targetActive
			LBConfig.Targets[idx].MaxConn = targetMaxConn
			LBConfig.Targets[idx].Timeout = targetTimeOut
		}
	}
	saveConfig(&LBConfig)
	LBConfigMutex.Unlock()

	return nil
}

func (s *Service) AddTargetInfo(r *http.Request, args *LBTargetCopy, reply *interface{}) error {
	targetId := calcTargetId(args.EndPointConn)
	targetEndPoint := args.EndPointConn
	targetActive := args.Active
	targetMaxConn := args.MaxConnCount
	targetTimeOut := args.Timeout

	if !verifyEndPointStr(targetEndPoint) {
		return errors.New("invalid args: wrong target endpoint")
	}

	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget != nil {
		return errors.New("invalid args, target endpoint has already existed")
	}

	lbTarget = new(LBTarget)
	lbTarget.Initialise(targetEndPoint, targetActive, targetMaxConn, targetTimeOut)

	LBTargetsMgrP.Set(targetId, lbTarget)

	LBConfigMutex.Lock()
	LBConfig.Targets = append(LBConfig.Targets, TargetConfig{targetEndPoint, targetMaxConn,
		targetTimeOut, targetActive})
	saveConfig(&LBConfig)
	LBConfigMutex.Unlock()

	return nil
}

func (s *Service) DelTargetInfo(r *http.Request, args *string, reply *interface{}) error {
	targetId := calcTargetId(*args)
	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget == nil {
		return errors.New("invalid args, can not find this target endpoint")
	}
	LBTargetsMgrP.Delete(targetId)

	active := lbTarget.GetActive()
	if !active {
		return nil
	} else {
		connPair := LBConnectionPairMgrP.GetTargetConnPairsByTargetId(targetId)
		for k, v := range connPair {
			LBConnectionPairMgrP.RemoveByTargetConn(k)
			targetConn := k.GetConnection()
			if targetConn != nil {
				_ = targetConn.Close()
			}
			nodeConn := v.GetConnection()
			if nodeConn != nil {
				_ = nodeConn.Close()
			}
		}
	}

	LBConfigMutex.Lock()
	var idxFound = -1
	for idx, _ := range LBConfig.Targets {
		if LBConfig.Targets[idx].ConnEndPoint == *args {
			idxFound = idx
			break
		}
	}
	if idxFound != -1 {
		LBConfig.Targets = append(LBConfig.Targets[0:idxFound], LBConfig.Targets[idxFound+1:]...)
		saveConfig(&LBConfig)
	}
	LBConfigMutex.Unlock()

	return nil
}

func (s *Service) GetTargetConnectPairsInfo(r *http.Request, args *string, reply *[]ConnectionPairStatInfo) error {
	targetId := calcTargetId(*args)
	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget == nil {
		return errors.New("invalid args, can not find this target endpoint")
	}

	connPair := LBConnectionPairMgrP.GetTargetConnPairsByTargetId(targetId)
	for k, v := range connPair {
		statInfo := getConnectionPairStatInfo(v, k)
		if statInfo != nil {
			*reply = append(*reply, *statInfo)
		}
	}
	return nil
}

func (s *Service) GetAllConnectPairsInfo(r *http.Request, args *interface{}, reply *[]ConnectionPairStatInfo) error {
	connPair := LBConnectionPairMgrP.GetAllTargetConnPairs()
	for k, v := range connPair {
		statInfo := getConnectionPairStatInfo(v, k)
		if statInfo != nil {
			*reply = append(*reply, *statInfo)
		}
	}
	return nil
}

func startApiServer(g goroutine_mgr.Goroutine, a interface{}) {
	defer g.OnQuit()

	c := a.(*Config)

	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	rpcService := new(Service)
	_ = rpcServer.RegisterService(rpcService, "")

	urlRouter := mux.NewRouter()
	urlRouter.Handle("/api", rpcServer)

	Info.Printf("Api listening on %s", c.Api.ListenEndPoint)
	_ = http.ListenAndServe(c.Api.ListenEndPoint, urlRouter)
}
