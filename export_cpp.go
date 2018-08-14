package logdefine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/djimenez/iconv-go"
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
		return "common::Object"
	case "bool":
		return "bool"
	case "string":
		return "std::string"
	}
	return tp
}

func cppGetType(key string, tp int8) string {
	switch tp {
	case UDTlist:
		keys := []byte(key)
		return fmt.Sprintf("std::list<%s>", cppGetMapValue(string(keys[2:])))
	case UDTplist:
		keys := []byte(key)
		return fmt.Sprintf("std::list<%s*>", cppGetMapValue(string(keys[3:])))
	case UDTmap:
		keys := []byte(key)
		index := strings.Index(key, "]")
		if index >= 0 {
			return fmt.Sprintf("std::map<%s, %s>", cppGetMapValue(string(keys[4:index])), cppGetMapValue(string(keys[index+1:])))
		}
	case UDTpmap:
		keys := []byte(key)
		index := strings.Index(key, "]")
		if index >= 0 {
			return fmt.Sprintf("std::map<%s, %s*>", cppGetMapValue(string(keys[4:index])), cppGetMapValue(string(keys[index+2:])))
		}
	}
	return key
}

func cppgetNodeType(node *XMLLogNode) (string, string) {
	if node != nil {
		switch node.Type {
		case TInt:
			return "int32_t", ""
		case TFloat:
			return "float", ""
		case TDouble:
			return "double", ""
		case TString:
			return "std::string", ""
		case TDateTime:
			return "int64_t", ""
		case TBool:
			return "bool", ""
		case TShort:
			return "short_t", ""
		case TLong:
			return "long_t", ""
		case TUserDef:
			switch node.UDType {
			case UDTlist, UDTplist:
				return cppGetType(node.SType, node.UDType), "list"
			case UDTmap, UDTpmap:
				return cppGetType(node.SType, node.UDType), "map"
			default:
				return node.SType, ""
			}
		}
	}
	return "std::string", ""
}

// 序列化结构
func cppfmortStruct(file *XMLLogFile, info *XMLLogStruct, incs *map[string]bool) string {
	complexTypes := make(map[string]string)
	var buffer, complexs bytes.Buffer
	for index, node := range info.Nodes {
		memtp, include := cppgetNodeType(&node)
		memtp = strings.Replace(memtp, file.MName+"_", "", 100)
		if strings.Index(memtp, ",") >= 0 {
			comptype, ok := complexTypes[memtp]
			if ok {
				memtp = comptype
			} else {
				smemtp := fmt.Sprintf("COMPLEX_%s", strings.ToUpper(node.Xname))
				complexTypes[memtp] = smemtp
				memtp = smemtp
			}
		}
		//strmember := fmt.Sprintf("\t%-20s %s;", memtp, node.Xname)
		strmember := fmt.Sprintf("\tLOGMEMDEF(%s, %s, %d);", memtp, node.Xname, index)
		if len(node.Desc) > 0 {
			buffer.WriteString(fmt.Sprintf("%-50s // desc:%s", strmember, node.Desc))
		} else {
			buffer.WriteString(strmember)
		}
		if incs != nil && len(include) > 0 {
			if _, ok := (*incs)[include]; ok == false {
				(*incs)[include] = true
			}
		}
		buffer.WriteString("\n")
	}
	for k, v := range complexTypes {
		complexs.WriteString(fmt.Sprintf("\ttypedef %s %s;\n", k, v))
	}
	return replace(`// --------------------------------------------------------------------
// #1# 结构定义
// #3#
struct #1# {	// version #2#
#5#
#4#
    void logExport(std::stringstream& stream) const;                                        // 序列化方法
    void logEntrance(const std::string& str, size_t size = 0, size_t* index = nullptr);     // 反序列化方法
};
`, "#", []interface{}{
		info.Name,
		info.Version,
		info.Desc,
		buffer.String(),
		complexs.String(),
	})
}

// 结构序列化方法
func cppfmortStrfunc(file *XMLLogFile, info *XMLLogStruct) string {
	var buffer bytes.Buffer
	for _, node := range info.Nodes {
		buffer.WriteString(fmt.Sprintf("%s, ", node.Xname))
	}
	nodestr := buffer.String()
	if len(nodestr) > 0 {
		nodestr = strings.Trim(nodestr, ", ")
	}
	return replace(`// --------------------------------------------------------------------
// #1#
void #1#::logExport(std::stringstream& stream) const {
	LOG(stream, #2#);
}

void #1#::logEntrance(const std::string& str, size_t size, size_t* index) {
    size_t sindex = 0;
    if (index == nullptr) {
        index = &sindex;
    }
    if (size <= 0) {
        size = str.length();
    }
	ENTRANCE(str, size, *index, #2#);
}
`, "#", []interface{}{
		info.Name,
		nodestr,
	})
}

// c++文件序列化方法
func cppfmortLogfile(file *XMLLogFile, includes []string) (string, string) {
	incs := make(map[string]bool)
	var bufStuH, bufStuF bytes.Buffer
	for _, node := range file.Stus {
		bufStuH.WriteString(cppfmortStruct(file, &node, &incs))
		bufStuF.WriteString(cppfmortStrfunc(file, &node))
	}
	var bufLogH, bufLogF bytes.Buffer
	for _, strnode := range file.Logs {
		bufLogH.WriteString(cppfmortStruct(file, &strnode, &incs))
		bufLogF.WriteString(cppfmortStrfunc(file, &strnode))
	}
	strincs := ""
	for k := range incs {
		strincs = strincs + fmt.Sprintf("#include <%s>\n", k)
	}
	for _, v := range includes {
		strincs = strincs + fmt.Sprintf("#include \"%s\"\n", v)
	}
	return fmt.Sprintf(
			`#pragma once

#include <stdint.h>
#include <string>
#include <sstream>
%s

namespace %s{

%s

%s
}
`, strincs, file.Name, bufStuH.String(), bufLogH.String()),
		fmt.Sprintf(
			`#include "logdef_%s.h"

namespace %s{

%s

%s
}
`, file.Name, file.Name, bufStuF.String(), bufLogF.String())
}

func saveFile(sdata, name, charset string) {
	fmt.Printf("save file: %s\n", name)
	var outstr bytes.Buffer
	limit, begin, hasxg := 120, 0, 0
	data := []byte(sdata)
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

	if charset != "utf-8" {
		converter, err := iconv.NewConverter("utf-8", charset)
		if err == nil {
			ostr, err := converter.ConvertString(outstr.String())
			if err == nil {
				outstr.Reset()
				outstr.WriteString(ostr)
			}
		}

	}

	os.MkdirAll(filepath.Dir(name), os.ModePerm)
	err := ioutil.WriteFile(name, []byte(outstr.String()), os.ModePerm)
	if err != nil {
		fmt.Printf("save file %s failed, error %v", name, err)
	}
	runCmd("astyle", "--style=java", "-SMpHnUoOY", "-k1W1", name)
}

// 导出 golang
func (file *XMLLogFile) exportCpp(outdir string, appends map[string]interface{}) bool {
	charset := ""
	any, ok := appends["charset"]
	if ok {
		charset = any.(string)
	}

	includes := make([]string, 0)
	any, ok = appends["includes"]
	if ok {
		includes = any.([]string)
	}

	hdata, pdata := cppfmortLogfile(file, includes)
	if len(hdata) > 0 {
		saveFile(hdata, fmt.Sprintf("%s/logdef_%s.h", outdir, file.Name), charset)
	}
	if len(pdata) > 0 {
		saveFile(pdata, fmt.Sprintf("%s/logdef_%s.cpp", outdir, file.Name), charset)
	}
	return true
}
