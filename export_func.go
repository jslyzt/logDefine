package logDefine

import (
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func allNumber(key string) bool {
	if len(key) <= 0 {
		return false
	}
	for _, tk := range []byte(key) {
		if tk < '0' || tk > '9' {
			return false
		}
	}
	return true
}

func reflect2string(any interface{}) string {
	anyt := reflect.TypeOf(any)
	anyk := anyt.Kind()
	anyv := reflect.ValueOf(any)
	switch anyk {
	case reflect.Bool:
		return any2string(anyv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return any2string(anyv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return any2string(anyv.Uint())
	case reflect.Float32, reflect.Float64:
		return any2string(anyv.Float())
	case reflect.Complex64, reflect.Complex128:
		return any2string(anyv.Complex())
	case reflect.String:
		return anyv.String()
	case reflect.Struct, reflect.Ptr:
		if anyv.IsValid() {
			tfunc := anyv.MethodByName("ToString")
			if tfunc.IsValid() {
				rel := tfunc.Call([]reflect.Value{})
				if len(rel) > 0 {
					return fmt.Sprintf("(%s)", rel[0].String())
				}
			}
		}
	case reflect.Map:
		ostr := "{"
		for _, v := range anyv.MapKeys() {
			ostr = ostr + fmt.Sprintf("%s:%s;", any2string(v), any2string(anyv.MapIndex(v)))
		}
		return ostr + "}"
	case reflect.Slice, reflect.Array:
		ostr := "["
		for i := 0; i < anyv.Len(); i++ {
			ostr = ostr + fmt.Sprintf("%s,", any2string(anyv.Index(i)))
		}
		return ostr + "]"
	}
	return ""
}

func any2string(any interface{}) string {
	value := ""
	switch anyi := any.(type) {
	case string:
		value = any.(string)
	case bool:
		value = strconv.FormatBool(any.(bool))
	case int:
		value = strconv.FormatInt(int64(any.(int)), 10)
	case int8:
		value = strconv.FormatInt(int64(any.(int8)), 10)
	case int16:
		value = strconv.FormatInt(int64(any.(int16)), 10)
	case int32:
		value = strconv.FormatInt(int64(any.(int32)), 10)
	case int64:
		value = strconv.FormatInt(any.(int64), 10)
	case uint:
		value = strconv.FormatUint(uint64(any.(uint)), 10)
	case uint8:
		value = strconv.FormatUint(uint64(any.(uint8)), 10)
	case uint16:
		value = strconv.FormatUint(uint64(any.(uint16)), 10)
	case uint32:
		value = strconv.FormatUint(uint64(any.(uint32)), 10)
	case uint64:
		value = strconv.FormatUint(any.(uint64), 10)
	case float32:
		value = strconv.FormatFloat(float64(any.(float32)), 'f', 2, 32)
	case float64:
		value = strconv.FormatFloat(any.(float64), 'f', 4, 64)
	case func() string:
		value = anyi()
	case func() bool:
		value = any2string(anyi())
	case reflect.Value:
		if anyi.IsValid() && anyi.CanInterface() {
			value = reflect2string(anyi.Interface())
		} else {
			value = reflect2string(any)
		}
	default:
		value = reflect2string(any)
	}
	return value
}

func structTmap(tp *reflect.Type, value *reflect.Value) map[string]interface{} {
	rnt := make(map[string]interface{})
	for k := 0; k < value.NumField(); k++ {
		field := value.Field(k)
		if field.CanInterface() {
			rnt[(*tp).Field(k).Name] = field.Interface()
		}
	}
	return rnt
}

func struct2map(any interface{}) map[string]interface{} {
	anyt := reflect.TypeOf(any)
	anyv := reflect.ValueOf(any)
	if anyv.IsValid() {
		switch anyt.Kind() {
		case reflect.Map:
			{
				rnt := make(map[string]interface{})
				for _, v := range anyv.MapKeys() {
					rnt[any2string(v)] = anyv.MapIndex(v).Interface()
				}
				return rnt
			}
		case reflect.Struct:
			{
				return structTmap(&anyt, &anyv)
			}
		case reflect.Ptr:
			{
				anyv = anyv.Elem()
				anyt = anyv.Type()
				return structTmap(&anyt, &anyv)
			}
		}
	}
	return nil
}

func anys2strings(anys []interface{}) []string {
	rtns := make([]string, len(anys))
	for index, any := range anys {
		rtns[index] = any2string(any)
	}
	return rtns
}

func replace(source, skey string, args []interface{}) string {
	if len(source) <= 0 {
		return source
	}
	strs := anys2strings(args)
	strslen := len(strs)
	outstr := ""
	for _, key := range strings.Split(source, skey) {
		if len(key) <= 0 {
			continue
		}
		if allNumber(key) == false {
			outstr = outstr + key
		} else {
			index, err := strconv.Atoi(key)
			if err == nil && index > 0 && index <= strslen {
				outstr = outstr + any2string(strs[index-1])
			}
		}
	}
	return outstr
}

func ToString(args ...interface{}) string {
	outstr := ""
	for _, arg := range args {
		outstr = outstr + any2string(arg) + "|"
	}
	return outstr
}

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

func menberName(name string) string {
	if len(name) <= 0 {
		return ""
	}
	keys := []byte(name)
	if len(keys) > 0 && keys[0] >= 'a' && keys[0] <= 'z' {
		keys[0] -= 32
	} else {
		return name
	}
	return string(keys)
}

func GetTime(tm *time.Time) string {
	if tm == nil {
		return time.Now().Format(TIME_FORMATE_UNIX)
	} else {
		return tm.Format(TIME_FORMATE_UNIX)
	}
}
