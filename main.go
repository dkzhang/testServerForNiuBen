package main

import (
	"fmt"
	"github.com/kataras/iris"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

type RealTimeInfo struct {
	InspectionID int
	RecordID     int
	DateTime     string
	TextContent  string
	ImageUrl     string
}

type RealTimeRecords struct {
	Records []RealTimeInfo
}

var lock sync.Mutex

func main() {

	//generate data
	roundNumber := 1
	currentIndex := 0
	dataArray := make([]RealTimeInfo, 100)

	app := iris.New()
	app.Logger().SetLevel("debug")
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	//设置数据驱动的模板引擎为 std html/template （go 语言标准模板 参见 html/template 标准库）
	// 当然 iris 还有很多数据驱动的模板引擎 以后会介绍到
	//被解析与加载的所有模板文件放在  ./web/views 文件夹里面，并且以 .html 为文件扩张
	// Reload 方法设置为 true 表示开启开发者模式 将会每一次请求都重新加载 views 文件下的所有模板
	// RegisterView 注册加载 模板文件 与加载配置
	app.RegisterView(iris.HTML("./web/views", ".html").Reload(true))
	// 为特定HTTP错误注册自定义处理程序方法
	// 当出现 StatusInternalServerError 500错误，将执行第二参数回调方法
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		// ctx.Values() 是一个很有用的东西，主要用来使 处理方法与中间件 通信 记住真的很重要
		// ctx.Values().GetString("error") 获取自定义错误提示信息
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}
		ctx.Writef("(Unexpected) internal server error")
	})
	app.OnErrorCode(iris.StatusBadRequest, func(ctx iris.Context) {
		// ctx.Values() 是一个很有用的东西，主要用来使 处理方法与中间件 通信 记住真的很重要
		// ctx.Values().GetString("error") 获取自定义错误提示信息
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Bad Request error: %s", errMessage)
			return
		}
		ctx.Writef("(Unexpected) error")
	})
	// context.Handler 类型 每一个请求都会先执行此方法 app.Use(context.Handler)
	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})
	// context.Handler 类型 每一个请求最后执行 app.Done(context.Handler)
	app.Done(func(ctx iris.Context) {})

	// Method:   GET
	// Resource: http://localhost:8080
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<h1>Welcome</h1>")

	})

	app.Handle("GET", "/cambrian001/inspection-log/realtime", func(ctx iris.Context) {

		inspectionID, err := strconv.Atoi(ctx.URLParam("inspectionID"))
		if err != nil {
			ctx.Values().Set("error",
				fmt.Sprintf("InspectionID %s is illegal: %s",
					ctx.URLParam("inspectionID"), err.Error()))
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		recordID, err := strconv.Atoi(ctx.URLParam("recordID"))
		if err != nil {
			ctx.Values().Set("error",
				fmt.Sprintf("RecordID %s is illegal: %s",
					ctx.URLParam("recordID"), err.Error()))
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		lock.Lock()
		defer lock.Unlock()

		if inspectionID < roundNumber {
			//忽视请求轮次，返回当前轮次的所有
			ctx.JSON(RealTimeRecords{
				Records: dataArray[:currentIndex+1]})
		} else if inspectionID > roundNumber {
			ctx.Values().Set("error",
				fmt.Sprintf("The InspectionID requested %d is too large, current InspectionID is %d.",
					inspectionID, roundNumber))
			ctx.StatusCode(iris.StatusBadRequest)
			return
		} else {
			if recordID > currentIndex {
				ctx.Values().Set("error",
					fmt.Sprintf("The RecordID requested %d is too large, current RecordID is %d.",
						recordID, currentIndex))
				ctx.StatusCode(iris.StatusBadRequest)
				return
			} else {
				ctx.JSON(RealTimeRecords{
					Records: dataArray[recordID : currentIndex+1]})
			}

		}
	})

	app.Handle("GET", "/cambrian001/inspection-log-view/realtime", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	//初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	//generate data
	go func() {
		for {
			lock.Lock()
			currentIndex++
			if currentIndex >= 100 {
				roundNumber++
				currentIndex = 0
			}

			imageID := rand.Intn(20)
			imageUrl := ""
			if imageID <= 9 {
				imageUrl = fmt.Sprintf("http://inspection-robot.gribgp.com:9981/cambrian001/static/test%d.png", imageID)
			}

			dataArray[currentIndex] = RealTimeInfo{
				InspectionID: roundNumber,
				RecordID:     currentIndex,
				DateTime:     time.Now().Format("2006-01-02 15:04:05"),
				TextContent:  fmt.Sprint(roundNumber, " + ", currentIndex),
				ImageUrl:     imageUrl,
			}
			lock.Unlock()

			fmt.Println(currentIndex)
			time.Sleep(time.Second * 1)
		}
	}()

	//图片文件服务
	app.HandleDir("/cambrian001/static", "./assets")

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
