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
	start(pars ...interface{}) error
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
				c.runChk()
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

	c.lk.Lock()
	defer c.lk.Unlock()
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
	// 目前只有timer需要自动Task
	err := comm.Db.Where("del!=1 and enable=1").And("types='timer'").Find(&ls)
	if err != nil {
		println("triggerManager err:" + err.Error())
		return
	}
	for _, v := range ls {
		c.lk.Lock()
		_, ok := c.tasks[v.Id]
		c.lk.Unlock()
		if !ok {
			c.StartOne(v)
		}
	}
}
func (c *triggerManager) Refresh(id int) {
	c.tmRfs = time.Time{}
	c.lk.Lock()
	defer c.lk.Unlock()
	v, ok := c.tasks[id]
	if ok {
		v.stop()
		delete(c.tasks, id)
	}
}

func (c *triggerManager) StartOne(trg *model.TTrigger, pars ...interface{}) {
	defer ruisUtil.Recovers("StartOne", nil)
	if trg.Del == 1 || trg.Enable != 1 {
		return
	}
	var i iTrigger
	switch trg.Types {
	case "timer":
		i = &trigTimeTask{tg: trg}
	case "hook":
		i = &trigHookTask{tg: trg}
	case "worked":
		i = &trigWorkedTask{tg: trg}
	}
	if i == nil {
		return
	}
	errs := i.start(pars...)
	trg.Errs = ""
	if errs != nil {
		trg.Errs = errs.Error()
		println("trigTimeTask start err:" + trg.Errs)
	} else {
		c.lk.Lock()
		c.tasks[trg.Id] = i
		c.lk.Unlock()
	}
	comm.Db.Cols("errs").Where("id=?", trg.Id).Update(trg)
}
