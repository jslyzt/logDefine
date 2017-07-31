package logDefine

import "encoding/xml"

const (
	T_INT      int8 = iota // int
	T_FLOAT                // float
	T_DOUBLE               // double
	T_STRING               // string
	T_DATETIME             // 时间日期
	T_USERDEF              // 自定义类型
)
const (
	ET_GO  int8 = iota // 导出go
	ET_CPP             // 导出c++
)

// 节点定义
type XmlLogNode struct {
	Xname    string      `xml:"name,attr"` // 节点名字
	Name     string      // 真正名字
	SType    string      `xml:"type,attr"` // 节点类型 -- xml
	Type     int8        // 节点类型 -- true
	Defvalue interface{} `xml:"defaultvalue,attr"` // 节点默认值
	Desc     string      `xml:"desc,attr"`         // 节点说明
}
type XmlLogNodes []XmlLogNode

// 日志描述定义
type XmlLogStruct struct {
	Name    string      `xml:"name,attr"`    // 日志名字
	Version int16       `xml:"version,attr"` // 版本号
	Desc    string      `xml:"desc,attr"`    // 说明
	Nodes   XmlLogNodes `xml:"entry"`        // 节点列表
}
type XmlLogStructs []XmlLogStruct

//type XmlLogStrMap map[string]*XmlLogStruct

// 日志文件定义
type XmlLogFile struct {
	file    string        // 日志文件
	XMLName xml.Name      `xml:"logs"`         // 入口节点
	Version int16         `xml:"version,attr"` // 版本号
	Name    string        `xml:"name,attr"`    // 名字
	Stus    XmlLogStructs `xml:"struct"`       // 日志结构数组
	Logs    XmlLogStructs `xml:"log"`          // 日志数组
}
