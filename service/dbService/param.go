package dbService

import (
	"gokins/comm"
	"gokins/model"
	"time"

	ruisUtil "github.com/mgr9525/go-ruisutil"
)

func FindParam(key string) *model.SysParam {
	if key == "" {
		return nil
	}
	e := new(model.SysParam)
	ok, err := comm.Db.Where("`key`=?", key).Get(e)
	if err != nil {
		return nil
	}
	if ok {
		return e
	}
	return nil
}

var mkey = []byte("QXQBPH1X6RRUWNRQ")
var miv = []byte("832S6MU5LTG0A20K")

func GetParam(key string) *ruisUtil.Map {
	ret := ruisUtil.NewMap()
	v := FindParam(key)
	if v != nil {
		bts, err := ruisUtil.AESDecrypt(v.Cont, mkey, miv)
		if err == nil {
			ret = ruisUtil.NewMapo(bts)
		}
	}
	return ret
}
func SetParam(key string, v *ruisUtil.Map) error {
	bts, err := ruisUtil.AESEncrypt(v.ToBytes(), mkey, miv)
	if err != nil {
		return err
	}
	para := FindParam(key)
	if para == nil {
		para = &model.SysParam{}
		para.Key = key
		para.Cont = bts
		para.Times = time.Now()
		_, err := comm.Db.Insert(para)
		return err
	} else {
		para.Cont = bts
		_, err := comm.Db.Where("id=?", para.Id).Update(para)
		return err
	}
}
