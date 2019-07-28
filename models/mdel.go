package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	//"github.com/astaxie/beego"
	//"database/sql"
)

//定义一个结构体对象
type User struct {
	Id       int
	Name     string
	PassWord string
	Articles []*Article `orm:"reverse(many)"`
}
type Article struct {
	Id          int          `orm:"pk;auto"`
	ArtiName    string       `orm:"size(20)"`
	Atime       time.Time    `orm:"auto_now"`
	Acount      int          `orm:"default(0)"`
	Acontent    string       `orm:"size(500)"`
	Aimg        string       `orm:"size(100)"`
	ArticleType *ArticleType `orm:"rel(fk);null;on_delete(set_null);null"`
	Users       []*User      `orm:"rel(m2m)"` //rel正向设置
}

//类型表
type ArticleType struct {
	Id       int
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	//操作数据库的代码
	//第一个参数是数据库驱动
	//"root:123456@tcp(127.0.0.1:3306)/数据库名称?charset=utf8"
	/*conn,err:=sql.Open("mysql","root:Xiaoyi522@tcp(115.29.176.53:3306)/test_for_g?charset=utf8")
	if err!=nil{
		beego.Info("连接错误",err)
		beego.Error("连接错误",err)
		return
	}
	//关闭数据库
	defer conn.Close()*/
	//创建表
	/*_,err=conn.Exec("create table customer(nam varchar(32),paword varchar(20))")
	if err!=nil{
		beego.Error("创建表失败",err)
	}*/

	//conn.Exec("insert into customer(name,paword) values (?,?)", "jack","123456")

	//查询
	/*res,err := conn.Query("select name from customer")
	var name string
	for res.Next(){
		res.Scan(&name)
		beego.Info(name)

	}*/

	//------------------------------ORM操作数据库
	//1-获取连接对象
	//在orm中双下划线__有特殊含义 pass_word不能这样定义结构体字段
	orm.RegisterDataBase("default", "mysql", "root:Xiaoyi522@tcp(115.29.176.53:3306)/test_for_g?charset=utf8")

	//2-创建表
	orm.RegisterModel(new(User), new(Article), new(ArticleType))

	//生产表
	//第一个参数是数据库别名,第二个参数是是否强制更新,第三个参数是否可见生成过程
	orm.RunSyncdb("default", false, true)

}
