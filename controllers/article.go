package controllers

import (
	"bxg01/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"log"
	"math"
	"path"
	"time"
)

type ArticleController struct {
	beego.Controller
}

//展示文章列表页
func (this *ArticleController) ShowArticleList() {
	//session判断
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	//获取数据
	//高级查询
	//---指定表
	o := orm.NewOrm()
	qs := o.QueryTable("Article") //queryseter查询集合
	var articles []models.Article
	/*_,err := qs.All(&articles)
	if err !=nil{
		beego.Info("查询数据错误")
	}*/
	//根据选择的文章类型查询相应的文章
	typeName := this.GetString("select")
	var count int64

	//获取总页数
	pageSize := 2

	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	start := (pageIndex - 1) * pageSize
	if typeName == "" {
		//查询总记录数
		count, _ = qs.Count()
	} else {
		count, _ = qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()
	}
	beego.Info(count)
	pageCount := int(math.Ceil(float64(count) / float64(pageSize)))
	//获取文章类型
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types

	if typeName == "" {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)
		//beego.Info(articles)
	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles) //RelatedSel("表名")
	}

	//传递数据
	this.Data["typeName"] = typeName
	this.Data["userName"] = userName
	this.Data["count"] = count
	this.Data["pageCount"] = pageCount
	this.Data["pageIndex"] = pageIndex
	this.Data["articles"] = articles

	//指定视图布局
	this.Layout = "layout.html"
	this.TplName = "index.html"
}

//展示添加文章页面
func (this *ArticleController) ShowAddArticle() {
	//查询所有类型数据，并展示
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	this.Data["types"] = types
	this.TplName = "add.html"
}

//获取添加文章数据
func (this *ArticleController) HandleAddArticle() {
	//1.获取数据
	articleName := this.GetString("articleName")
	content := this.GetString("content")
	//2.校验数据
	if articleName == "" || content == "" {
		this.Data["errmsg"] = "添加数据不完整"
		this.TplName = "add.html"
		return
	}
	beego.Info(articleName, content)
	//处理文件
	file, head, err := this.GetFile("uploadname")
	defer file.Close()
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}
	//this.SaveToFile("uploadname","./static/img/"+head.Filename)
	// 1.文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "文件太大"
		this.TplName = "add.html"
		return
	}
	// 2.文件格式
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		this.Data["errmsg"] = "文件格式错误"
		this.TplName = "add.html"
		return
	}
	// 3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	// 存储
	//this.SaveToFile("uploadname","./static/img/"+head.Filename)
	this.SaveToFile("uploadname", "./static/img/"+fileName)
	//3.处理数据
	// 插入操作
	o := orm.NewOrm()
	var article models.Article
	article.ArtiName = articleName
	article.Acontent = content
	article.Aimg = "/static/img/" + fileName
	//给文章添加类型
	typeName := this.GetString("select")
	//根据
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Read(&articleType, "TypeName")
	article.ArticleType = &articleType

	o.Insert(&article)
	//4.返回页面
	this.Redirect("/showArticleList", 302)
}

//显示文章详情页
func (this *ArticleController) ShowArticleDetail() {
	//获取数据
	id, err := this.GetInt("articleId")
	//数据校验
	if err != nil {
		beego.Info("传递的链接错误")
		return
	}
	//操作数据-查询
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	//o.Read(&article)
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id", id).One(&article)

	//修改阅读量
	article.Acount += 1
	o.Update(&article)
	//多对多插入浏览记录
	m2m := o.QueryM2M(&article, "Users") // 第一个参数的对象，主键必须有值，第二个参数为对象需要操作的 M2M 字段
	userName := this.GetSession("userName")
	if userName == nil {
		this.Redirect("/login", 302)
		return
	}
	var user models.User
	user.Name = userName.(string)
	o.Read(&user, "Name")
	//插入操作
	m2m.Add(user)
	//查询
	//o.LoadRelated(&article,"Users")//不能去重复
	var users []models.User
	//Filter(结构体对象的字段__结构体名称__字段可以是主键字段)
	//过滤出看过这篇文章的所有人
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)
	this.Data["article"] = article
	this.Data["users"] = users
	//返回数据给视图
	userLayout := this.GetSession("userName")
	this.Data["userName"] = userLayout.(string)
	this.Layout = "layout.html"
	this.TplName = "content.html"
}

//显示编辑文章页面
func (this *ArticleController) ShowUpdateArticle() {
	//获取数据
	id, err := this.GetInt("articleId")
	//校验数据
	if err != nil {
		beego.Info("请求文章错误")
		return
	}
	//数据处理
	//--查询相应的文章
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Read(&article)

	//返回视图
	this.Data["article"] = article
	this.TplName = "update.html"
}

//封装上传文件函数
func UploadFile(this *beego.Controller, filePath string) string {
	log.Println("call UploadFile。。。。", filePath)
	file, head, err := this.GetFile(filePath)
	if err != nil {
		beego.Info(err)
		this.Data["errmsg"] = err
		return "NoImg"
	}
	if head.Filename == "" {
		return "NoImg"
	}
	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return ""
	}
	defer file.Close()
	// 1.文件大小
	if head.Size > 5000000 {
		this.Data["errmsg"] = "文件太大"
		this.TplName = "add.html"
		return ""
	}
	// 2.文件格式
	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		this.Data["errmsg"] = "文件格式错误"
		this.TplName = "add.html"
		return ""
	}
	// 3.防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	// 存储
	this.SaveToFile(filePath, "./static/img/"+fileName)
	return "/static/img/" + fileName
}

//处理编辑界面数据
func (this *ArticleController) HandleUpdateArticle() {
	//1获取数据
	log.Println("call HandleUpdateArticle...")
	id, err := this.GetInt("articleId")
	articleName := this.GetString("articleName")
	content := this.GetString("content")

	filePath := UploadFile(&this.Controller, "uploadname")
	//2数据校验
	if err != nil || articleName == "" || content == "" || filePath == "" {
		beego.Info("请求路径错误")
		return
	}
	//3数据处理
	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	err = o.Read(&article)
	if err != nil {
		beego.Info("更新的文章不存在")
		return
	}
	article.ArtiName = articleName
	article.Acontent = content
	if filePath != "NoImg" {
		article.Aimg = filePath
	}
	o.Update(&article)
	//4返回视图
	this.Redirect("/article/showArticleList", 302)
}

//删除文章处理
func (this *ArticleController) DeleteArticle() {
	//1获取数据
	id, err := this.GetInt("articleId")
	//2校验数据
	if err != nil {
		beego.Info("删除文件请求路径错误")
		return
	}
	//3数据处理
	o := orm.NewOrm()
	var artilce models.Article
	artilce.Id = id
	o.Delete(&artilce)
	//4返回视图
	this.Redirect("/showArticleList", 302)
}

//展示添加类型页面
func (this *ArticleController) ShowAddType() {
	//查询
	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)
	//传递数据
	this.Data["types"] = types
	this.TplName = "addType.html"
}

//处理添加类型数据
func (this *ArticleController) HandleAddType() {
	//1获取数据
	typeName := this.GetString("typeName")
	//2校验数据
	if typeName == "" {
		beego.Info("信息不完整，重新输入")
		return
	}
	//3处理数据
	//插入数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeName = typeName
	o.Insert(&articleType)
	//4返回数据
	//this.Redirect("/ShowAddType",302)
	this.Redirect("/article/addType", 302)
}

//删除文章类型
func (this *ArticleController) DeleteType(){
	//获取数据
	id,err := this.GetInt("id")
	//校验数据
	if err != nil{
		beego.Error("删除类型错误",err)
		return
	}
	//处理数据
	//----1。获取orm对象
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType)
	//返回视图
	this.Redirect("/article/addType",302)
}
