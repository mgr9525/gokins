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
	"regexp"
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
	if c.Mr == nil || c.cncl != nil {
		return
	}
	c.Md = dbService.GetModel(c.Mr.Tid)
	if c.Md == nil {
		return
	}
	c.ctx, c.cncl = context.WithCancel(mgrCtx)
	go func() {
		c.run()
		c.cncl = nil
	}()
}

func (c *RunTask) run() {
	defer ruisUtil.Recovers("RunTask start", nil)
	c.plugs = nil
	if c.Md.Wrkdir != "" {
		if ruisIo.PathExists(c.Md.Wrkdir) {
			if c.Md.Clrdir == 1 {
				err := rmDirFiles(c.Md.Wrkdir)
				if err != nil {
					c.end(2, "工作目录创建失败:"+err.Error())
					return
				}
			}
		} else {
			/*if c.Md.Clrdir != 1 {
				c.end(2, "工作目录不存在")
				return
			}*/
			err := os.MkdirAll(c.Md.Wrkdir, 0755)
			if err != nil {
				c.end(2, "工作目录创建失败:"+err.Error())
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

var regPATH = regexp.MustCompile(`^PATH=`)

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
		noPath := true
		str := strings.ReplaceAll(c.Md.Envs, "\t", "")
		envs := strings.Split(str, "\n")
		for i, s := range envs {
			if regPATH.MatchString(s) {
				noPath = false
				envs[i] = strings.ReplaceAll(s, "$PATH", os.Getenv("PATH"))
				envs[i] = strings.ReplaceAll(s, "${PATH}", os.Getenv("PATH"))
			}
		}
		if noPath {
			envs = append(envs, "PATH="+os.Getenv("PATH"))
		}
		envs = append(envs, "WORKDIR="+c.Md.Wrkdir)
		cmd.Env = envs
	}
	if c.Md.Wrkdir != "" {
		cmd.Dir = c.Md.Wrkdir
	} else if comm.Dir != "" {
		cmd.Dir = comm.Dir
	}
	err = cmd.Run()
	rn.State = 4
	if err != nil {
		println("cmd.run err:" + err.Error())
		// rn.State = 2
		// return rn, err
	}
	fmt.Println(fmt.Sprintf("cmdRun(%s)dir:%s", pgn.Title, cmd.Dir))
	if cmd.ProcessState != nil {
		rn.Excode = cmd.ProcessState.ExitCode()
	}
	if pgn.Exend == 1 && (err != nil || rn.Excode != 0) {
		rn.State = 2
		return rn, fmt.Errorf("程序执行错误(exit:%d)：%+v", rn.Excode, err)
	}
	return rn, nil
}
func (c *RunTask) end(stat int, errs string) {
	defer c.stop()
	c.Mr.State = stat
	c.Mr.Errs = errs
	c.Mr.Timesd = time.Now()
	_, err := comm.Db.Cols("state", "errs", "timesd").Where("id=?", c.Mr.Id).Update(c.Mr)
	if err != nil {
		println("db err:", err.Error())
		return
	}

	var ls []*model.TTrigger
	err = comm.Db.Where("del!=1 and enable=1 and meid=?", c.Md.Id).Find(&ls)
	if err != nil {
		println("db err:", err.Error())
		return
	}
	for _, v := range ls {
		TriggerMgr.StartOne(v, c.Md, c.Mr)
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
