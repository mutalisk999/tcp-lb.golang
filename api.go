package main

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"net/http"
)

type ConnectionPairInfo struct {
	NodeConnFrom              string
	NodeConnTo                string
	NodeReadSpeedPeriod1m     float64
	NodeWriteSpeedPeriod1m    float64
	NodeReadSpeedPeriod5m     float64
	NodeWriteSpeedPeriod5m    float64
	NodeReadSpeedPeriod30m    float64
	NodeWriteSpeedPeriod30m   float64
	TargetConnFrom            string
	TargetConnTo              string
	TargetReadSpeedPeriod1m   float64
	TargetWriteSpeedPeriod1m  float64
	TargetReadSpeedPeriod5m   float64
	TargetWriteSpeedPeriod5m  float64
	TargetReadSpeedPeriod30m  float64
	TargetWriteSpeedPeriod30m float64
	TargetId                  string
}

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
	targetId := CalcTargetId(*args)
	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget == nil {
		return errors.New("invalid args, can not find this targetId")
	}
	active := lbTarget.GetActive()
	if active {
		return nil
	} else {
		lbTarget.SetActive(true)
		return nil
	}
}

func (s *Service) SetTargetDisable(r *http.Request, args *string, reply *interface{}) error {
	targetId := CalcTargetId(*args)
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
		return nil
	}
}

func (s *Service) SetTargetInfo(r *http.Request, args *LBTargetCopy, reply *interface{}) error {
	targetId := CalcTargetId(args.EndPointConn)
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
	return nil
}

func (s *Service) AddTargetInfo(r *http.Request, args *LBTargetCopy, reply *interface{}) error {
	targetId := CalcTargetId(args.EndPointConn)
	targetEndPoint := args.EndPointConn
	targetActive := args.Active
	targetMaxConn := args.MaxConnCount
	targetTimeOut := args.Timeout

	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget != nil {
		return errors.New("invalid args, target endpoint has already existed")
	}

	lbTarget = new(LBTarget)
	lbTarget.Initialise(targetEndPoint, targetActive, targetMaxConn, targetTimeOut)

	LBTargetsMgrP.Set(targetId, lbTarget)
	return nil
}

func (s *Service) GetTargetConnectPairsInfo(r *http.Request, args *string, reply *interface{}) error {
	targetId := CalcTargetId(*args)
	lbTarget := LBTargetsMgrP.Get(targetId)
	if lbTarget == nil {
		return errors.New("invalid args, can not find this target endpoint")
	}

	connPair := LBConnectionPairMgrP.GetTargetConnPairsByTargetId(targetId)
	for k, v := range connPair {

	}
	return nil
}

func (s *Service) GetAllConnectPairsInfo(r *http.Request, args *interface{}, reply *interface{}) error {
	connPair := LBConnectionPairMgrP.GetAllTargetConnPairs()
	for k, v := range connPair {

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
