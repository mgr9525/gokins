package mgr

import (
	"context"
	"encoding/json"
	"errors"
	"gokins/comm"
	"gokins/model"
	"gokins/service/dbService"
	"time"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

type confWorkedBean struct {
	Forced bool `json:"forced"`
}
type trigWorkedTask struct {
	md   *model.TModel
	tg   *model.TTrigger
	conf *confWorkedBean
	ctx  context.Context
	cncl context.CancelFunc

	mde *model.TModel
	mr  *model.TModelRun
}

func (c *trigWorkedTask) start(pars ...interface{}) error {
	if c.tg == nil || c.cncl != nil {
		return errors.New("already run")
	}
	if len(pars) < 2 {
		return errors.New("param err")
	}
	c.mde = pars[0].(*model.TModel)
	c.mr = pars[1].(*model.TModelRun)
	c.conf = &confWorkedBean{}
	err := json.Unmarshal([]byte(c.tg.Config), c.conf)
	if err != nil {
		return err
	}
	c.md = dbService.GetModel(c.tg.Mid)
	if c.md == nil {
		return errors.New("not found model")
	}
	c.ctx, c.cncl = context.WithCancel(mgrCtx)
	go func() {
		c.tg.Errs = ""
		err := c.run()
		if err != nil {
			c.tg.Errs = err.Error()
		}
		comm.Db.Cols("errs").Where("id=?", c.tg.Id).Update(c.tg)
		c.stop()
		println("ctx end!")
	}()
	return nil
}
func (c *trigWorkedTask) stop() {
	if c.cncl != nil {
		c.cncl()
		c.cncl = nil
	}
}
func (c *trigWorkedTask) isRun() bool {
	return c.cncl == nil
}
func (c *trigWorkedTask) run() error {
	defer ruisUtil.Recovers("RunTask start", func(errs string) {
		println("trigWorkedTask run err:" + errs)
	})

	var err error
	if c.conf.Forced || c.mr.State == 4 {
		rn := &model.TModelRun{}
		rn.Tid = c.md.Id
		rn.Uid = c.tg.Uid
		rn.Times = time.Now()
		rn.Tgid = c.tg.Id
		rn.Tgtyps = "流水线触发"
		_, err = comm.Db.Insert(rn)
		ExecMgr.Refresh()
		println("trigWorkedTask model run id:", rn.Id)
	}
	return err
}
