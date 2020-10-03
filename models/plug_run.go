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
