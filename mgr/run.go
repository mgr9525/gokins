package mgr

import (
	"bytes"
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
	"sync"
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
	lk    sync.Mutex
	stds  map[int]*bytes.Buffer
	//outs  map[int]*string
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
	c.stds = make(map[int]*bytes.Buffer)
	//c.outs = make(map[int]*string)
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
			err := c.runs(v)
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

func (c *RunTask) runs(pgn *model.TPlugin) (rterr error) {
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
			return err
		}
	}
	if rn.State != 1 {
		return nil
	}
	select {
	case <-c.ctx.Done():
		return nil
	default:
	}
	name := "sh"
	par0 := "-c"
	if runtime.GOOS == "windows" {
		name = "cmd"
		par0 = "/c"
	}
	stdout := &bytes.Buffer{}
	c.lk.Lock()
	c.stds[rn.Id] = stdout
	c.lk.Unlock()
	cmd := exec.CommandContext(c.ctx, name, par0, pgn.Cont)
	cmd.Stdout = stdout
	cmd.Stderr = stdout
	if c.Md.Envs != "" {
		str := strings.ReplaceAll(c.Md.Envs, "\t", "")
		envs := strings.Split(str, "\n")
		cmd.Env = envs
	}
	if c.Md.Wrkdir != "" {
		cmd.Dir = c.Md.Wrkdir
	}
	err := cmd.Run()
	rn.State = 4
	if err != nil {
		println("cmd.run err:" + err.Error())
		rn.State = 2
	}
	println(fmt.Sprintf("cmdRun(%s)dir:%s", pgn.Title, cmd.Dir))
	if cmd.ProcessState != nil {
		rn.Excode = cmd.ProcessState.ExitCode()
	}

	c.lk.Lock()
	defer c.lk.Unlock()
	rn.Output = stdout.String()
	/*if c.outs[rn.Id] == nil {
		rn.Output = stdout.String()
	} else {
		*c.outs[rn.Id] += stdout.String()
		rn.Output = *c.outs[rn.Id]
	}*/

	rn.Timesd = time.Now()
	_, err = comm.Db.Cols("state", "excode", "timesd", "output").Where("id=?", rn.Id).Update(rn)
	if err != nil {
		return err
	}
	delete(c.stds, rn.Id)
	if pgn.Exend == 1 && rn.Excode != 0 {
		return fmt.Errorf("程序执行错误：%d", rn.Excode)
	}
	return nil
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

func (c *RunTask) read(id int) string {
	c.lk.Lock()
	defer c.lk.Unlock()
	if c.stds[id] == nil {
		return ""
	}
	return c.stds[id].String()
	/*if c.outs[id] == nil {
		s := ""
		c.outs[id] = &s
	}
	out := c.stds[id].String()
	*c.outs[id] += out
	return *c.outs[id]*/
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
