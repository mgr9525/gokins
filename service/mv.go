package service

import (
	"fmt"
	"gokins/comm"
	"gokins/model"
	"time"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

func MoveModels() {
	var olds []*ruisUtil.Map
	err := comm.Dbold.SQL("select * from t_model").Find(&olds)
	if err != nil {
		fmt.Println("find model err:" + err.Error())
		return
	}
	for _, v := range olds {
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
			mvPlugin(ne)
		}
	}
}

func mvPlugin(md *model.TModel) {
	var olds []*ruisUtil.Map
	err := comm.Dbold.SQL("select * from t_plugin").Find(&olds)
	if err != nil {
		fmt.Println("find model err:" + err.Error())
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
			fmt.Sprintf("insert plug err:" + err.Error())
			break
		}
	}
}
