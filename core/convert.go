package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	ruisUtil "github.com/mgr9525/go-ruisutil"

	"git.code.oa.com/cloud_energy_studio/util/logger"
	"google.golang.org/protobuf/proto"
)

// Maps2Struct 转换
func Maps2Struct(mp *ruisUtil.Map, dist interface{}) error {
	return Map2Struct(mp.Map(), dist)
}

// Map2Struct 转换
func Map2Struct(mp map[string]interface{}, dist interface{}) (rterr error) {
	defer func() {
		if errs := recover(); errs != nil {
			rterr = errors.New(fmt.Sprintf("%v", errs))
		}
	}()
	if mp == nil || dist == nil {
		return errors.New("map/dist is nil")
	}

	tf := reflect.TypeOf(dist)
	vf := reflect.ValueOf(dist)
	if vf.IsNil() {
		return errors.New("dist is nil")
	}
	if tf.Kind() != reflect.Ptr {
		return errors.New("dist is not pointer")
	}
	tf = tf.Elem()
	vf = vf.Elem()
	if tf.Kind() != reflect.Struct {
		return errors.New("dist is not struct")
	}

	for i := 0; i < tf.NumField(); i++ {
		tfid := tf.Field(i)
		v, ok := mp[tfid.Name]
		if !ok {
			name := strings.Split(tfid.Tag.Get("json"), ",")[0]
			if name == "" {
				continue
			}
			v, ok = mp[name]
			if !ok {
				continue
			}
		}
		tfd := tfid.Type
		vfd := vf.FieldByIndex(tfid.Index)
		if !vfd.CanSet() {
			continue
		}
		if reflect.TypeOf(v).AssignableTo(tfd) {
			vfd.Set(reflect.ValueOf(v))
			continue
		}
		tfdc := tfd
		vfdc := vfd
		if tfd.Kind() == reflect.Ptr {
			tfdc = tfd.Elem()
			if vfd.Elem().Kind() == reflect.Invalid {
				nv := reflect.New(tfdc)
				vfd.Set(nv)
			}
			vfdc = vfd.Elem()
		}
		switch tfdc.Kind() {
		case reflect.String:
			vfdc.SetString(getString(v))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vfdc.SetInt(getInt(v))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vfdc.SetUint(getUint(v))
		case reflect.Float32, reflect.Float64:
			vfdc.SetFloat(getFloat(v))
		case reflect.Bool:
			vfdc.SetBool(getBool(v))
		case reflect.Struct:
			mpd, ok := v.(map[string]interface{})
			if ok {
				if err := Map2Struct(mpd, vfdc.Addr().Interface()); err != nil {
					println(err.Error())
				}
			}
		case reflect.Slice:
			ls, err := Obj2Slice(v, tfdc)
			if err != nil {
				println(err.Error())
			}
			vfd.Set(ls)
		default:
			setValue(vfdc, v)
		}
	}
	return nil
}

// pbStruct2Map p为protobuf生成的结构体指针
func pbStruct2Map(p interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if p == nil {
		return res
	}

	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)
	if t.Kind() == reflect.Ptr {
		if v.Elem().Kind() == reflect.Invalid {
			return res
		}
		v = reflect.Indirect(v)
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		tags := strings.Split(tag, ",")
		if len(tags) == 0 {
			continue
		}

		fValue := v.Field(i)
		fType := t.Field(i).Type.Kind()
		if fType == reflect.Ptr {
			if fValue.Kind() != reflect.Invalid {
				res[tags[0]] = pbStruct2Map(fValue.Interface())
				continue
			}
		} else if fType == reflect.Slice {
			lenth := fValue.Len()
			if lenth > 0 {
				out := make([]interface{}, lenth)
				for j := 0; j < lenth; j++ {
					tmpv := fValue.Index(j).Interface()
					if fValue.Index(j).Kind() == reflect.Ptr {
						out[j] = pbStruct2Map(tmpv)
					} else {
						out[j] = tmpv
					}
				}
				res[tags[0]] = out
			} else {
				res[tags[0]] = []interface{}{}
			}
			continue
		}

		res[tags[0]] = fValue.Interface()
	}
	return res
}

// Struct2Struct 转换
func Struct2Struct(src interface{}, dist interface{}) (rterr error) {
	dugName := ""
	defer func() {
		if errs := recover(); errs != nil {
			logger.Errorf("Struct2Struct dugName=%s, panic:%+v", dugName, errs)
			rterr = errors.New(fmt.Sprintf("dugName=%s, errs=%+v", dugName, errs))
		}
	}()
	if src == nil {
		return nil
	}
	if dist == nil {
		return errors.New("dist is nil")
	}

	tf := reflect.TypeOf(dist)
	vf := reflect.ValueOf(dist)
	if vf.IsNil() {
		return errors.New("dist is nil")
	}
	if tf.Kind() != reflect.Ptr {
		return errors.New("dist is not pointer")
	}
	tf = tf.Elem()
	vf = vf.Elem()
	if tf.Kind() != reflect.Struct {
		return errors.New("dist is not struct")
	}

	tf1 := reflect.TypeOf(src)
	vf1 := reflect.ValueOf(src)
	if vf1.IsNil() {
		return nil
	}
	if tf1.Kind() == reflect.Ptr {
		tf1 = tf1.Elem()
		vf1 = vf1.Elem()
	}
	if tf1.Kind() != reflect.Struct {
		return errors.New("src is not struct")
	}

	for i := 0; i < tf.NumField(); i++ {
		tfid := tf.Field(i)
		tf1fd, ok := tf1.FieldByName(tfid.Name)
		if !ok {
			lws := strings.ToLower(tfid.Name)
			tags := strings.Split(tfid.Tag.Get("conv"), ",")
			if lws == "id" {
				tags = []string{"ID", "Id", "id"}
			}
			if strings.Contains(lws, "id") {
				tags = append(tags, strings.ReplaceAll(tfid.Name, "Id", "ID"))
				tags = append(tags, strings.ReplaceAll(tfid.Name, "ID", "Id"))
			}
			for _, v := range tags {
				if v == "" {
					continue
				}
				tf1fd, ok = tf1.FieldByName(v)
				if ok {
					break
				}
			}
			if !ok {
				continue
			}
		}
		dugName = tf1fd.Name
		vf1fd := vf1.FieldByIndex(tf1fd.Index)
		if !vf1fd.CanInterface() {
			continue
		}
		v := vf1fd.Interface()
		tfd := tfid.Type
		vfd := vf.FieldByIndex(tfid.Index)
		if !vfd.CanSet() {
			continue
		}
		if tf1fd.Type.AssignableTo(tfd) {
			vfd.Set(vf1fd)
			continue
		}
		tfdc := tfd
		vfdc := vfd
		if tfd.Kind() == reflect.Ptr {
			tfdc = tfd.Elem()
			if vf1fd.Elem().Kind() == reflect.Invalid {
				continue
			}
			if vfd.Elem().Kind() == reflect.Invalid {
				nv := reflect.New(tfdc)
				vfd.Set(nv)
			}
			vfdc = vfd.Elem()
		}
		switch vfdc.Kind() {
		case reflect.String:
			vfdc.SetString(getString(v))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vfdc.SetInt(getInt(v))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vfdc.SetUint(getUint(v))
		case reflect.Float32, reflect.Float64:
			vfdc.SetFloat(getFloat(v))
		case reflect.Bool:
			vfdc.SetBool(getBool(v))
		case reflect.Struct:
			if err := Struct2Struct(v, vfdc.Addr().Interface()); err != nil {
				println(err.Error())
			}
		case reflect.Slice:
			ls, err := Obj2Slice(v, tfdc)
			if err != nil {
				println(err.Error())
			}
			vfdc.Set(ls)
		default:
			setValue(vfd, v)
		}
	}
	return nil
}

func setValue(vf reflect.Value, v interface{}) {
	defer func() {
		if errs := recover(); errs != nil {
			logger.Errorf("setValue name:%s,err=%s", vf.String(), errs)
		}
	}()
	vf.Set(reflect.ValueOf(v))
}
func getBool(v interface{}) bool {
	vf := reflect.ValueOf(v)
	switch vf.Kind() {
	case reflect.Bool:
		return v.(bool)
	case reflect.String:
		return v.(string) == "true"
	}
	return false
}
func getString(v interface{}) string {
	vf := reflect.ValueOf(v)
	switch vf.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", v)
	}
	return fmt.Sprintf("%v", v)
}
func getInt(v interface{}) int64 {
	vf := reflect.ValueOf(v)
	switch vf.Kind() {
	case reflect.Int:
		return int64(v.(int))
	case reflect.Int8:
		return int64(v.(int8))
	case reflect.Int16:
		return int64(v.(int16))
	case reflect.Int32:
		return int64(v.(int32))
	case reflect.Int64:
		return v.(int64)
	case reflect.Uint:
		return int64(v.(uint))
	case reflect.Uint8:
		return int64(v.(uint8))
	case reflect.Uint16:
		return int64(v.(uint16))
	case reflect.Uint32:
		return int64(v.(uint32))
	case reflect.Uint64:
		return int64(v.(uint64))
	case reflect.Float32:
		return int64(v.(float32))
	case reflect.Float64:
		return int64(v.(float64))
	case reflect.String:
		vc, _ := strconv.ParseInt(v.(string), 10, 64)
		return vc
	}
	return 0
}
func getUint(v interface{}) uint64 {
	vf := reflect.ValueOf(v)
	switch vf.Kind() {
	case reflect.Int:
		return uint64(v.(int))
	case reflect.Int8:
		return uint64(v.(int8))
	case reflect.Int16:
		return uint64(v.(int16))
	case reflect.Int32:
		return uint64(v.(int32))
	case reflect.Int64:
		return uint64(v.(int64))
	case reflect.Uint:
		return uint64(v.(uint))
	case reflect.Uint8:
		return uint64(v.(uint8))
	case reflect.Uint16:
		return uint64(v.(uint16))
	case reflect.Uint32:
		return uint64(v.(uint32))
	case reflect.Uint64:
		return v.(uint64)
	case reflect.Float32:
		return uint64(v.(float32))
	case reflect.Float64:
		return uint64(v.(float64))
	case reflect.String:
		vc, _ := strconv.ParseUint(v.(string), 10, 64)
		return vc
	}
	return 0
}
func getFloat(v interface{}) float64 {
	vf := reflect.ValueOf(v)
	switch vf.Kind() {
	case reflect.Int:
		return float64(v.(int))
	case reflect.Int8:
		return float64(v.(int8))
	case reflect.Int16:
		return float64(v.(int16))
	case reflect.Int32:
		return float64(v.(int32))
	case reflect.Int64:
		return float64(v.(int64))
	case reflect.Uint:
		return float64(v.(uint))
	case reflect.Uint8:
		return float64(v.(uint8))
	case reflect.Uint16:
		return float64(v.(uint16))
	case reflect.Uint32:
		return float64(v.(uint32))
	case reflect.Uint64:
		return float64(v.(uint64))
	case reflect.Float32:
		return float64(v.(float32))
	case reflect.Float64:
		return v.(float64)
	case reflect.String:
		vc, _ := strconv.ParseFloat(v.(string), 64)
		return vc
	}
	return 0
}

// Bytes2Struct 转换
func Bytes2Struct(bts []byte, dist interface{}) error {
	mp := make(map[string]interface{})
	err := json.Unmarshal(bts, &mp)
	if err != nil {
		return err
	}
	return Map2Struct(mp, dist)
}

// Obj2Slice 转换
func Obj2Slice(obj interface{}, dtf reflect.Type) (ret reflect.Value, rterr error) {
	defer func() {
		if errs := recover(); errs != nil {
			logger.Errorf("Obj2Slice dtfName=%s, panic:%+v", dtf.Name(), errs)
			rterr = errors.New(fmt.Sprintf("dtfName=%s, errs=%+v", dtf.Name(), errs))
		}
	}()
	objs := reflect.ValueOf(obj)
	if objs.Kind() != reflect.Slice {
		return reflect.Value{}, errors.New("objs is not slice")
	}
	dist := reflect.MakeSlice(dtf, 0, 0)

	dtf1 := dtf.Elem()
	dtf2 := dtf1
	if dtf1.Kind() == reflect.Ptr {
		dtf2 = dtf1.Elem()
	}

	for i := 0; i < objs.Len(); i++ {
		valf := objs.Index(i)
		val := valf.Interface()
		nv := reflect.New(dtf2)
		tf := reflect.TypeOf(val)
		vf := reflect.ValueOf(val)
		if tf.Kind() == reflect.Ptr {
			tf = tf.Elem()
			if vf.Elem().Kind() == reflect.Invalid {
				continue
			}
			vf = vf.Elem()
		}
		switch dtf2.Kind() {
		case reflect.String:
			nv = nv.Elem()
			nv.SetString(getString(val))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			nv = nv.Elem()
			nv.SetInt(getInt(val))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			nv = nv.Elem()
			nv.SetUint(getUint(val))
		case reflect.Float32, reflect.Float64:
			nv = nv.Elem()
			nv.SetFloat(getFloat(val))
		case reflect.Bool:
			nv = nv.Elem()
			nv.SetBool(getBool(val))
		case reflect.Struct:
			mpd, ok := val.(map[string]interface{})
			if ok {
				if err := Map2Struct(mpd, nv.Interface()); err != nil {
					println(err.Error())
				}
			} else {
				if err := Struct2Struct(val, nv.Interface()); err != nil {
					println(err.Error())
				}
			}
		case reflect.Slice:
			ls, err := Obj2Slice(val, dtf2)
			if err != nil {
				println(err.Error())
			}
			nv.Set(ls)
		default:
			nv = valf
		}
		dist = reflect.Append(dist, nv)
	}
	return dist, nil
}

/*// Obj2Slice 转换
func Obj2Slice(obj interface{}, dist interface{}) error {
	bts, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bts, dist)
	if err != nil {
		return err
	}
	return nil
}*/

// PutProtoMsgData 填充rsp.Data字段
func PutProtoMsgData(rsp proto.Message) {
	defer func() {
		if errs := recover(); errs != nil {
			logger.Errorf("ProtoMsgData ,err=%s", errs)
		}
	}()
	tf := reflect.TypeOf(rsp)
	vf := reflect.ValueOf(rsp)
	if tf.Kind() != reflect.Ptr {
		return
	}

	tf = tf.Elem()
	vf = vf.Elem()
	if vf.Kind() == reflect.Invalid {
		return
	}

	tfd, ok := tf.FieldByName("Data")
	if !ok || tfd.Type.Kind() != reflect.Ptr {
		return
	}
	vfd := vf.FieldByIndex(tfd.Index)
	if vfd.Elem().Kind() != reflect.Invalid {
		return
	}
	vfd.Set(reflect.New(tfd.Type.Elem()))
}
