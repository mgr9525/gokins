package service

import (
	"fmt"
	"gokins/comm"
	"gokins/model"
	"gokins/service/dbService"
	"time"

	"github.com/go-xorm/xorm"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

var dbold *xorm.Engine

func MoveModels() {
	var olds []*ruisUtil.Map
	err := dbold.SQL("select * from t_model").Find(&olds)
	if err != nil {
		fmt.Println("find model err:" + err.Error())
		return
	}
	for _, v := range olds {
		id, err := v.GetInt("id")
		if err != nil {
			continue
		}
		del, err := v.GetInt("del")
		if err != nil {
			continue
		}
		clrdir, err := v.GetInt("clrdir")
		if err != nil {
			continue
		}
		ne := &model.TModel{}
		ne.Uid = v.GetString("uid")
		ne.Title = v.GetString("title")
		ne.Desc = v.GetString("desc")
		ne.Del = int(del)
		ne.Clrdir = int(clrdir)
		ne.Envs = v.GetString("envs")
		ne.Wrkdir = v.GetString("wrkdir")
		if tm, ok := v.Get("times").(time.Time); ok {
			ne.Times = tm
		}
		_, err = comm.Db.Insert(ne)
		if err == nil {
			mvPlugin(id, ne)
		}
	}
}

func mvPlugin(tid int64, md *model.TModel) {
	var olds []*ruisUtil.Map
	err := dbold.SQL("select * from t_plugin where tid=?", tid).Find(&olds)
	if err != nil {
		println("find model err:" + err.Error())
		return
	}
	for _, v := range olds {
		del, err := v.GetInt("del")
		if err != nil {
			continue
		}
		typ, err := v.GetInt("type")
		if err != nil {
			continue
		}
		sort, err := v.GetInt("sort")
		if err != nil {
			break
		}
		exend, err := v.GetInt("exend")
		if err != nil {
			continue
		}
		ne := &model.TPlugin{}
		ne.Tid = md.Id
		ne.Title = v.GetString("title")
		ne.Type = int(typ)
		ne.Del = int(del)
		ne.Sort = int(sort)
		ne.Exend = int(exend)
		ne.Para = v.GetString("para")
		ne.Cont = v.GetString("cont")
		if tm, ok := v.Get("times").(time.Time); ok {
			ne.Times = tm
		}
		_, err = comm.Db.Insert(ne)
		if err != nil {
			println("insert plug err:" + err.Error())
			break
		}
	}
}

func MoveTrigger() {
	var olds []*ruisUtil.Map
	err := dbold.SQL("select * from t_trigger").Find(&olds)
	if err != nil {
		//fmt.Println("find trigger err:" + err.Error())
		return
	}
	for _, v := range olds {
		/*id, err := v.GetInt("id")
		if err != nil {
			continue
		}*/
		del, err := v.GetInt("del")
		if err != nil {
			continue
		}
		enable, err := v.GetInt("enable")
		if err != nil {
			continue
		}
		ne := &model.TTrigger{}
		ne.Uid = v.GetString("uid")
		ne.Types = v.GetString("types")
		ne.Title = v.GetString("title")
		ne.Desc = v.GetString("desc")
		ne.Config = v.GetString("config")
		ne.Del = int(del)
		ne.Enable = int(enable)
		ne.Errs = v.GetString("errs")
		if tm, ok := v.Get("times").(time.Time); ok {
			ne.Times = tm
		}
		_, err = comm.Db.Insert(ne)
		if err != nil {
			println("MoveTrigger err:" + err.Error())
			return
		}
	}
}

func MoveUser() {
	var olds []*ruisUtil.Map
	err := dbold.SQL("select * from sys_user").Find(&olds)
	if err != nil {
		fmt.Println("find trigger err:" + err.Error())
		return
	}
	for _, v := range olds {
		isup := true
		ne := dbService.FindUserName(v.GetString("name"))
		if ne == nil {
			isup = false
			ne = &model.SysUser{}
			ne.Xid = v.GetString("xid")
			ne.Name = v.GetString("name")
		}
		ne.Pass = v.GetString("pass")
		ne.Phone = v.GetString("phone")
		ne.Avat = v.GetString("avat")
		if tm, ok := v.Get("times").(time.Time); ok {
			ne.Times = tm
		}
		if isup {
			_, err = comm.Db.Cols("pass", "phone", "avat", "times").Where("id=?", ne.Id).Update(ne)
		} else {
			_, err = comm.Db.Insert(ne)
		}
		if err != nil {
			println("MoveTrigger err:" + err.Error())
			return
		}
	}
}
func MoveParam() {
	var olds []*ruisUtil.Map
	err := dbold.SQL("select * from sys_param").Find(&olds)
	if err != nil {
		fmt.Println("find trigger err:" + err.Error())
		return
	}
	for _, v := range olds {
		cont, ok := v.Get("cont").([]byte)
		if !ok {
			continue
		}
		ne := &model.SysParam{}
		ne.Key = v.GetString("key")
		ne.Cont = cont
		if tm, ok := v.Get("times").(time.Time); ok {
			ne.Times = tm
		}
		_, err = comm.Db.Insert(ne)
		if err != nil {
			println("MoveTrigger err:" + err.Error())
			return
		}
	}
}
