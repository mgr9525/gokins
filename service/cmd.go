package service

import (
	"fmt"
	"gokins/comm"
	"gokins/service/dbService"
	"os"
	"path/filepath"

	"github.com/go-xorm/xorm"
)

func ClearUPass(nm string) {
	if nm == "" {
		return
	}
	usr := dbService.FindUserName(nm)
	if usr == nil {
		fmt.Printf("user(%s) not found\n", nm)
	} else {
		usr.Pass = ""
		_, err := comm.Db.Cols("pass").Where("id=?", usr.Id).Update(usr)
		if err != nil {
			fmt.Println("clear password err:" + err.Error())
		} else {
			fmt.Println("clear password ok")
		}
	}
}

func MoveData(pth string) {
	if pth == "" {
		return
	}
	pths, err := filepath.Abs(pth)
	if err != nil {
		fmt.Println("old db path err:" + err.Error())
		return
	}
	stat, err := os.Stat(pths)
	if err != nil {
		fmt.Println("old db path err:" + err.Error())
		return
	}
	if stat.IsDir() {
		pths += "/db.dat"
		stat, err = os.Stat(pths)
		if err != nil {
			fmt.Println("find old db path err:" + err.Error())
			return
		}
		if stat.IsDir() {
			fmt.Println("find old db path err:is dir")
			return
		}
	}

	db, err := xorm.NewEngine("sqlite3", pths)
	if err != nil {
		fmt.Println("old db path err:" + err.Error())
		return
	}
	dbold = db
	MoveModels()
}
