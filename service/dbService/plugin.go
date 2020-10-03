package dbService

import (
	"gokins/comm"
	"gokins/model"
)

func GetPlugin(id int) *model.TPlugin {
	if id <= 0 {
		return nil
	}
	e := new(model.TPlugin)
	ok, err := comm.Db.Where("id=?", id).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}
func GetPluginRun(id int) *model.TPluginRun {
	if id <= 0 {
		return nil
	}
	e := new(model.TPluginRun)
	ok, err := comm.Db.Where("id=?", id).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}
func FindPluginRun(mid, tid, pid int) *model.TPluginRun {
	if mid <= 0 || tid <= 0 || pid <= 0 {
		return nil
	}
	e := new(model.TPluginRun)
	ok, err := comm.Db.Where("mid=? and tid=? and pid=?", mid, tid, pid).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}
