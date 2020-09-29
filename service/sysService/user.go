package sysService

import (
	"gokins/comm"
	"gokins/model"
)

func FindUser(xid string) *model.SysUser {
	if xid == "" {
		return nil
	}
	e := new(model.SysUser)
	ok, err := comm.Db.Where("xid=?", xid).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}
func FindUserName(nm string) *model.SysUser {
	if nm == "" {
		return nil
	}
	e := new(model.SysUser)
	ok, err := comm.Db.Where("name=?", nm).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}
