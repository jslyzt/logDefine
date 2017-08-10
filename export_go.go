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
		case T_BOOL:
			return "bool"
		case T_SHORT:
			return "int8"
		case T_LONG:
			return "int64"
		case T_USERDEF:
			return node.SType
		}
	}
	return "string"
}

// 序列化结构
func gofmortStruct(file *XmlLogFile, info *XmlLogStruct) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		buffer.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`//",
			node.Name,
			gogetNodeType(&node),
			node.Xname))
		/*
			defstr := any2string(node.Defvalue)
			if len(defstr) > 0 {
				buffer.WriteString(fmt.Sprintf(" default: %s", defstr))
			}
		*/
		if len(node.Desc) > 0 {
			buffer.WriteString(fmt.Sprintf(" desc: %s", node.Desc))
		}
		buffer.WriteString("\n")
	}
	return replace(`
// #1# #2#结构定义
// #4#
type #6#_#2# struct {	// version #3#
#5#
}
`, "#", []interface{}{
		file.Name,
		info.Name,
		info.Version,
		info.Desc,
		buffer.String(),
		file.MName,
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
func gofmortstrFuncStruct() string {
	return `
// #1# #2#序列化方法
func (node #4#_#2#) ToString() string {
	return logDefine.ToString(#3#)
}
`
}

func gofmortstrFuncLog() string {
	return `
// #1# #2#序列化方法
func (node #4#_#2#) ToString() string {
	return node.ToAString(node.GetAppend())
}
func (node #4#_#2#) ToAString(arr []interface{}) string {
	strlog := logDefine.ToString(#3#)
	if len(arr) > 0 {
		strlog = logDefine.ToString(arr...) + strlog
	}
	return strlog
}
func (node #4#_#2#) GetAlias() string {
	return "#5#"
}
func (node #4#_#2#) GetAppend() []interface{} {
	return []interface{}{
		node.GetAlias(),
		logDefine.GetTime(nil),
	}
}
`
}

func gofmort2String(file *XmlLogFile, info *XmlLogStruct, bstu bool) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		buffer.WriteString(fmt.Sprintf("node.%s, ", node.Name))
	}
	nodestr := buffer.String()
	if len(nodestr) > 0 {
		nodestr = strings.Trim(nodestr, ", ")
	}
	var fmortstr string
	if bstu == true {
		fmortstr = gofmortstrFuncStruct()
	} else {
		fmortstr = gofmortstrFuncLog()
	}
	return replace(fmortstr, "#", []interface{}{
		file.Name,
		info.Name,
		nodestr,
		file.MName,
		info.Alias,
	})
}

// 序列化json
func gofmort2Json(file *XmlLogFile, info *XmlLogStruct) string {
	return replace(`
// #1# #2#序列化方法
func (node #3#_#2#) ToJson() string {
	data, err := json.Marshal(node)
	if err == nil {
		return string(data)
	}
	return ""
}
`, "#", []interface{}{
		file.Name,
		info.Name,
		file.MName,
	})
}

// 结构序列化方法
func gofmortStrfunc(file *XmlLogFile, info *XmlLogStruct, bstu bool) string {
	return gofmort2String(file, info, bstu) + gofmort2Json(file, info)
}

// go文件序列化方法
func gofmortLogfile(file *XmlLogFile) string {
	var bufStu bytes.Buffer
	for _, node := range file.Stus {
		bufStu.WriteString(gofmortStruct(file, &node))
		bufStu.WriteString(gofmortStrfunc(file, &node, true))
	}
	var bufLog bytes.Buffer
	for _, node := range file.Logs {
		bufLog.WriteString(gofmortStruct(file, &node))
		bufLog.WriteString(gofmortStrfunc(file, &node, false))
	}
	return fmt.Sprintf(`
package %s

%s

%s

%s
`, file.Name, gofmortDeffunc(), bufStu.String(), bufLog.String())
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
