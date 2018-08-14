package logdefine

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// 工具函数
func (node *XMLLogNode) alsisStuType(file *XMLLogFile) {
	node.Type = TUserDef
	node.UDType = UDTnone
	if len(file.StuMp) > 0 {
		keys := []byte(node.SType)
		lkey := len(keys)
		tkey := node.SType
		if lkey > 2 && keys[0] == '[' && keys[1] == ']' { // 数组
			tkey = string(keys[2:])
			node.UDType = UDTlist
			if lkey > 3 && keys[2] == '*' {
				tkey = string(keys[3:])
				node.UDType = UDTplist
			}
		} else if lkey > 4 && string(keys[:4]) == "map[" { // map
			mindex := strings.Index(node.SType, "]")
			if mindex >= 0 {
				tkey = string(keys[mindex+1:])
				node.UDType = UDTmap
				if len(tkey) > 0 && keys[mindex+2] == '*' {
					tkey = string(keys[mindex+2:])
					node.UDType = UDTpmap
				}
			}
		}
		if len(tkey) > 0 {
			if _, ok := file.StuMp[tkey]; ok {
				node.SType = strings.Replace(node.SType, tkey, fmt.Sprintf("%s_%s", file.MName, tkey), 100)
			}
		}
	}
}

func (node *XMLLogNode) analysis(file *XMLLogFile) {
	switch node.SType {
	case "string":
		node.Type = TString
	case "int":
		node.Type = TInt
	case "float":
		node.Type = TFloat
	case "double":
		node.Type = TDouble
	case "datetime":
		node.Type = TDateTime
	case "bool":
		node.Type = TBool
	case "short":
		node.Type = TShort
	case "long":
		node.Type = TLong
	default:
		node.alsisStuType(file)
	}

	node.Name = menberName(node.Xname)
}

func (info *XMLLogStruct) analysis(file *XMLLogFile) {
	for index := range info.Nodes {
		info.Nodes[index].analysis(file)
	}
	if len(info.Alias) <= 0 {
		info.Alias = strings.ToLower(info.Name)
	}
}

// 分析
func (file *XMLLogFile) analysis() error {
	if file == nil || len(file.file) <= 0 {
		return errors.New("file is nil")
	}
	fd, err := os.Open(file.file)
	if err != nil {
		return err
	}
	defer fd.Close()
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(data, file)
	if err != nil {
		return err
	}
	file.MName = menberName(file.Name)
	for index := range file.Stus {
		node := &file.Stus[index]
		node.analysis(file)
		file.StuMp[node.Name] = node
	}
	for index := range file.Logs {
		file.Logs[index].analysis(file)
	}
	return nil
}

// Export 导出
func (file *XMLLogFile) Export(types []int8, outdir string, appends map[string]interface{}) {
	if file != nil {
		for _, ntp := range types {
			switch ntp {
			case ETgo:
				file.exportGo(outdir)
			case ETcpp:
				file.exportCpp(outdir, appends)
			case ETjava:
				file.exportJava(outdir)
			default:
				fmt.Printf("no support export type: %d\n", ntp)
			}
		}
	}
}

// AnalysisFile 分析文件
func AnalysisFile(file string) *XMLLogFile {
	xmllog := &XMLLogFile{
		file:  file,
		Stus:  make(XMLLogStructs, 0),
		Logs:  make(XMLLogStructs, 0),
		StuMp: make(XMLLogStrMap, 0),
	}
	err := xmllog.analysis()
	if err != nil {
		fmt.Printf("analysis file: %s error, error is %v\n", file, err)
		return nil
	}
	return xmllog
}
