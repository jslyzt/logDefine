package logger

import (
	"encoding/json"

	"github.com/jslyzt/logDefine"
)

// --------------------------------------------------------------------
// logger ResultImage结构定义
// 识别图片信息
type Logger_ResultImage struct { // version 1
	ImageID  string  `json:"imageID"`  // desc: 识别服务器返回的图片target_id
	ServerIP string  `json:"serverIP"` // desc: 识别服务器IP
	ImageURL string  `json:"imageURL"` // desc: 识别服务器返回的图片地址
	Score    float32 `json:"score"`    // desc: 图片得分

}

// logger ResultImage序列化方法
// ToString
func (node Logger_ResultImage) ToString() string {
	return logDefine.ToString(node.ImageID, node.ServerIP, node.ImageURL, node.Score)
}

// ToJson
func (node Logger_ResultImage) ToJson() string {
	data, err := json.Marshal(node)
	if err == nil {
		return string(data)
	}
	return ""
}

// logger ResultImage反序列化方法
// FromString
func (node Logger_ResultImage) FromString(data []byte, index int) int {
	return logDefine.FromString(data, index, &node)
}

// FromJson
func (node Logger_ResultImage) FromJson(data []byte) {
	json.Unmarshal(data, node)
}

// --------------------------------------------------------------------
// logger LogData结构定义
// 发送到存储系统的日志结构体
type Logger_LogData struct { // version 1
	RequestID              string                 `json:"requestID"`              // desc: 请求的ID
	Token                  string                 `json:"token"`                  // desc: 鉴权令牌
	Latitude               float64                `json:"latitude"`               // desc: 纬度
	Longitude              float64                `json:"longitude"`              // desc: 经度
	Collection             string                 `json:"collection"`             // desc: 识别图片，多图集使用逗号隔开
	Number                 int                    `json:"number"`                 // desc: 返回结果的top
	ClientIP               string                 `json:"clientIP"`               // desc: 客户端IP
	Image                  string                 `json:"image"`                  // desc: 用户请求图片
	ResultImage            []*Logger_ResultImage  `json:"resultImage"`            // desc: 识别服务器返回的图片
	Result                 map[string]interface{} `json:"result"`                 // desc: 返回客户的端的结果json
	CreateTime             string                 `json:"createTime"`             // desc: 请求的时间
	Timeconst              float64                `json:"timeconst"`              // desc: 请求总耗时
	AppKey                 string                 `json:"appKey"`                 // desc: 应用ID
	Appname                string                 `json:"appname"`                // desc: 应用名称
	Useragent              string                 `json:"useragent"`              // desc: 用户代理
	Version                string                 `json:"version"`                // desc: 版本号
	RecognizeTimeConsuming float64                `json:"recognizeTimeConsuming"` // desc: getFeature时间

}

// logger LogData序列化方法
// ToString
func (node Logger_LogData) ToString() string {
	return logDefine.ToString(node.RequestID, node.Token, node.Latitude, node.Longitude, node.Collection, node.Number, node.ClientIP, node.Image, node.ResultImage, node.Result, node.CreateTime, node.Timeconst, node.AppKey, node.Appname, node.Useragent, node.Version, node.RecognizeTimeConsuming)
}

// ToJson
func (node Logger_LogData) ToJson() string {
	data, err := json.Marshal(node)
	if err == nil {
		return string(data)
	}
	return ""
}

// logger LogData反序列化方法
// FromString
func (node Logger_LogData) FromString(data []byte, index int) int {
	return logDefine.FromString(data, index, &node)
}

// FromJson
func (node Logger_LogData) FromJson(data []byte) {
	json.Unmarshal(data, node)
}

// --------------------------------------------------------------------
// logger sdkReco结构定义
// sdk reco记录日志
type Logger_sdkReco struct { // version 1
	Business  Logger_LogData         `json:"business"`  // desc: 日志信息
	OauthInfo map[string]interface{} `json:"oauthInfo"` // desc: 鉴别信息

}

// logger sdkReco序列化方法
// ToString
func (node Logger_sdkReco) ToString() string {
	return node.ToAString(node.GetAppend())
}
func (node Logger_sdkReco) ToAString(arr []interface{}) string {
	strlog := logDefine.ToString(node.Business, node.OauthInfo)
	if len(arr) > 0 {
		strlog = logDefine.ToString(arr...) + strlog
	}
	return strlog
}
func (node Logger_sdkReco) GetAlias() string {
	return "sdk-reco"
}
func (node Logger_sdkReco) GetAppend() []interface{} {
	return []interface{}{
		node.GetAlias(),
		logDefine.GetTime(nil),
	}
}

// ToJson
func (node Logger_sdkReco) ToJson() string {
	data, err := json.Marshal(node)
	if err == nil {
		return string(data)
	}
	return ""
}

// logger sdkReco反序列化方法
// FromString
func (node Logger_sdkReco) FromString(data []byte, index int) int {
	var alias, stime string
	return node.FromAString(data, index, &alias, &stime)
}

func (node Logger_sdkReco) FromAString(data []byte, index int, nodes ...interface{}) int {
	slen := logDefine.FromString(data, index, nodes...)
	if slen < len(data) {
		return logDefine.FromString(data, slen, &node)
	}
	return slen
}

// FromJson
func (node Logger_sdkReco) FromJson(data []byte) {
	json.Unmarshal(data, node)
}

// --------------------------------------------------------------------
// logger cloudReco结构定义
// cloud reco记录日志
type Logger_cloudReco struct { // version 1
	Business Logger_LogData `json:"business"` // desc: 日志信息

}

// logger cloudReco序列化方法
// ToString
func (node Logger_cloudReco) ToString() string {
	return node.ToAString(node.GetAppend())
}
func (node Logger_cloudReco) ToAString(arr []interface{}) string {
	strlog := logDefine.ToString(node.Business)
	if len(arr) > 0 {
		strlog = logDefine.ToString(arr...) + strlog
	}
	return strlog
}
func (node Logger_cloudReco) GetAlias() string {
	return "cloud-reco"
}
func (node Logger_cloudReco) GetAppend() []interface{} {
	return []interface{}{
		node.GetAlias(),
		logDefine.GetTime(nil),
	}
}

// ToJson
func (node Logger_cloudReco) ToJson() string {
	data, err := json.Marshal(node)
	if err == nil {
		return string(data)
	}
	return ""
}

// logger cloudReco反序列化方法
// FromString
func (node Logger_cloudReco) FromString(data []byte, index int) int {
	var alias, stime string
	return node.FromAString(data, index, &alias, &stime)
}

func (node Logger_cloudReco) FromAString(data []byte, index int, nodes ...interface{}) int {
	slen := logDefine.FromString(data, index, nodes...)
	if slen < len(data) {
		return logDefine.FromString(data, slen, &node)
	}
	return slen
}

// FromJson
func (node Logger_cloudReco) FromJson(data []byte) {
	json.Unmarshal(data, node)
}