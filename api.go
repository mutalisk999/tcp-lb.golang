package main

import "github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"

func startApiServer(g goroutine_mgr.Goroutine, a interface{}) {
	defer g.OnQuit()

	c := a.(*Config)

	_ = c
}
