package mgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gokins/comm"
	"gokins/model"
	"gokins/service/dbService"
	"time"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

type confTimeBean struct {
	Mid      int       `json:"mid"`
	Repeated string    `json:"repeated"`
	Dates    string    `json:"dates"`
	Date     time.Time `json:"-"`
}
type trigTimeTask struct {
	md    *model.TModel
	tg    *model.TTrigger
	conf  *confTimeBean
	ctx   context.Context
	cncl  context.CancelFunc
	runtm time.Time
}

func (c *trigTimeTask) start() error {
	if c.tg == nil || c.cncl != nil {
		return errors.New("already run")
	}
	c.conf = &confTimeBean{}
	err := json.Unmarshal([]byte(c.tg.Config), c.conf)
	if err != nil {
		return err
	}
	tms, err := time.Parse(comm.TimeFmtpck, c.conf.Dates)
	if err != nil {
		return err
	}
	c.conf.Date = tms.Local()
	println(fmt.Sprintf("%d-%d-%d %d:%d:%d", c.conf.Date.Year(), c.conf.Date.Month(), c.conf.Date.Day(), c.conf.Date.Hour(), c.conf.Date.Minute(), c.conf.Date.Second()))
	c.md = dbService.GetModel(c.conf.Mid)
	if c.md == nil {
		return errors.New("not found model")
	}
	c.ctx, c.cncl = context.WithCancel(mgrCtx)
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				goto end
			default:
				c.run()
				time.Sleep(time.Millisecond * 200)
			}
		}
	end:
		println("ctx end!")
	}()
	return nil
}
func (c *trigTimeTask) stop() {
	if c.cncl != nil {
		c.cncl()
		c.cncl = nil
	}
}
func (c *trigTimeTask) isRun() bool {
	return c.cncl == nil
}
func (c *trigTimeTask) run() {
	defer ruisUtil.Recovers("RunTask start", func(errs string) {
		println("trigTimeTask run err:" + errs)
	})
	if time.Since(c.runtm).Seconds() < 5 {
		return
	}

	isend := false
	match := false
	switch c.conf.Repeated {
	case "1":
		match = c.check(0, 0, 0, 1, 1, 1, 0)
	case "2":
		match = c.check(0, 0, 0, 1, 1, 1, 1)
	case "3":
		match = c.check(0, 0, 1, 1, 1, 1, 0)
	case "4":
		match = c.check(0, 1, 1, 1, 1, 1, 0)
	default:
		isend = true
		match = c.check(1, 1, 1, 1, 1, 1, 0)
	}

	if match {
		if isend {
			defer c.end()
		}
		c.runtm = time.Now()
		rn := &model.TModelRun{}
		rn.Tid = c.md.Id
		rn.Uid = c.tg.Uid
		rn.Times = time.Now()
		rn.Tgtyp = "timer"
		comm.Db.Insert(rn)
		ExecMgr.Refresh()
		println("trigTimeTask model run id:", rn.Id)
	}
}
func (c *trigTimeTask) check(y, m, d, h, min, s, w int) bool {
	now := time.Now()
	if y == 1 && now.Year() != c.conf.Date.Year() {
		return false
	}
	if m == 1 && now.Month() != c.conf.Date.Month() {
		return false
	}
	if d == 1 && now.Day() != c.conf.Date.Day() {
		return false
	}
	if h == 1 && now.Hour() != c.conf.Date.Hour() {
		return false
	}
	if min == 1 && now.Minute() != c.conf.Date.Minute() {
		return false
	}
	if s == 1 && now.Second() != c.conf.Date.Second() {
		return false
	}
	if w == 1 && now.Weekday() != c.conf.Date.Weekday() {
		return false
	}
	return true
}
func (c *trigTimeTask) end() {
	c.tg.Enable = 2
	comm.Db.Cols("enable").Where("id=?", c.tg.Id).Update(c.tg)
	c.stop()
}
