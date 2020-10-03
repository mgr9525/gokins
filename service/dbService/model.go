package dbService

import (
	"gokins/comm"
	"gokins/model"
)

func GetModel(id int) *model.TModel {
	if id <= 0 {
		return nil
	}
	e := new(model.TModel)
	ok, err := comm.Db.Where("id=?", id).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}
func GetModelRun(id int) *model.TModelRun {
	if id <= 0 {
		return nil
	}
	e := new(model.TModelRun)
	ok, err := comm.Db.Where("id=?", id).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}
