package comm

import (
	"bytes"
	"gokins/model"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func InitDb() error {
	db, err := xorm.NewEngine("sqlite3", Dir+"/db.dat")
	if err != nil {
		return err
	}
	Db = db
	isext, err := Db.IsTableExist(model.SysUser{})
	if err == nil && !isext {
		_, err := Db.Import(bytes.NewBufferString(sqls))
		if err != nil {
			println("Db.Import err:" + err.Error())
		}
		//e:=&models.SysUser{}
		//e.Times=time.Now()
		//db.Cols("times").Where("xid=?","admin").Update(e)
	}
	return nil
}
