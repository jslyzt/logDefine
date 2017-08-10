package logDefine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// 导出入口文件
func (file *XmlLogFile) javafmortEntrance(outdir string) {
	modelName := file.Name
	baseName := file.MName
	fileName := fmt.Sprintf("%s/%s.java", outdir, baseName)
	fmt.Printf("save file: %s\n", fileName)

	casestr := ""
	for _, node := range file.Logs {
		casestr = casestr + replace(
			`            if (type.equals("#1#")) {
                return Func.Trans2String(logger, headers, data, #2#.GetType());
            }
`, "#", []interface{}{
				node.Alias,
				fmt.Sprintf("%s_%s", baseName, node.Name),
			})
	}

	basestr := replace(
		`package com.hiscene.common.#1#;

import org.apache.commons.lang.StringUtils;

import java.util.Map;

public class #2# {

    public static byte[] Json2String(org.slf4j.Logger logger, Map<String, String> headers, String data) {
        String type = headers.get("indexName");
        if (StringUtils.isNotBlank(type)) {
#3#
        }
        return new byte[0];
    }
}
`, "#", []interface{}{
			modelName,
			baseName,
			casestr,
		})

	err := ioutil.WriteFile(fileName, []byte(basestr), os.ModePerm)
	if err != nil {
		fmt.Printf("save file %s failed, error %v", fileName, err)
	}
}

// 类导出共用方法
func javafmortSave(outdir, data string, file *XmlLogFile, info *XmlLogStruct) {
	fileName := fmt.Sprintf("%s/%s_%s.java", outdir, file.MName, info.Name)
	fmt.Printf("save file: %s\n", fileName)
	err := ioutil.WriteFile(fileName, []byte(data), os.ModePerm)
	if err != nil {
		fmt.Printf("save file %s failed, error %v", fileName, err)
	}
}

func javaGetMapValue(tp string) string {
	switch tp {
	case "int", "int16", "int32":
		return "int"
	case "int8":
		return "short"
	case "int64", "datetime":
		return "long"
	case "interface{}":
		return "Object"
	case "bool":
		return "boolean"
	case "string":
		return "String"
	}
	return tp
}

func javeGetType(key string, tp int8) string {
	switch tp {
	case UDT_LIST:
		keys := []byte(key)
		return fmt.Sprintf("List<%s>", javaGetMapValue(string(keys[2:])))
	case UDT_PLIST:
		keys := []byte(key)
		return fmt.Sprintf("List<%s>", javaGetMapValue(string(keys[3:])))
	case UDT_MAP:
		keys := []byte(key)
		index := strings.Index(key, "]")
		if index >= 0 {
			return fmt.Sprintf("Map<%s, %s>", javaGetMapValue(string(keys[4:index])), javaGetMapValue(string(keys[index+1:])))
		}
	case UDT_PMAP:
		keys := []byte(key)
		index := strings.Index(key, "]")
		if index >= 0 {
			return fmt.Sprintf("Map<%s, %s>", javaGetMapValue(string(keys[4:index])), javaGetMapValue(string(keys[index+2:])))
		}
	}
	return key
}

func javaGetNodeType(node *XmlLogNode) (string, string) {
	if node != nil {
		switch node.Type {
		case T_INT:
			return "int", ""
		case T_FLOAT:
			return "float", ""
		case T_DOUBLE:
			return "double", ""
		case T_STRING:
			return "String", ""
		case T_DATETIME:
			return "long", ""
		case T_BOOL:
			return "boolean", ""
		case T_SHORT:
			return "short", ""
		case T_LONG:
			return "long", ""
		case T_USERDEF:
			switch node.UDType {
			case UDT_LIST, UDT_PLIST:
				return javeGetType(node.SType, node.UDType), "java.util.List"
			case UDT_MAP, UDT_PMAP:
				return javeGetType(node.SType, node.UDType), "java.util.Map"
			default:
				return node.SType, ""
			}
		}
	}
	return "", ""
}

func javafmortMembers(info *XmlLogStruct) (string, string) {
	var bfimport bytes.Buffer
	var bfmember bytes.Buffer
	for _, node := range info.Nodes {
		stype, simport := javaGetNodeType(&node)
		if len(simport) > 0 {
			bfimport.WriteString(fmt.Sprintf("import %s;\n", simport))
		}
		if len(stype) > 0 {
			bfmember.WriteString(fmt.Sprintf("    private %s %s;\n", stype, node.Xname))
		}
	}
	return bfimport.String(), bfmember.String()
}

func javafmortTostring(info *XmlLogStruct) string {
	var bfstring bytes.Buffer
	index := 0
	for _, node := range info.Nodes {
		if index == 0 {
			bfstring.WriteString(fmt.Sprintf("        writer.append(%s)", node.Xname))
		} else {
			bfstring.WriteString(fmt.Sprintf("\n                .append(%s)", node.Xname))
		}
		index = index + 1
	}
	return bfstring.String() + ";"
}

// 导出struct文件
func javafmortStruct(outdir string, file *XmlLogFile, info *XmlLogStruct) {
	imports, members := javafmortMembers(info)
	tostring := javafmortTostring(info)
	javafmortSave(outdir, replace(
		`package com.hiscene.common.#1#;

#2#
public class #3#_#4# {
#5#

    public String toString() {
        FStringWriter writer = new FStringWriter();
#6#
        return writer.toString();
    }
}
`, "#", []interface{}{
			file.Name,
			imports,
			file.MName,
			info.Name,
			members,
			tostring,
		}), file, info)
}

// 导出log文件
func javafmortLog(outdir string, file *XmlLogFile, info *XmlLogStruct) {
	imports, members := javafmortMembers(info)
	tostring := javafmortTostring(info)
	javafmortSave(outdir, replace(`
package com.hiscene.common.#1#;

import com.google.gson.reflect.TypeToken;
import java.lang.reflect.Type;
#2#
public class #3#_#4# {
#5#
    private static Type typeToken = new TypeToken<#3#_#4#>(){}.getType();

    public static Type GetType() {
        return typeToken;
    }

    public String toString() {
        FStringWriter writer = new FStringWriter();
#6#
        return writer.toString();
    }
}
`, "#", []interface{}{
		file.Name,
		imports,
		file.MName,
		info.Name,
		members,
		tostring,
	}), file, info)
}

// 导出 java
func (file *XmlLogFile) exportJava(outdir string) bool {
	file.javafmortEntrance(outdir)

	for _, node := range file.Stus {
		javafmortStruct(outdir, file, &node)
	}
	for _, node := range file.Logs {
		javafmortLog(outdir, file, &node)
	}
	return true
}
