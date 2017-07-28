package logDefine

import (
	"strconv"
	"strings"
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

func any2string(any interface{}) string {
	value := ""
	switch anyi := any.(type) {
	case string:
		value = any.(string)
	case bool:
		value = strconv.FormatBool(any.(bool))
	case int64, int, int8, int16, int32:
		value = strconv.FormatInt(any.(int64), 10)
	case uint64, uint, uint8, uint16, uint32:
		value = strconv.FormatUint(any.(uint64), 10)
	case float32:
		value = strconv.FormatFloat(any.(float64), 'f', 2, 32)
	case float64:
		value = strconv.FormatFloat(any.(float64), 'f', 4, 64)
	case func() string:
		value = anyi()
	case func() bool:
		value = any2string(anyi())
	}
	return value
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
			outstr = outstr + skey + key
		} else {
			index, err := strconv.Atoi(key)
			if err == nil && index >= 0 && index < strslen {
				outstr = outstr + any2string(strs[index])
			}
		}
	}
	return outstr
}
