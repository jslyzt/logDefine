package logdefine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 导出入口文件
func (file *XMLLogFile) javafmortEntrance(outdir string) {
	modelName := file.Name
	baseName := file.MName
	fileName := fmt.Sprintf("%s/%s/%s.java", outdir, modelName, baseName)
	fmt.Printf("save file: %s\n", fileName)

	casestr := ""
	for _, node := range file.Logs {
		casestr = casestr + replace(
			`            if (type.equals("#1#")) {
                return Func.Trans2String(logger, headers, data, #2#.GetType());
            }
`, "#", []interface{}{
				node.Alias,
				node.Name,
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

	os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	err := ioutil.WriteFile(fileName, []byte(basestr), os.ModePerm)
	if err != nil {
		fmt.Printf("save file %s failed, error %v", fileName, err)
	}
}

// 类导出共用方法
func javafmortSave(outdir, data string, file *XMLLogFile, info *XMLLogStruct) {
	fileName := fmt.Sprintf("%s/%s/%s.java", outdir, file.Name, info.Name)
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
	case UDTlist:
		keys := []byte(key)
		return fmt.Sprintf("List<%s>", javaGetMapValue(string(keys[2:])))
	case UDTplist:
		keys := []byte(key)
		return fmt.Sprintf("List<%s>", javaGetMapValue(string(keys[3:])))
	case UDTmap:
		keys := []byte(key)
		index := strings.Index(key, "]")
		if index >= 0 {
			return fmt.Sprintf("Map<%s, %s>", javaGetMapValue(string(keys[4:index])), javaGetMapValue(string(keys[index+1:])))
		}
	case UDTpmap:
		keys := []byte(key)
		index := strings.Index(key, "]")
		if index >= 0 {
			return fmt.Sprintf("Map<%s, %s>", javaGetMapValue(string(keys[4:index])), javaGetMapValue(string(keys[index+2:])))
		}
	}
	return key
}

func javaGetNodeType(node *XMLLogNode) (string, string) {
	if node != nil {
		switch node.Type {
		case TInt:
			return "int", ""
		case TFloat:
			return "float", ""
		case TDouble:
			return "double", ""
		case TString:
			return "String", ""
		case TDateTime:
			return "long", ""
		case TBool:
			return "boolean", ""
		case TShort:
			return "short", ""
		case TLong:
			return "long", ""
		case TUserDef:
			switch node.UDType {
			case UDTlist, UDTplist:
				return javeGetType(node.SType, node.UDType), "java.util.List"
			case UDTmap, UDTpmap:
				return javeGetType(node.SType, node.UDType), "java.util.Map"
			default:
				return node.SType, ""
			}
		}
	}
	return "", ""
}

func javafmortMembers(info *XMLLogStruct) (string, string) {
	var bfimport bytes.Buffer
	var bfmember bytes.Buffer
	importmap := make(map[string]bool)
	for _, node := range info.Nodes {
		stype, simport := javaGetNodeType(&node)
		if len(simport) > 0 {
			if _, ok := importmap[simport]; ok == false {
				bfimport.WriteString(fmt.Sprintf("import %s;\n", simport))
				importmap[simport] = true
			}
		}
		if len(stype) > 0 {
			bfmember.WriteString(fmt.Sprintf("    private %s %s;\n", stype, node.Xname))
		}
	}
	return bfimport.String(), bfmember.String()
}

func javafmortTostring(info *XMLLogStruct) string {
	var bfstring bytes.Buffer
	index := 0
	for _, node := range info.Nodes {
		if index == 0 {
			bfstring.WriteString(fmt.Sprintf("        writer.appendSplit(%s)", node.Xname))
		} else {
			bfstring.WriteString(fmt.Sprintf("\n                .appendSplit(%s)", node.Xname))
		}
		index = index + 1
	}
	return bfstring.String() + ";"
}

// 导出struct文件
func javafmortStruct(outdir string, file *XMLLogFile, info *XMLLogStruct) {
	imports, members := javafmortMembers(info)
	tostring := javafmortTostring(info)
	javafmortSave(outdir, replace(
		`package com.hiscene.common.#1#;

#2#
public class #3# {
#4#

    public String toString() {
        FStringWriter writer = new FStringWriter();
#5#
        return writer.toString();
    }
}
`, "#", []interface{}{
			file.Name,
			imports,
			info.Name,
			members,
			tostring,
		}), file, info)
}

// 导出log文件
func javafmortLog(outdir string, file *XMLLogFile, info *XMLLogStruct) {
	imports, members := javafmortMembers(info)
	tostring := javafmortTostring(info)
	javafmortSave(outdir, replace(`
package com.hiscene.common.#1#;

import com.google.gson.reflect.TypeToken;
import java.lang.reflect.Type;
#2#
public class #3# {
#4#
    private static Type typeToken = new TypeToken<#3#>(){}.getType();

    public static Type GetType() {
        return typeToken;
    }

    public String toString() {
        FStringWriter writer = new FStringWriter();
#5#
        return String.format("(%s)", writer.toString());
    }
}
`, "#", []interface{}{
		file.Name,
		imports,
		info.Name,
		members,
		tostring,
	}), file, info)
}

// 导出 java
func (file *XMLLogFile) exportJava(outdir string) bool {
	file.javafmortEntrance(outdir)

	for _, node := range file.Stus {
		javafmortStruct(outdir, file, &node)
	}
	for _, node := range file.Logs {
		javafmortLog(outdir, file, &node)
	}
	return true
}
