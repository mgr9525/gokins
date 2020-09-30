package core

import (
	"errors"
	"reflect"

	"github.com/go-xorm/xorm"
)

type Page struct {
	Page  int64
	Size  int64
	Pages int64
	Total int64
	Data  interface{}
}

func XormFindCount(ses *xorm.Session, rowsSlicePtr interface{}) (int64, error) {
	sliceValue := reflect.Indirect(reflect.ValueOf(rowsSlicePtr))
	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Map {
		return 0, errors.New("needs a pointer to a slice or a map")
	}

	sliceElementType := sliceValue.Type().Elem()

	if sliceElementType.Kind() == reflect.Ptr {
		if sliceElementType.Elem().Kind() == reflect.Struct {
			pv := reflect.New(sliceElementType.Elem())
			return ses.Clone().Count(pv.Interface())
		}
	} else if sliceElementType.Kind() == reflect.Struct {
		pv := reflect.New(sliceElementType)
		return ses.Clone().Count(pv.Interface())
	}
	return 0, errors.New("not found table")
}

func XormFindPage(ses *xorm.Session, ls interface{}, page int64, size ...int64) (*Page, error) {
	count, err := XormFindCount(ses, ls)
	if err != nil {
		return nil, err
	}
	return XormFindPages(ses, ls, count, page, size...)
}
func XormFindPages(ses *xorm.Session, ls interface{}, count, page int64, size ...int64) (*Page, error) {
	var pageno int64 = 1
	var sizeno int64 = 10
	var pagesno int64 = 0
	//var count=c.FindCount(pars)
	if page > 0 {
		pageno = page
	}
	if len(size) > 0 && size[0] > 0 {
		sizeno = size[0]
	}
	start := (pageno - 1) * sizeno
	err := ses.Limit(int(sizeno), int(start)).Find(ls)
	if err != nil {
		return nil, err
	}
	pagest := count / sizeno
	if count%sizeno > 0 {
		pagesno = pagest + 1
	} else {
		pagesno = pagest
	}
	return &Page{
		Page:  pageno,
		Pages: pagesno,
		Size:  sizeno,
		Total: count,
		Data:  ls,
	}, nil
}
