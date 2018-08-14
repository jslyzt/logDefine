package logdefine

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestStruct struct {
	Name  string
	Value int
}

func printJSON(val interface{}) []byte {
	data, err := json.Marshal(val)
	if err == nil {
		fmt.Println(string(data))
		return data
	}
	return nil
}

func Test_struct2map(t *testing.T) {
	mp1 := struct2map(TestStruct{
		Name:  "1111",
		Value: 1,
	})
	mpt1 := &TestStruct{}
	json.Unmarshal(printJSON(mp1), mpt1)

	mp2 := struct2map(map[int]interface{}{
		1: "111",
		2: 2,
	})
	printJSON(mp2)
}

func Test_Base(t *testing.T) {

	sinfo := []byte("1.1|2|{aaa:10;}")
	index := 0
	x := 3.4
	index = bytes2anyptr(sinfo, index, &x, 0)
	y := 0
	index = bytes2anyptr(sinfo, index, &y, 0)
	mp := make(map[string]int)
	index = bytes2anyptr(sinfo, index, &mp, 0)
	fmt.Println(x, y)
}
