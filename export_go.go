package logDefine

import (
	"fmt"
)

func goget_nodeType(node *XmlLogNode) string {
	if node != nil {
		switch node.Type {
		case T_INT:
			return "int"
		case T_FLOAT:
			return "float32"
		case T_DOUBLE:
			return "float64"
		case T_STRING:
			return "string"
		}
	}
	return "string"
}

// 序列化结构
func gofmort_struct(file *XmlLogFile, info *XmlLogStruct) string {
	nodestr := ""
	for _, node := range info.Nodes {
		nodestr = nodestr + fmt.Sprintf("\t%s %s // default: %s, desc: %s", node.Name,
			goget_nodeType(&node), any2string(node.Defvalue), node.Desc)
	}
	return replace(`
// #1# #2#结构定义
// #4#
type #1#_#2# struct {	// version #3#
#5#
}

`, "#", []interface{}{
		file.Name,
		info.Name,
		info.Version,
		info.Desc,
		nodestr,
	})
}

// 导出 golang
func (file *XmlLogFile) export_go() {

}
