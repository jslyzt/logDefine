package logDefine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func gogetNodeType(node *XmlLogNode) string {
	if node != nil {
		switch node.Type {
		case T_INT:
			return "int"
		case T_FLOAT:
			return "float32"
		case T_DOUBLE:
			return "float64"
		case T_STRING:
			return "string"
		case T_DATETIME:
			return "time.Time"
		}
	}
	return "string"
}

// 序列化结构
func gofmortStruct(file *XmlLogFile, info *XmlLogStruct) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		buffer.WriteString(fmt.Sprintf("\t%s %s //",
			node.Name,
			gogetNodeType(&node)))
		defstr := any2string(node.Defvalue)
		if len(defstr) > 0 {
			buffer.WriteString(fmt.Sprintf(" default: %s", defstr))
		}
		if len(node.Desc) > 0 {
			buffer.WriteString(fmt.Sprintf(" desc: %s", node.Desc))
		}
		buffer.WriteString("\n")
	}
	return replace(`
// #1# #2#结构定义
// #4#
type #1#_#2# struct {	// version #3#
#5#
}

`, "#", []interface{}{
		file.Name,
		info.Name,
		info.Version,
		info.Desc,
		buffer.String(),
	})
}

// 通用序列化方法
func gofmortDeffunc() string {
	return `
	import(
		"github.com/jslyzt/logDefine"
	)
	`
}

// 序列化string
func gofmort2String(file *XmlLogFile, info *XmlLogStruct) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		buffer.WriteString(fmt.Sprintf("node.%s, ", node.Name))
	}
	nodestr := buffer.String()
	if len(nodestr) > 0 {
		nodestr = strings.Trim(nodestr, ", ")
	}
	return replace(`
// #1# #2#序列化方法
func (node *#1#_#2#) ToString() string {
	return logDefine.ToString(#3#)
}
`, "#", []interface{}{
		file.Name,
		info.Name,
		nodestr,
	})
}

// 序列化json
func gofmort2Json(file *XmlLogFile, info *XmlLogStruct) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		buffer.WriteString(fmt.Sprintf("\"%s\": node.%s,\n", node.Name, node.Name))
	}
	return replace(`
// #1# #2#序列化方法
func (node *#1#_#2#) ToJson() string {
	return json.Marshal(map[string]interface{}{
#3#
	})
}
`, "#", []interface{}{
		file.Name,
		info.Name,
		buffer.String(),
	})
}

// 结构序列化方法
func gofmortStrfunc(file *XmlLogFile, info *XmlLogStruct) string {
	return gofmort2String(file, info) + gofmort2Json(file, info)
}

// go文件序列化方法
func gofmortLogfile(file *XmlLogFile) string {
	var buffer bytes.Buffer
	for _, strnode := range file.Logs {
		buffer.WriteString(gofmortStruct(file, &strnode))
		buffer.WriteString(gofmortStrfunc(file, &strnode))
	}
	return fmt.Sprintf(`
package %s

%s

%s
`, file.Name, gofmortDeffunc(), buffer.String())
}

// 导出 golang
func (file *XmlLogFile) exportGo(outdir string) bool {
	fileName := fmt.Sprintf("%s/logdef_%s.go", outdir, file.Name)
	fmt.Printf("save file: %s\n", fileName)
	err := ioutil.WriteFile(fileName, []byte(gofmortLogfile(file)), os.ModePerm)
	if err != nil {
		fmt.Printf("save file %s failed, error %v", fileName, err)
		return false
	}
	//runCmd("goreturns", fmt.Sprintf("-w %s", fileName))
	runCmd("goreturns", "-w", fileName)
	return true
}
