package logdefine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func gogetNodeType(file *XMLLogFile, node *XMLLogNode) string {
	if node != nil {
		switch node.Type {
		case TInt:
			return "int"
		case TFloat:
			return "float32"
		case TDouble:
			return "float64"
		case TString:
			return "string"
		case TDateTime:
			return "time.Time"
		case TBool:
			return "bool"
		case TShort:
			return "int8"
		case TLong:
			return "int64"
		case TUserDef:
			return strings.Replace(node.SType, fmt.Sprintf("%v_", file.MName), "", 100)
		}
	}
	return "string"
}

// 序列化结构
func gofmortStruct(file *XMLLogFile, info *XMLLogStruct) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		buffer.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`//",
			node.Name,
			gogetNodeType(file, &node),
			node.Xname))
		if len(node.Desc) > 0 {
			buffer.WriteString(fmt.Sprintf(" desc: %s", node.Desc))
		}
		buffer.WriteString("\n")
	}
	return replace(`
// --------------------------------------------------------------------
// #1# #2#结构定义

// #2# #4#
type #2# struct {	// version #3#
#5#
}
`, "#", []interface{}{
		file.Name,
		info.UName,
		info.Version,
		info.Desc,
		buffer.String(),
	})
}

// 通用序列化方法
func gofmortDeffunc() string {
	return `
	import(
		logDefine "github.com/jslyzt/logdefine"
	)
	`
}

// 序列化string
func gofmortstrFuncStruct() string {
	return `// ToString 转换string
func (node *#1#) ToString() string {
	return logDefine.ToString(node)
}
`
}

func gofmortstrFuncLog() string {
	return `// ToString 转换string
func (node *#1#) ToString() string {
	return logDefine.ToString(node.GetAlias(), logDefine.GetTime(nil), *node)
}

// GetAlias 获取别名
func (node *#1#) GetAlias() string {
	return "#2#"
}
`
}

func gofmort2String(file *XMLLogFile, info *XMLLogStruct, bstu bool) string {
	var fmortstr string
	if bstu == true {
		fmortstr = gofmortstrFuncStruct()
	} else {
		fmortstr = gofmortstrFuncLog()
	}
	return replace(fmortstr, "#", []interface{}{
		info.UName,
		info.Alias,
	})
}

// 序列化json
func gofmort2Json(file *XMLLogFile, info *XMLLogStruct) string {
	return replace(`// ToJSON 转换json
func (node *#1#) ToJSON() string {
	data, err := json.Marshal(*node)
	if err == nil {
		return string(data)
	}
	return ""
}
`, "#", []interface{}{
		info.UName,
	})
}

// 反序列化string
func gofmortstrFuncUStruct() string {
	return `// FromString string初始化
func (node *#1#) FromString(data []byte, index int) int  {
	return logDefine.FromString(data, index, node)
}
`
}

func gofmortstrFuncULog() string {
	return `// FromString string初始化
func (node *#1#) FromString(data []byte, index int) (size int, alias, stime string) {
	size = logDefine.FromString(data, index, &alias, &stime, node)
	return
}
`
}

func gofmortFString(file *XMLLogFile, info *XMLLogStruct, bstu bool) string {
	var fmortstr string
	if bstu == true {
		fmortstr = gofmortstrFuncUStruct()
	} else {
		fmortstr = gofmortstrFuncULog()
	}
	return replace(fmortstr, "#", []interface{}{
		info.UName,
	})
}

// 反序列化json
func gofmortFJson(file *XMLLogFile, info *XMLLogStruct) string {
	return replace(`// FromJSON json初始化
func (node *#1#) FromJSON(data []byte)  {
	json.Unmarshal(data, node)
}
`, "#", []interface{}{
		info.UName,
	})
}

// 结构序列化方法
func gofmortStrfunc(file *XMLLogFile, info *XMLLogStruct, bstu bool) string {
	return replace(`
// #1# #2#序列化方法

#3#

// #1# #2#反序列化方法

#4#
`, "#", []interface{}{
		file.Name,
		info.UName,
		gofmort2String(file, info, bstu) + gofmort2Json(file, info),
		gofmortFString(file, info, bstu) + gofmortFJson(file, info),
	})
}

// go文件序列化方法
func gofmortLogfile(file *XMLLogFile) string {
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
func (file *XMLLogFile) exportGo(outdir string) bool {
	fileName := fmt.Sprintf("%s/%s/%s.go", outdir, file.Name, file.Name)
	fmt.Printf("save file: %s\n", fileName)

	os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	err := ioutil.WriteFile(fileName, []byte(gofmortLogfile(file)), os.ModePerm)
	if err != nil {
		fmt.Printf("save file %s failed, error %v", fileName, err)
		return false
	}
	//runCmd("goreturns", fmt.Sprintf("-w %s", fileName))
	runCmd("goreturns", "-w", fileName)
	return true
}
