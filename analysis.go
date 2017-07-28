package logDefine

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

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
	for _, node := range file.Logs {
		file.maps[node.Name] = &node
	}
	return nil
}

// 导出
func (file *XmlLogFile) export(types []int8) {
	if file != nil {
		for _, ntp := range types {
			switch ntp {
			case ET_GO:
				file.export_go()
			case ET_CPP:
				file.export_cpp()
			default:
				fmt.Printf("no support export type: %d\n", ntp)
			}
		}
	}
}

// 分析文件
func AnalysisFile(file string) *XmlLogFile {
	xmllog := &XmlLogFile{
		file: file,
		Logs: make(XmlLogStructs, 0),
		maps: make(XmlLogStrMap),
	}
	err := xmllog.analysis()
	if err != nil {
		fmt.Printf("analysis file: %s error, error is %v", file, err)
		return nil
	}
	return xmllog
}
