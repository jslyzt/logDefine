package logDefine

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// 工具函数
func (node *XmlLogNode) alsisStuType(file *XmlLogFile) {
	node.Type = T_USERDEF
	if len(file.StuMp) > 0 {
		keys := []byte(node.SType)
		lkey := len(keys)
		tkey := ""
		if lkey > 2 && keys[0] == '[' && keys[1] == ']' { // 数组
			tkey = string(keys[2:])
			if lkey > 3 && keys[2] == '*' {
				tkey = string(keys[3:])
			}
		} else if lkey > 4 && string(keys[:5]) == "map[" { // map
			mindex := strings.Index(node.SType, "]")
			if mindex >= 0 {
				tkey = string(keys[mindex+1:])
				if len(tkey) > 0 && keys[mindex+2] == '*' {
					tkey = string(keys[mindex+2:])
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

func (node *XmlLogNode) analysis(file *XmlLogFile) {
	switch node.SType {
	case "string":
		node.Type = T_STRING
	case "int":
		node.Type = T_INT
	case "float":
		node.Type = T_FLOAT
	case "double":
		node.Type = T_DOUBLE
	case "datetime":
		node.Type = T_DATETIME
	default:
		node.alsisStuType(file)
	}

	node.Name = menberName(node.Xname)
}

func (info *XmlLogStruct) analysis(file *XmlLogFile) {
	for index := range info.Nodes {
		info.Nodes[index].analysis(file)
	}
	if len(info.Alias) <= 0 {
		info.Alias = strings.ToLower(info.Name)
	}
}

// 分析
func (file *XmlLogFile) analysis() error {
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

// 导出
func (file *XmlLogFile) Export(types []int8, outdir string) {
	if file != nil {
		for _, ntp := range types {
			switch ntp {
			case ET_GO:
				file.exportGo(outdir)
			case ET_CPP:
				file.exportCpp(outdir)
			default:
				fmt.Printf("no support export type: %d\n", ntp)
			}
		}
	}
}

func DefaultExport() []int8 {
	return []int8{
		ET_GO,
		ET_CPP,
	}
}

// 分析文件
func AnalysisFile(file string) *XmlLogFile {
	xmllog := &XmlLogFile{
		file:  file,
		Stus:  make(XmlLogStructs, 0),
		Logs:  make(XmlLogStructs, 0),
		StuMp: make(XmlLogStrMap, 0),
	}
	err := xmllog.analysis()
	if err != nil {
		fmt.Printf("analysis file: %s error, error is %v", file, err)
		return nil
	}
	return xmllog
}
