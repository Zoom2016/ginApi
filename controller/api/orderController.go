package api

import (
	"ginApi/common/response"
	"ginApi/common/tools"
	"ginApi/controller"
	"ginApi/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type OrderController struct {
	controller.BaseController
}

func (this OrderController) Lists(c *gin.Context) {
	var param service.OrderParam
	// 根据 body为空会报错EOF err.Error() != "EOF" ,查询条件是否必填为判断
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil && err.Error() != "EOF" {
		tools.GetError(err, param)
		return
	}
	// 这里的userId是middleware中设置的当前用户Id, 例:c.Set("userId", xxx)
	param.UserId = c.GetInt("userId")
	data, lastPage, total, _ := service.OrderService{}.Lists(&param)
	response.Success(c, &response.Response{Data: map[string]interface{}{
		"lists":    data,
		"lastPage": lastPage,
		"total":    total,
	}})
}
