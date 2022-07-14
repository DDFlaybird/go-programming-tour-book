package ginsample

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var dst string = "/home/coder/project/go_project/go-programming-tour-book/static/img/"

// redirect 路由组使用重定向方法
func redirect(ctx *gin.Context) {
	url := "https://www.baidu.com"
	ctx.Redirect(http.StatusMovedPermanently, url)
}

// redirect2 重定向指向的请求
func redirect2(ctx *gin.Context) {
	ctx.String(http.StatusOK, "重定向")
}

// third_party_content 获取三方数据
func third_party_content(ctx *gin.Context) {
	url := "https://pic3.zhimg.com/v2-58d652598269710fa67ec8d1c88d8f03_r.jpg?source=1940ef5c"
	r, err := http.Get(url)
	if err != nil || r.StatusCode != http.StatusOK {
		ctx.Status(http.StatusServiceUnavailable)
		return
	}
	body := r.Body
	contentLength := r.ContentLength
	contentType := r.Header.Get("Content-Type")
	ctx.DataFromReader(http.StatusOK, contentLength, contentType, body, nil)
}

// 没什么卵用就是简单展示如何使用中间件
// Middleware 简易的计时中间件
func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		begin := time.Now().Nanosecond()
		ctx.Next()
		end := time.Now().Nanosecond()
		totalCost := end - begin
		fmt.Printf("总耗时：%v ns \n", totalCost)
	}
}

func Run() {
	r := gin.Default()
	// middleware 中间件
	r.Use(Middleware())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"code": "2000", "msg": "sucess", "data": "Pong~@!"})
	})

	r.GET("/r", func(ctx *gin.Context) {
		url := "http://www.baidu.com"
		ctx.Redirect(http.StatusMovedPermanently, url)
	})

	v1 := r.Group("v1")
	{
		// 直接修改结构体中url来实现
		v1.GET("/redirect", func(ctx *gin.Context) {
			ctx.Request.URL.Path = "/v1/redirect_sub"
			fmt.Println(ctx.Request.URL.RawPath)
			r.HandleContext(ctx)
		})
		v1.GET("/redirect_sub", redirect2)

		// 使用Redirect方法实现。注意：需要带http不然默认凭借在路由下面
		v1.GET("/rd", redirect)
	}

	// 获取第三方返回数据
	r.GET("/read", third_party_content)

	// 文件上传
	r.POST("/upload", func(ctx *gin.Context) {
		fh, err := ctx.FormFile("filename")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": "4000", "msg": "请求失败", "data": nil})
		}
		dst := "/home/coder/project/go_project/go-programming-tour-book/"
		ctx.SaveUploadedFile(fh, dst+fh.Filename)
		ctx.JSON(http.StatusOK, gin.H{"code": "2000", "msg": fmt.Sprintf("文件: %s , 上传成功", fh.Filename), "data": nil})
	})

	// 多文件上传
	r.MaxMultipartMemory = 8 << 20
	r.POST("multUpload", func(ctx *gin.Context) {
		f, err := ctx.MultipartForm()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"code": "4000", "msg": "请求失败", "data": nil})

		}
		files := f.File["upload[]"]
		for _, fs := range files {
			ctx.SaveUploadedFile(fs, dst+fs.Filename)
		}
		ctx.JSON(http.StatusOK, gin.H{"code": "2000", "msg": fmt.Sprintf("%d 个文件上传成功", len(files)), "data": nil})

	})

	r.Run(":8082")
}
