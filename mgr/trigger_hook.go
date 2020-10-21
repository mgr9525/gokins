package mgr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gokins/comm"
	"gokins/model"
	"gokins/service/dbService"
	"net/http"
	"time"

	"github.com/dop251/goja"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

type confHookBean struct {
	Mid  int    `json:"mid"`
	Plug string `json:"plug"`
	Conf string `json:"conf"`
}
type trigHookTask struct {
	md   *model.TModel
	tg   *model.TTrigger
	conf *confHookBean
	ctx  context.Context
	cncl context.CancelFunc

	js string
	vm *goja.Runtime

	querys string
	header http.Header
	bodys  []byte

	body  *ruisUtil.Map
	confs *ruisUtil.Map
}

func (c *trigHookTask) stop() {
	if c.cncl != nil {
		c.cncl()
		c.cncl = nil
	}
}
func (c *trigHookTask) isRun() bool {
	return c.cncl == nil
}
func (c *trigHookTask) start(pars ...interface{}) error {
	if len(pars) < 3 {
		return errors.New("param err")
	}
	c.querys = pars[0].(string)
	c.header = pars[1].(http.Header)
	c.bodys = pars[2].([]byte)
	c.conf = &confHookBean{}
	err := json.Unmarshal([]byte(c.tg.Config), c.conf)
	if err != nil {
		return err
	}
	c.confs = ruisUtil.NewMapo(c.conf.Conf)
	c.md = dbService.GetModel(c.conf.Mid)
	if c.md == nil {
		return errors.New("not found model")
	}
	js, ok := hookjsMap[c.conf.Plug]
	if !ok {
		return errors.New("not found plugin:" + c.conf.Plug)
	}
	c.js = js
	c.ctx, c.cncl = context.WithCancel(mgrCtx)
	go func() {
		defer ruisUtil.Recovers("gorun", nil)
		c.vm = goja.New()
		c.initVm()
		c.tg.Errs = ""
		err := c.run()
		if err != nil {
			c.tg.Errs = err.Error()
		}
		comm.Db.Cols("errs").Where("id=?", c.tg.Id).Update(c.tg)
		c.stop()
	}()
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				if c.vm != nil {
					c.vm.ClearInterrupt()
				}
				goto end
			default:
				time.Sleep(time.Millisecond)
			}
		}
	end:
		println("ctx end")
	}()
	return nil
}
func (c *trigHookTask) run() error {
	defer ruisUtil.Recovers("RunTask start", func(errs string) {
		println("trigHookTask run err:" + errs)
	})
	rslt, err := c.vm.RunString(c.js)
	if err != nil {
		println("vm.RunString err:" + err.Error())
		return err
	}
	println(fmt.Sprintf("js result:%+v", rslt.Export()))
	mainFun, ok := goja.AssertFunction(c.vm.Get("main"))
	if !ok {
		println("not found main err")
		return errors.New("not found main err")
	}

	ret, err := mainFun(goja.Null())
	if err != nil {
		println("vm mainFun err:" + err.Error())
		return err
	}
	rets := ruisUtil.NewMapo(ret.Export())
	fmt.Printf("rets:%+v\n", rets)

	if rets.GetBool("check") {
		rn := &model.TModelRun{}
		rn.Tid = c.md.Id
		rn.Uid = c.tg.Uid
		rn.Times = time.Now()
		rn.Tgid = c.tg.Id
		rn.Tgtyps = c.conf.Plug + "触发"
		comm.Db.Insert(rn)
		ExecMgr.Refresh()
		println("trigHookTask model run id:", rn.Id)
	} else {
		errs := rets.GetString("errs")
		if errs == "" {
			errs = "插件逻辑返回错误!"
		}
		return errors.New(errs)
	}
	return nil
}
func (c *trigHookTask) initVm() {
	csl := c.vm.NewObject()
	c.vm.Set("console", csl)
	csl.Set("log", func(args ...interface{}) {
		fmt.Println(args)
	})

	c.vm.Set("getHeader", func(key string) string {
		return c.header.Get(key)
	})
	c.vm.Set("getBodys", func(key string) string {
		return string(c.bodys)
	})
	c.vm.Set("getBody", func() interface{} {
		if c.body == nil {
			c.body = ruisUtil.NewMapo(c.bodys)
		}
		return c.body
	})
	c.vm.Set("getConf", func() interface{} {
		return c.confs
	})
}
