package main

import (
	"blog/controllers"
	"blog/route"
	"blog/models"
	log "github.com/alecthomas/log4go"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"runtime"
	"sync"
	"flag"
)

var(

	logfile = flag.String("log","./conf/blog-log.xml","log4go file path!")
	configfile = flag.String("config","./conf/blog.ini","config file path")

)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	log.LoadConfiguration(*logfile)
	defer log.Close()

	gin.SetMode(gin.DebugMode) //全局设置环境，此为开发环境，线上环境为gin.ReleaseMode

	err := models.InitBaseConfig(*configfile)
	if err != nil {
		log.Error("InitBaseConfig has error:%v", err)
		return
	}

	//前台文章浏览量入库
	go controllers.ProcessReadData()

	router := route.InitRouter()

	http.ListenAndServe(models.AppPort, router)
}
