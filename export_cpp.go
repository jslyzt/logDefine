package logDefine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func cppGetMapValue(tp string) string {
	switch tp {
	case "int", "int32":
		return "int32_t"
	case "int8":
		return "int8_t"
	case "int16":
		return "int16_t"
	case "int64", "datetime":
		return "int64_t"
	case "interface{}":
		return "commom::Object"
	case "bool":
		return "bool"
	case "string":
		return "std::string"
	}
	return tp
}

func cppGetType(key string, tp int8) string {
	switch tp {
	case UDT_LIST:
		keys := []byte(key)
		return fmt.Sprintf("std::list<%s>", cppGetMapValue(string(keys[2:])))
	case UDT_PLIST:
		keys := []byte(key)
		return fmt.Sprintf("std::list<%s*>", cppGetMapValue(string(keys[3:])))
	case UDT_MAP:
		keys := []byte(key)
		index := strings.Index(key, "]")
		if index >= 0 {
			return fmt.Sprintf("std::map<%s, %s>", cppGetMapValue(string(keys[4:index])), cppGetMapValue(string(keys[index+1:])))
		}
	case UDT_PMAP:
		keys := []byte(key)
		index := strings.Index(key, "]")
		if index >= 0 {
			return fmt.Sprintf("std::map<%s, %s*>", cppGetMapValue(string(keys[4:index])), cppGetMapValue(string(keys[index+2:])))
		}
	}
	return key
}

func cppgetNodeType(node *XmlLogNode) (string, string) {
	if node != nil {
		switch node.Type {
		case T_INT:
			return "int32_t", ""
		case T_FLOAT:
			return "float_t", ""
		case T_DOUBLE:
			return "double_t", ""
		case T_STRING:
			return "std::string", ""
		case T_DATETIME:
			return "int64_t", ""
		case T_BOOL:
			return "bool", ""
		case T_SHORT:
			return "short_t", ""
		case T_LONG:
			return "long_t", ""
		case T_USERDEF:
			switch node.UDType {
			case UDT_LIST, UDT_PLIST:
				return cppGetType(node.SType, node.UDType), "list"
			case UDT_MAP, UDT_PMAP:
				return cppGetType(node.SType, node.UDType), "map"
			default:
				return node.SType, ""
			}
		}
	}
	return "std::string", ""
}

// 序列化结构
func cppfmortStruct(file *XmlLogFile, info *XmlLogStruct, incs *map[string]bool) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		memtp, include := cppgetNodeType(&node)
		memtp = strings.Replace(memtp, file.MName+"_", "", 100)
		buffer.WriteString(fmt.Sprintf("\t%s %s; //", memtp, node.Xname))
		if len(node.Desc) > 0 {
			buffer.WriteString(fmt.Sprintf(" desc: %s", node.Desc))
		}
		if incs != nil && len(include) > 0 {
			if _, ok := (*incs)[include]; ok == false {
				(*incs)[include] = true
			}
		}
		buffer.WriteString("\n")
	}
	return replace(`// --------------------------------------------------------------------
// #1# 结构定义
// #3#
struct #1# {	// version #2#
#4#
};

`, "#", []interface{}{
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
	return replace(`// #1# 序列化方法
static void logExport(const #1#& node, std::string& str) {
	str.clear();
	std::stringstream stream;
	commom::log(stream, #2#);
	str = stream.str();
}

// #1# 反序列化方法
static void logEntrance(const std::string& str, #1#& node) {
	size_t index = 0;
	commom::entrance(str, str.length(), index, #2#);
}
`, "#", []interface{}{
		info.Name,
		nodestr,
	})
}

// c++文件序列化方法
func cppfmortLogfile(file *XmlLogFile) string {
	incs := make(map[string]bool)
	var bufStu bytes.Buffer
	for _, node := range file.Stus {
		bufStu.WriteString(cppfmortStruct(file, &node, &incs))
		bufStu.WriteString(cppfmortStrfunc(file, &node))
	}
	var bufLog bytes.Buffer
	for _, strnode := range file.Logs {
		bufLog.WriteString(cppfmortStruct(file, &strnode, &incs))
		bufLog.WriteString(cppfmortStrfunc(file, &strnode))
	}
	strincs := ""
	for k := range incs {
		strincs = strincs + fmt.Sprintf("#include <%s>\n", k)
	}
	return fmt.Sprintf(
		`#pragma once

#include <stdint.h>
#include <string>
#include <sstream>
%s
#include "logDef.h"

namespace %s{

%s

%s
}
`, strincs, file.Name, bufStu.String(), bufLog.String())
}

// 导出 golang
func (file *XmlLogFile) exportCpp(outdir string) bool {
	fileName := fmt.Sprintf("%s/logdef_%s.h", outdir, file.Name)
	fmt.Printf("save file: %s\n", fileName)

	var outstr bytes.Buffer
	limit, begin, hasxg := 120, 0, 0
	data := []byte(cppfmortLogfile(file))
	for index := 0; index < len(data); index++ {
		if data[index] == '/' && data[index+1] == '/' && index+1 < len(data) {
			hasxg = 1
		}
		if data[index] == '\n' || (index-begin > limit && data[index] == ' ' && hasxg == 0) {
			outstr.Write(data[begin:index])
			begin = index
			hasxg = 0
			if data[index] != '\n' {
				outstr.WriteByte('\n')
			}
		}
	}

	err := ioutil.WriteFile(fileName, []byte(outstr.String()), os.ModePerm)
	if err != nil {
		fmt.Printf("save file %s failed, error %v", fileName, err)
		return false
	}
	runCmd("astyle", "--style=java", "-SMpHnUoOY", "-k1W1", fileName)
	return true
}
