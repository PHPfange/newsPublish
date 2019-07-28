package main
import (
	_ "bxg01/routers"
	"github.com/astaxie/beego"
	_ "bxg01/models"
)

func main() {
	beego.AddFuncMap("prepage",showPrePage)
	beego.AddFuncMap("nextpage",showNextPage)
	beego.Run()
}
//后台
func showPrePage(pageIndex int)int{
	if pageIndex == 1{
		return 1
	}
	return pageIndex-1
}

func showNextPage(pageIndex int,pageCount int)int{
	if pageIndex == pageCount{
		return pageCount
	}
	return pageIndex+1
}
