package logDefine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func cppgetNodeType(node *XmlLogNode) string {
	if node != nil {
		switch node.Type {
		case T_INT:
			return "int32_t"
		case T_FLOAT:
			return "float_t"
		case T_DOUBLE:
			return "double_t"
		case T_STRING:
			return "std::string"
		case T_DATETIME:
			return "uint64_t"
		}
	}
	return "std::string"
}

// 序列化结构
func cppfmortStruct(file *XmlLogFile, info *XmlLogStruct) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		buffer.WriteString(fmt.Sprintf("\t%s %s; //",
			cppgetNodeType(&node),
			node.Name))
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
struct #1#_#2# {	// version #3#
#5#
};

`, "#", []interface{}{
		file.Name,
		info.Name,
		info.Version,
		info.Desc,
		buffer.String(),
	})
}

// 结构序列化方法
func cppfmortStrfunc(file *XmlLogFile, info *XmlLogStruct) string {
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
static void logExport(const #1#_#2#& node, std::string& str) {
	str.clear();
	std::stringstream stream;
	commom::log(stream, #3#);
	str = stream.str();
}
`, "#", []interface{}{
		file.Name,
		info.Name,
		nodestr,
	})
}

// c++文件序列化方法
func cppfmortLogfile(file *XmlLogFile) string {
	var buffer bytes.Buffer
	for _, strnode := range file.Logs {
		buffer.WriteString(cppfmortStruct(file, &strnode))
		buffer.WriteString(cppfmortStrfunc(file, &strnode))
	}
	return fmt.Sprintf(
		`#pragma once

#include <stdint.h>
#include <string>
#include <sstream>

#include "logDef.h"

namespace %s{

%s
}
`, file.Name, buffer.String())
}

// 导出 golang
func (file *XmlLogFile) exportCpp(outdir string) bool {
	fileName := fmt.Sprintf("%s/logdef_%s.h", outdir, file.Name)
	fmt.Printf("save file: %s\n", fileName)
	err := ioutil.WriteFile(fileName, []byte(cppfmortLogfile(file)), os.ModePerm)
	if err != nil {
		fmt.Printf("save file %s failed, error %v", fileName, err)
		return false
	}
	runCmd("astyle", "--style=java", "-SMpHnUoOY", "-k1W1", fileName)
	return true
}
