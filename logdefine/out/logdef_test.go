package logger

import (
	"fmt"
	"reflect"
	"testing"
)

func createLogData() *Logger_LogData {
	images := make([]*Logger_ResultImage, 0)
	images = append(images,
		&Logger_ResultImage{
			ImageID:  "imageid1",
			ServerIP: "127.0.0.1",
			ImageURL: "image://1/xxxx",
			Score:    1.99001,
		},
		&Logger_ResultImage{
			ImageID:  "imageid2",
			ServerIP: "127.0.0.2",
			ImageURL: "image://2/xxxx",
			Score:    1.99002,
		})

	logdata := &Logger_LogData{
		RequestID:   "xxxxxxxxx",
		Token:       "123456",
		Latitude:    1.0009,
		Longitude:   2.9991,
		Collection:  "collection",
		Number:      1,
		ClientIP:    "127.0.0.1",
		Image:       "image://xxxx",
		ResultImage: images,
		Result: map[string]interface{}{
			"success": 0,
			"error":   "none",
		},
		CreateTime:             "2017-08-01T10:15:12+08:00",
		Timeconst:              0.198,
		AppKey:                 "testkey1",
		Appname:                "testname",
		Useragent:              "testagent",
		Version:                "1.0.1",
		RecognizeTimeConsuming: 99.129,
	}
	return logdata
}

func Test_fromString(t *testing.T) {
	sdkReco := Logger_sdkReco{
		Business: *createLogData(),
		OauthInfo: map[string]interface{}{
			"success": 0,
			"desc":    "it is a test",
		},
	}
	data := sdkReco.ToString()

	sdkReco2 := Logger_sdkReco{}
	var alias, stime string
	datas := []byte(data)
	sdkReco2.FromAString(datas, 0, &alias, &stime)
	sdkReco2.FromString(datas, 0)
	fmt.Println("alias: ", alias, ", stime: ", stime, ", sdkReco: ", sdkReco2)
}

type TStruct struct {
	Name string
}

func Test_struct(t *testing.T) {
	node := TStruct{
		Name: "name",
	}
	value := reflect.ValueOf(&node)
	value = reflect.Indirect(value)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		field.SetString("111")
	}
	fmt.Println(node, value)
}
