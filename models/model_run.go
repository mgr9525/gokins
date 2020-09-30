package models

import (
	"gokins/comm"
	"gokins/service/dbService"
	"time"
)

type ModelRun struct {
	Id     int `xorm:"pk autoincr"`
	Tid    int
	Uid    string
	Times  time.Time
	Timesd time.Time
	State  int

	Nick   string `xorm:"-"`
	Times1 string `xorm:"-"`
	Times2 string `xorm:"-"`
}

func (ModelRun) TableName() string {
	return "t_model_run"
}

func (c *ModelRun) Add() error {
	c.State = 0
	c.Times = time.Now()
	_, err := comm.Db.Insert(c)
	return err
}
func (c *ModelRun) ToUI() {
	c.Times1 = c.Times.Format(comm.TimeFmt)
	if !c.Timesd.IsZero() {
		c.Times2 = c.Timesd.Format(comm.TimeFmt)
	}
	usr := dbService.FindUser(c.Uid)
	if usr != nil {
		c.Nick = usr.Nick
	}
}
