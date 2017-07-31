package logDefine

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// 工具函数
func (node *XmlLogNode) analysis() {
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
		node.Type = T_USERDEF
	}

	node.Name = menberName(node.Xname)
}

func (info *XmlLogStruct) analysis() {
	for index := range info.Nodes {
		info.Nodes[index].analysis()
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
	for index := range file.Stus {
		file.Stus[index].analysis()
	}
	for index := range file.Logs {
		file.Logs[index].analysis()
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
		file: file,
		Logs: make(XmlLogStructs, 0),
	}
	err := xmllog.analysis()
	if err != nil {
		fmt.Printf("analysis file: %s error, error is %v", file, err)
		return nil
	}
	return xmllog
}
