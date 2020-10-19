package mgr

import (
	"context"
	"errors"
	"fmt"
	"gokins/comm"
	"gokins/model"
	"gokins/service/dbService"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ruisIo "github.com/mgr9525/go-ruisutil/ruisio"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

type RunTask struct {
	Md    *model.TModel
	Mr    *model.TModelRun
	plugs []*model.TPlugin
	ctx   context.Context
	cncl  context.CancelFunc
}

func (c *RunTask) start() {
	if c.cncl != nil {
		return
	}
	c.Md = dbService.GetModel(c.Mr.Tid)
	if c.Md == nil {
		return
	}
	c.ctx, c.cncl = context.WithCancel(context.Background())
	go c.run()
}

func (c *RunTask) run() {
	defer ruisUtil.Recovers("RunTask start", nil)
	c.plugs = nil
	if c.Md.Wrkdir != "" {
		if ruisIo.PathExists(c.Md.Wrkdir) {
			if c.Md.Clrdir == 1 {
				err := rmDirFiles(c.Md.Wrkdir)
				if err != nil {
					c.end(2, "运行目录创建失败:"+err.Error())
					return
				}
			}
		} else {
			if c.Md.Clrdir != 1 {
				c.end(2, "运行目录不存在")
				return
			}
			err := os.MkdirAll(c.Md.Wrkdir, 0755)
			if err != nil {
				c.end(2, "运行目录创建失败:"+err.Error())
				return
			}
		}
	}
	err := comm.Db.Where("del!='1' and tid=?", c.Mr.Tid).OrderBy("sort ASC,id ASC").Find(&c.plugs)
	if err != nil {
		c.end(2, "db err:"+err.Error())
		return
	}
	if len(c.plugs) <= 0 {
		c.end(2, "无插件")
		return
	}
	for _, v := range c.plugs {
		select {
		case <-c.ctx.Done():
			c.end(-1, "手动停止")
			return
		default:
			rn, err := c.runs(v)
			if rn != nil {
				rn.Timesd = time.Now()
				_, errs := comm.Db.Cols("state", "excode", "timesd").Where("id=?", rn.Id).Update(rn)
				if err == nil && errs != nil {
					err = errs
				}
			}
			if err != nil {
				println("cmd run err:", err.Error())
				c.end(2, err.Error())
				return
			}
			time.Sleep(time.Second)
		}
	}
	c.end(4, "")
}

func (c *RunTask) runs(pgn *model.TPlugin) (rns *model.TPluginRun, rterr error) {
	defer ruisUtil.Recovers("RunTask.run", func(errs string) {
		rterr = errors.New(errs)
	})
	rn := dbService.FindPluginRun(c.Mr.Tid, c.Mr.Id, pgn.Id)
	if rn == nil {
		rn = &model.TPluginRun{Mid: c.Mr.Tid, Tid: c.Mr.Id, Pid: pgn.Id}
		rn.Times = time.Now()
		rn.State = 1
		_, err := comm.Db.Insert(rn)
		if err != nil {
			return nil, err
		}
	} else if rn.State != 0 {
		rn.State = 2
		return rn, errors.New("already run")
	}
	select {
	case <-c.ctx.Done():
		rn.State = -1
		return rn, nil
	default:
	}
	logpth := fmt.Sprintf("%s/data/logs/%d/%d.log", comm.Dir, rn.Tid, rn.Id)
	if !ruisIo.PathExists(filepath.Dir(logpth)) {
		err := os.MkdirAll(filepath.Dir(logpth), 0755)
		if err != nil {
			println("MkdirAll err:" + err.Error())
			rn.State = 2
			return rn, err
		}
	}
	logfl, err := os.OpenFile(logpth, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		println("MkdirAll err:" + err.Error())
		rn.State = 2
		return rn, err
	}
	defer logfl.Close()
	name := "sh"
	par0 := "-c"
	if runtime.GOOS == "windows" {
		name = "cmd"
		par0 = "/c"
	}
	cmd := exec.CommandContext(c.ctx, name, par0, pgn.Cont)
	cmd.Stdout = logfl
	cmd.Stderr = logfl
	if c.Md.Envs != "" {
		str := strings.ReplaceAll(c.Md.Envs, "\t", "")
		envs := strings.Split(str, "\n")
		cmd.Env = envs
	}
	if c.Md.Wrkdir != "" {
		cmd.Dir = c.Md.Wrkdir
	}
	err = cmd.Run()
	rn.State = 4
	if err != nil {
		println("cmd.run err:" + err.Error())
		rn.State = 2
		return rn, err
	}
	fmt.Println(fmt.Sprintf("cmdRun(%s)dir:%s", pgn.Title, cmd.Dir))
	if cmd.ProcessState != nil {
		rn.Excode = cmd.ProcessState.ExitCode()
	}
	if rn.Excode != 0 {
		rn.State = 2
		if pgn.Exend == 1 {
			return rn, fmt.Errorf("程序执行错误：%d", rn.Excode)
		}
	}
	return rn, nil
}
func (c *RunTask) end(stat int, errs string) {
	c.stop()
	v := &model.TModelRun{}
	v.State = stat
	v.Errs = errs
	v.Timesd = time.Now()
	_, err := comm.Db.Cols("state", "errs", "timesd").Where("id=?", c.Mr.Id).Update(v)
	if err != nil {
		println("db err:", err.Error())
	}
}

func (c *RunTask) stop() {
	if c.cncl != nil {
		c.cncl()
		c.cncl = nil
	}
}

func rmDirFiles(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
