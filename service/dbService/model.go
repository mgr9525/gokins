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
