package logdefine

import "encoding/xml"

// 常量定义
const (
	TInt      int8 = iota // int
	TFloat                // float
	TDouble               // double
	TString               // string
	TDateTime             // 时间日期
	TBool                 // 布尔类型
	TShort                // 短整型
	TLong                 // 长整型
	TUserDef              // 自定义类型
)

// 常量定义
const (
	ETgo   int8 = iota // 导出go
	ETcpp              // 导出c++
	ETjava             // 导出java
)

// 常量定义
const (
	UDTnone  int8 = iota // 无
	UDTlist              // 列表
	UDTplist             // 指针列表
	UDTmap               // 键值对
	UDTpmap              // 指针键值对
)

// 变量定义
var (
	TimeFormateUnix = "2006-01-02T15:04:05+08:00"
)

// XMLLogNode 节点定义
type XMLLogNode struct {
	Xname string `xml:"name,attr"` // 节点名字
	Name  string // 真正名字
	SType string `xml:"type,attr"` // 节点类型 -- xml
	Type  int8   // 节点类型 -- true
	//Defvalue interface{} `xml:"defaultvalue,attr"` // 节点默认值
	Desc   string `xml:"desc,attr"` // 节点说明
	UDType int8   // 扩展类型
}

// XMLLogNodes 节点数组定义
type XMLLogNodes []XMLLogNode

// XMLLogStruct 日志描述定义
type XMLLogStruct struct {
	Name    string      `xml:"name,attr"`    // 名字
	UName   string      `xml:"uname,attr"`   // 名字
	Alias   string      `xml:"alias,attr"`   // 别名
	Version int16       `xml:"version,attr"` // 版本号
	Desc    string      `xml:"desc,attr"`    // 说明
	Nodes   XMLLogNodes `xml:"entry"`        // 节点列表
}

// XMLLogStructs 日志描述数组定义
type XMLLogStructs []XMLLogStruct

// XMLLogStrMap 日志描述map定义
type XMLLogStrMap map[string]*XMLLogStruct

// XMLLogFile 日志文件定义
type XMLLogFile struct {
	file    string        // 日志文件
	XMLName xml.Name      `xml:"logs"`         // 入口节点
	Version int16         `xml:"version,attr"` // 版本号
	Name    string        `xml:"name,attr"`    // 名字
	MName   string        // 大写名字
	Stus    XMLLogStructs `xml:"struct"` // 日志结构数组
	Logs    XMLLogStructs `xml:"log"`    // 日志数组
	StuMp   XMLLogStrMap  // 日志结构map
}
