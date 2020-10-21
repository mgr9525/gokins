package dbService

import (
	"gokins/comm"
	"gokins/model"
)

func GetTrigger(id int) *model.TTrigger {
	if id <= 0 {
		return nil
	}
	e := new(model.TTrigger)
	ok, err := comm.Db.Where("id=?", id).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}
