package controllers

import (
	"bxg01/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

//处理注册数据
func (this *UserController) HandlePost() {
	//1-获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//beego.Info(userName,pwd)
	//2-校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "注册数据不完整,重新注册"
		beego.Info("注册数据不完整，重新注册")
		this.TplName = "register.html"
		return
	}
	//3-操作数据
	//获取ORM对象
	o := orm.NewOrm()
	//获取插入对象
	var user models.User
	//给插入对象赋值
	user.Name = userName
	user.PassWord = pwd
	//插入数据
	o.Insert(&user)
	//返回结果

	//4-返回页面
	//this.Ctx.WriteString("注册成功")
	//跳转函数
	this.Redirect("/login", 302) //不能传递数据给视图
	//this.TplName = "login.html"//能传递数据给视图
}

func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	if userName == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}
	this.TplName = "login.html"
}
func (this *UserController) HandleLogin() {
	//1-获取数据
	userName := this.GetString("userName")
	pwd := this.GetString("password")
	//2-校验数据
	if userName == "" || pwd == "" {
		this.Data["errmsg"] = "登录数据不完整,重新登录"
		beego.Info("登录数据不完整，重新登录")
		this.TplName = "login.html"
		return
	}
	//3-操作数据
	//获取ORM对象
	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		this.Data["errmsg"] = "用户不存在"
		this.TplName = "login.html"
		return
	}
	if user.PassWord != pwd {
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}
	//4-返回页面
	//this.Ctx.WriteString("登录成功")
	data := this.GetString("remember")
	beego.Info(data)
	if data == "on" {
		this.Ctx.SetCookie("userName", userName, 100)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}
	this.SetSession("userName", userName)
	this.Redirect("/article/showArticleList", 302)
}

func (this *UserController) Logout() {
	//删除session
	this.DelSession("userName")
	//跳转到登录页
	this.Redirect("/login", 302)
}
