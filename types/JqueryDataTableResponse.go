package types

import "encoding/json"

// JqueryDataTableResponse 标准响应结构
type JqueryDataTableResponse struct {
	Draw            string
	Start           int
	Length          int
	Content         map[string]any
	Max             *int64
	RecordsTotal    *int64
	RecordsFiltered *int64
}

// ToMap 转map[string]interface
func (receiver JqueryDataTableResponse) ToMap() map[string]any {
	return map[string]any{
		"content":         receiver.Content,
		"draw":            receiver.Draw,
		"start":           receiver.Start,
		"length":          receiver.Length,
		"max":             receiver.Max,
		"recordsTotal":    receiver.RecordsTotal,
		"recordsFiltered": receiver.RecordsFiltered,
	}
}

// ToJson 转json
func (receiver JqueryDataTableResponse) ToJson() []byte {
	data := receiver.ToMap()
	if marshal, err := json.Marshal(data); err != nil {
		return []byte{}
	} else {
		return marshal
	}
}

// ToJsonStr 转json字符串
func (receiver JqueryDataTableResponse) ToJsonStr() string {
	return string(receiver.ToJson())
}

// ToGinResponse 转gin响应格式 int+map[string]any
func (receiver JqueryDataTableResponse) ToGinResponse() (int, map[string]any) {
	return 200, receiver.ToMap()
}
