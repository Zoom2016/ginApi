package service

import (
	"encoding/json"
	"errors"
	"ginApi/common/enum"
	"ginApi/common/response"
	"ginApi/common/tools"
	"ginApi/models"
)

type UserService struct {
}

type UserCopy struct {
	models.User
}

type customerUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	UserName string `json:"username"`
}

type IdParam struct {
	Id int `form:"id" json:"id" binding:"required" msg:"id必填"`
}

type UserNameParam struct {
	UserName string `form:"username" json:"username" binding:"required,min=1" msg:"用户名不能为空"`
}

type PasswordParam struct {
	Password string `form:"password" json:"password" binding:"required,min=1" msg:"密码不能为空"`
}

// AddParam 添加用户参数
type AddParam struct {
	Age  int    `form:"age" json:"age"`
	Name string `form:"name" json:"name" binding:"required"`
	Page int    `form:"page" json:"page"`
	UserNameParam
	PasswordParam
	Phone string `form:"phone" json:"phone" binding:"required" msg:"手机不能为空"`
}

// ListParam 用户列表查询参数
type ListParam struct {
	Name string `form:"name" json:"name"`
	Age  int    `form:"age" json:"age"`
	Page int    `form:"page" json:"page"`
}

// EditParam 修改用户的参数
type EditParam struct {
	AddParam
	IdParam
}

// DelParam 删除用户的参数
type DelParam struct {
	IdParam
}

// LoginParam 登录所需参数
type LoginParam struct {
	UserNameParam
	PasswordParam
}

// MarshalJSON 如果一个类型实现了 Marshaler 接口，即具有 MarshalJSON() ([]byte, error) 方法，
// 那么编码器在遇到该类型时会调用 MarshalJSON 方法来自定义该类型的 JSON 编码行为
func (u *UserCopy) MarshalJSON() ([]byte, error) {
	// lists 返回的字段
	user := customerUser{
		Id:       u.Id,
		Name:     u.Name,
		Age:      u.Age,
		UserName: u.UserName,
	}
	return json.Marshal(user)
}

func (this UserService) Lists(userParam *ListParam) ([]*UserCopy, error) {
	var users []*models.User
	err := models.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	list := make([]*UserCopy, 0)
	for _, user := range users {
		list = append(list, &UserCopy{*user})
	}
	return list, nil
}

func (this UserService) Add(param *AddParam) int {
	var user models.User // User model
	var userCopy models.User
	var phone models.Phone // Phone model
	user.Name = param.Name
	user.Age = param.Age
	user.UserName = param.UserName
	user.Password = tools.Sha1(param.Password)

	phone.Phone = param.Phone

	//事务提交
	tx := models.DB.Begin()
	// 先查一下用户是否已经存在
	tx.Model(&userCopy).Where("username=?", user.UserName).First(&userCopy)
	if userCopy.Id != 0 {
		tx.Rollback()
		panic(&response.Response{
			Code: enum.CodeParamError,
			Msg:  "用户名已经注册",
		})
	}
	tx.Create(&user) // 创建用户记录
	phone.UserId = user.Id
	tx.Create(&phone) // 拿到用户id, 创建Phone记录
	tx.Commit()
	return user.Id
}

func (this UserService) Edit(param *EditParam) {
	var user models.User
	var phone models.Phone
	user.Name = param.Name
	user.UserName = param.UserName
	user.Password = tools.Sha1(param.Password)
	user.Age = param.Age
	phone.Phone = param.Phone
	// 开启事务
	tx := models.DB.Begin()
	tx.Model(&user).Where("id=?", param.Id).Updates(&user)
	tx.Model(&phone).Where("userId=?", param.Id).Updates(&phone)
	tx.Commit()
}

func (this UserService) Del(param *DelParam) {
	tx := models.DB.Begin()
	tx.Delete(&models.User{}, param.Id)
	tx.Where("userId=?", param.Id).Delete(&models.Phone{})
	tx.Commit()
}

func (this UserService) Login(param *LoginParam) (models.User, error) {
	var user models.User
	// 根据用户名查询用户记录，对比密码是否正确，对比前需要将传入密码用Sha1转一下
	models.DB.Model(&user).Where("username=?", param.UserName).First(&user)
	if tools.Sha1(param.Password) != user.Password {
		panic(&response.Response{
			Code: enum.CodeParamError,
			Msg:  "密码错误",
		})
		return user, errors.New("密码错误")
	}
	return user, nil
}
