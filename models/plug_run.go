package models

import (
	"fmt"
	"gokins/service/dbService"
	"time"
)

type PluginUI struct {
	Id    int `xorm:"pk autoincr"`
	Tid   int //model id
	Type  int
	Title string
	Times time.Time
	Sort  int
	Exend int

	RunStat int    `xorm:"-"`
	Hstm    string `xorm:"-"`
}

func (PluginUI) TableName() string {
	return "t_plugin"
}

func (c *PluginUI) ToUI(mrid int) {
	run := dbService.FindPluginRun(c.Tid, mrid, c.Id)
	if run != nil {
		c.RunStat = run.State
		c.Hstm = "0"
		if run.State >= 2 {
			c.Hstm = fmt.Sprintf("%.3f", run.Timesd.Sub(run.Times).Seconds())
		} else if run.State == 1 {
			c.Hstm = fmt.Sprintf("%.3f", time.Since(run.Times).Seconds())
		}
	}
}

/*type PluginRunUI struct {
	Id    int `xorm:"pk autoincr"`
	Pid    int //plugin id
	Mid    int //model id
	Tid    int //modelRun id
	Times  time.Time
	Timesd time.Time
	State  int //-1已停止，0等待，1运行，2运行失败，4运行成功
	Excode int

	Title int    `xorm:"-"`
	Hstm    string `xorm:"-"`
}

func (PluginRunUI) TableName() string {
	return "t_plugin_run"
}

func (c *PluginRunUI) ToUI() {
	run := dbService.GetPlugin(c.Pid)
	if run != nil {
		c.Hstm = "0"
		if c.State >= 2 {
			c.Hstm = fmt.Sprintf("%.3f", c.Timesd.Sub(run.Times).Seconds())
		} else if c.State == 1 {
			c.Hstm = fmt.Sprintf("%.3f", time.Since(run.Times).Seconds())
		}
	}
}
*/
