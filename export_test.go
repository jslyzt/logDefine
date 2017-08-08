package logDefine

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestStruct struct {
	Name  string
	Value int
}

func printJson(val interface{}) {
	data, err := json.Marshal(val)
	if err == nil {
		fmt.Println(string(data))
	}
}

func Test_struct2map(t *testing.T) {
	mp1 := struct2map(TestStruct{
		Name:  "1111",
		Value: 1,
	})
	printJson(mp1)
	mp2 := struct2map(map[int]interface{}{
		1: "111",
		2: 2,
	})
	printJson(mp2)
}
