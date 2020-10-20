package mgr

import (
	"gokins/comm"
	"gokins/model"
	"sync"
	"time"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

var TriggerMgr = &triggerManager{}

type triggerManager struct {
	tmChk time.Time
	tmRfs time.Time
	lk    sync.Mutex
	tasks map[int]iTrigger
}

type iTrigger interface {
	start() error
	stop()
	isRun() bool
}

func (c *triggerManager) Start() {
	c.tasks = make(map[int]iTrigger)
	go func() {
		for {
			select {
			case <-mgrCtx.Done():
				goto end
			default:
				c.runRfs()
				time.Sleep(time.Second)
			}
		}
	end:
		println("ctx end!")
	}()
}
func (c *triggerManager) runChk() {
	defer ruisUtil.Recovers("triggerManager runChk", nil)
	if time.Since(c.tmChk).Seconds() < 30 {
		return
	}
	c.tmChk = time.Now()

	for k, v := range c.tasks {
		if !v.isRun() {
			delete(c.tasks, k)
		}
	}
}
func (c *triggerManager) runRfs() {
	defer ruisUtil.Recovers("run", func(errs string) {
		println("triggerManager run err:" + errs)
	})
	if time.Since(c.tmRfs).Minutes() < 30 {
		return
	}
	c.tmRfs = time.Now()
	var ls []*model.TTrigger
	err := comm.Db.Where("del!=1 and enable=1").Find(&ls)
	if err != nil {
		println("triggerManager err:" + err.Error())
		return
	}
	c.lk.Lock()
	defer c.lk.Unlock()
	for _, v := range ls {
		_, ok := c.tasks[v.Id]
		if !ok {
			var i iTrigger
			switch v.Types {
			case "git":
				//TODO: git not get
				break
			case "timer":
				i = &trigTimeTask{tg: v}
			}
			if i == nil {
				continue
			}
			errs := i.start()
			v.Errs = ""
			if errs != nil {
				v.Errs = errs.Error()
				println("trigTimeTask start err:" + v.Errs)
			} else {
				c.tasks[v.Id] = i
			}
			comm.Db.Cols("errs").Where("id=?", v.Id).Update(v)
		}
	}
}
func (c *triggerManager) Refresh(id int) {
	c.tmRfs = time.Time{}
	c.lk.Lock()
	defer c.lk.Unlock()
	v, ok := c.tasks[id]
	delete(c.tasks, id)
	if ok {
		v.stop()
	}
}
