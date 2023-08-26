package tools

import (
	"math"
	"strconv"
	"sync"

	uuid "github.com/satori/go.uuid"

	"rnv-mmq/types"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type correct struct {
	msg          string
	messageId    string
	traceId      string
	context      *gin.Context
	businessType string
}

var ins *correct
var once sync.Once

// NewCorrectWithGinContext 正确返回值
func NewCorrectWithGinContext(msg string, ctx *gin.Context) *correct {
	once.Do(func() { ins = &correct{msg: ""} })
	ins.msg = msg
	ins.context = ctx
	ins.messageId = uuid.NewV4().String()
	return ins
}

// NewCorrectWithBusiness 返回正确值（不使用gin context）
func NewCorrectWithBusiness(msg, businessType, traceId string) *correct {
	once.Do(func() { ins = &correct{msg: ""} })
	ins.msg = msg
	ins.businessType = businessType
	ins.traceId = traceId
	ins.messageId = uuid.NewV4().String()
	return ins
}

// Ok 操作成功
func (receiver correct) Ok() types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "OK"
	}

	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      0,
		Status:       200,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}

// Datum 读取成功
func (receiver correct) Datum(content any) types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "OK"
	}

	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      content,
		Status:       200,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}

// DataForPager 返回分页数据
func (receiver correct) DataForPager(db *gorm.DB, read func(*gorm.DB) types.MapStringToAny) types.StdResponse {
	var count int64

	if receiver.msg == "" {
		receiver.msg = "OK"
	}

	pageStr := receiver.context.Query("__page__")
	limitStr := receiver.context.DefaultQuery("__limit__", "1000")
	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)

	if pageStr != "" {
		db.Count(&count)

		return types.StdResponse{
			MessageId: receiver.messageId,
			TraceId:   receiver.traceId,
			Msg:       receiver.msg,
			Content:   read(db.Offset((page - 1) * limit).Limit(limit)),
			Status:    200,
			ErrorCode: 0,
			Pagination: map[string]interface{}{
				"page_curr": page,
				"page_last": math.Ceil(float64(count/int64(limit))) + 1,
				"page_prev": page - 1,
				"page_next": page + 1,
				"page_size": limit,
				"count":     count,
			},
			BusinessType: receiver.businessType,
		}
	} else {
		if limit <= 0 {
			return types.StdResponse{
				MessageId:    receiver.messageId,
				Msg:          receiver.msg,
				Content:      read(db),
				Status:       200,
				ErrorCode:    0,
				Pagination:   nil,
				BusinessType: receiver.businessType,
			}
		} else {
			return types.StdResponse{
				MessageId:    receiver.messageId,
				Msg:          receiver.msg,
				Content:      read(db.Limit(limit)),
				Status:       200,
				ErrorCode:    0,
				Pagination:   nil,
				BusinessType: receiver.businessType,
			}
		}
	}
}

// DataForJqueryDataTable 返回jquery-dataTable格式分页数据
func (receiver correct) DataForJqueryDataTable(db *gorm.DB, read func(*gorm.DB) types.MapStringToAny) types.JqueryDataTableResponse {
	var count int64

	if receiver.msg == "" {
		receiver.msg = "OK"
	}

	startStr := receiver.context.DefaultQuery("start", "1")
	limitStr := receiver.context.DefaultQuery("length", "1000")

	limit, _ := strconv.Atoi(limitStr)
	start, _ := strconv.Atoi(startStr)
	db.Count(&count)

	return types.JqueryDataTableResponse{
		Content:         read(db.Offset(start).Limit(limit)),
		Draw:            receiver.context.DefaultQuery("draw", "1"),
		Start:           start,
		Length:          limit,
		Max:             &count,
		RecordsTotal:    &count,
		RecordsFiltered: &count,
	}
}

// Created 新建成功
func (receiver correct) Created(content any) types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "新建成功"
	}

	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      content,
		Status:       201,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}

// Updated 更新成功
func (receiver correct) Updated(content any) types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "编辑成功"
	}

	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      content,
		Status:       202,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}

// Deleted 删除成功
func (receiver correct) Deleted() types.StdResponse {
	if receiver.msg == "" {
		receiver.msg = "删除成功"
	}
	return types.StdResponse{
		MessageId:    receiver.messageId,
		TraceId:      receiver.traceId,
		Msg:          receiver.msg,
		Content:      nil,
		Status:       204,
		ErrorCode:    0,
		Pagination:   nil,
		BusinessType: receiver.businessType,
	}
}
