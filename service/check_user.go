package service

import (
	"github.com/gin-gonic/gin"
	"Miniprogram-server-Golang/model"
	"Miniprogram-server-Golang/serializer"
)

// CheckIsRegisteredService 管理用户注册服务
type CheckUserService struct {
	UserId string `form:"userid" json:"userid"`
	Corpid string `form:"corpid" json:"corpid"`
	Uid    string `form:"uid" json:"uid"`
	Token  string `form:"token" json:"token"`
}

// 用于检测用户标识是否已经被绑定
func (service *CheckUserService) CheckUser(c *gin.Context) serializer.Response {
	if !model.CheckToken(service.Uid, service.Token) {
		return serializer.ParamErr("token验证错误", nil)
	}

	//	根据corpid找到公司名称
	var corp model.Corp
	if err := model.DB.Where(&model.Corp{Corpid: service.Corpid}).First(&corp); err != nil {
		return serializer.Err(10006, "获取企业信息失败", nil)
	}

	corpid := corp.ID
	//	根据corpid查找用户-企业绑定信息
	var corpBind model.WxMpBindInfo
	if err := model.DB.Where(&model.WxMpBindInfo{OrgId:corpid, Username:service.UserId, Isbind:1}).First(&corpBind); err != nil {
		//	错误码未知，张老师没有写到，有待修改
		return serializer.Err(100019, "用户和企业未绑定", nil)
	}

	wxuid := corpBind.WxUid

	if wxuid == service.Uid {
		//	正确的返回结果
		return serializer.BuildUserCheckResponse(0, service.Corpid, service.UserId)
	} else {
		//	这里不确定是返回错误信息还是显示用户已存在。接口文档和php代码不一致，目前以php代码为准
		return serializer.Err(100020, "该用户已被其他微信绑定，每个用户只能被一个微信绑定", nil)
		//return serializer.BuildUserCheckResponse(1, service.Corpid, service.UserId)
	}
}
