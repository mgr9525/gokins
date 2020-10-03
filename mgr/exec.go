package mgr

import (
	"gokins/comm"
	"gokins/model"
	"sync"
	"time"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

var ExecMgr = &execManager{}

type execManager struct {
	tm    time.Time
	lk    sync.Mutex
	tasks map[int]*RunTask
}

func (c *execManager) Start() {
	c.tasks = make(map[int]*RunTask)
	go func() {
		for {
			select {
			case <-mgrCtx.Done():
				break
			default:
				c.run()
				time.Sleep(time.Second)
			}
		}
	}()
}
func (c *execManager) run() {
	defer ruisUtil.Recovers("run", func(errs string) {
		println("execManager run err:" + errs)
	})
	if time.Since(c.tm).Seconds() < 20 {
		return
	}
	c.tm = time.Now()

	c.lk.Lock()
	defer c.lk.Unlock()
	for k, v := range c.tasks {
		if v.cncl == nil {
			delete(c.tasks, k)
		}
	}

	if len(c.tasks) >= comm.RunTaskLen {
		return
	}
	var ls []*model.TModelRun
	err := comm.Db.Where("state=0 or state=1").Find(&ls)
	if err != nil {
		println("execManager err:" + err.Error())
		return
	}
	for _, v := range ls {
		if v.State == 0 {
			v.State = 1
			comm.Db.Cols("state").Where("id=?", v.Id).Update(v)
		}
		_, ok := c.tasks[v.Id]
		if !ok {
			e := &RunTask{Mr: v}
			c.tasks[v.Id] = e
			e.start()
		}
	}
}
func (c *execManager) Refresh() {
	c.tm = time.Time{}
}
func (c *execManager) StopTask(id int) {
	c.lk.Lock()
	defer c.lk.Unlock()
	e, ok := c.tasks[id]
	if ok {
		e.stop()
	}
	//v := &model.TModelRun{}
	//v.State = -1
	//_, err := comm.Db.Cols("state").Where("id=?", v.Id).Update(v)
	//return err
}
func (c *execManager) TaskRead(id, pid int) string {
	c.lk.Lock()
	defer c.lk.Unlock()
	e, ok := c.tasks[id]
	if ok {
		return e.read(pid)
	}
	return ""
}
