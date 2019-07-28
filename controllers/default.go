package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"bxg01/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["data"] = "China"
	c.TplName = "test.html"
}
func (c *MainController) Post() {
	c.Data["data"] ="上海一期"
	c.TplName = "test.html"
}
func (c *MainController) ShowGet() {
	//获取ORM对象
	o:=orm.NewOrm()
	beego.Info(o)
	//执行某个函数操作，增删改查
	//1-添加
	/*var user models.User
	user.Name = "heima"
	user.PassWord = "hehe"

	count,err:=o.Insert(&user)*/
	//2-查询操作
	/*var user models.User
	user.Id=1
	err:=o.Read(&user,"Id")
	if err!=nil{
		beego.Error("查询失败")
	}*/

	//3-更新操作
	/*var user models.User
	user.Id=1
	err:=o.Read(&user)
	if err!=nil{
		beego.Error("更新的数据不存在")
	}
	user.Name ="Tom"
	count,err:= o.Update(&user)
	if err!=nil{
		beego.Error("更新失败")
	}
	beego.Info(count)*/
	//返回结果
	//beego.Info(user)
	/*if err!=nil{
		beego.Error("insert fail")
	}
	beego.Info(count)*/

	//4-删除操作
	var user models.User
	user.Id=1
	err:=o.Read(&user)
	if err!=nil{
		beego.Error("删除的数据不存在")
	}
	count,err :=o.Delete(&user)
	if err!=nil{
		beego.Error("删除失败")
	}
	beego.Info(count)
	c.Data["data"] ="上海一期"
	c.TplName = "test.html"
}