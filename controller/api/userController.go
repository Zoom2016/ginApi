package api

import (
	"ginApi/common/config"
	"ginApi/common/enum"
	"ginApi/common/jwt"
	"ginApi/common/response"
	"ginApi/common/tools"
	"ginApi/controller"
	"ginApi/models"
	"ginApi/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type UserController struct {
	controller.BaseController
}

func (this UserController) Lists(c *gin.Context) {
	var ListParam service.ListParam
	// 根据 body为空会报错EOF err.Error() != "EOF" ,查询条件是否必填为判断
	if err := c.ShouldBindBodyWith(&ListParam, binding.JSON); err != nil && err.Error() != "EOF" {
		tools.GetError(err, ListParam)
		return
	}
	data, _ := service.UserService{}.Lists(&ListParam)
	response.Success(c, &response.Response{Data: data})
}

// Add 新增用户
func (this UserController) Add(c *gin.Context) {
	var param service.AddParam
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		tools.GetError(err, param)
		return
	}
	id := service.UserService{}.Add(&param)
	response.Success(c, &response.Response{
		Data: map[string]int{"id": id},
	})
}

func (this UserController) Edit(c *gin.Context) {
	var param service.EditParam
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		tools.GetError(err, param)
		return
	}
	service.UserService{}.Edit(&param)

	response.Success(c, &response.Response{
		Data: map[string]string{},
	})
}

func (this UserController) Del(c *gin.Context) {
	var delParam service.DelParam
	err := c.ShouldBindBodyWith(&delParam, binding.JSON)
	if err != nil {
		panic(err)
	}
	service.UserService{}.Del(&delParam)

	response.Success(c, &response.Response{
		Data: map[string]string{},
	})
}

// Login 用户登陆
func (this UserController) Login(c *gin.Context) {
	var loginParam service.LoginParam
	err := c.ShouldBindBodyWith(&loginParam, binding.JSON)
	if err != nil {
		tools.GetError(err, loginParam)
	}
	data, _ := service.UserService{}.Login(&loginParam)

	// 根据配置文件中的类型(token.type)，生成token
	var token string
	if config.Viper.GetString("token.type") == "jwt" {
		// 生成jwt token
		result, err := jwt.Jwt{}.CreateToken(data.Id)
		if err != nil {
			response.Fail(c, &response.Response{Code: enum.CodeParamError, Msg: enum.ErrMsg[enum.CodeSystemError]})
			return
		}
		token = result
	} else if config.Viper.GetString("token.type") == "token" {
		// 生成指定长度(token.length)的随机token
		var res map[string]string
		for {
			token = tools.RandString(config.Viper.GetInt("token.length"))
			res, _ = models.RedisDb.HGetAll(token).Result()
			if len(res) != 0 {
				token = tools.RandString(config.Viper.GetInt("token.length"))
			} else {
				break
			}
		}
		// 将token保存至Redis缓存
		models.RedisDb.HMSet("token:"+token, map[string]interface{}{
			"userId": data.Id,
			"token":  token,
		})
		// 设置token有效期(token.expire)
		models.RedisDb.Expire(
			"token:"+token,
			time.Duration(config.Viper.GetInt64("token.expire"))*time.Second,
		)
	} else {
		response.Fail(c, &response.Response{Code: enum.CodeSystemError, Msg: "token参数配置错误"})
	}

	response.Success(c, &response.Response{
		Data: map[string]string{"token": token},
	})
}
