// Copyright 2021 CloudJ Company Limited. All rights reserved.

package web

import (
	"cloudiac/common"
	"cloudiac/configs"
	_ "cloudiac/docs" // 千万不要忘了导入你上一步生成的docs
	"cloudiac/portal/consts"
	"cloudiac/portal/libs/ctrl"
	"cloudiac/portal/libs/ctx"
	api_v1 "cloudiac/portal/web/api/v1"
	"cloudiac/portal/web/middleware"
	"cloudiac/utils"
	"cloudiac/utils/logs"
	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"io"
)

var logger = logs.Get()

func GetRouter() *gin.Engine {
	w := ctrl.WrapHandler

	e := gin.New()
	e.Use(gin.RecoveryWithWriter(io.MultiWriter(
		gin.DefaultWriter,
		logs.MustGetLogWriter("error"),
	)))

	// 允许跨域
	e.Use(w(middleware.Cors))
	e.Use(w(middleware.Operation))
	e.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	e.GET("/system/info", w(func(c *ctx.GinRequest) {
		c.JSONSuccess(gin.H{
			"version": common.VERSION,
			"build":   common.BUILD,
		})
	}))
	api_v1.Register(e.Group("/api/v1"))

	// 直接提供静态文件访问，生产环境部署时也可以使用 nginx 反代
	e.StaticFS(consts.ReposUrlPrefix, gin.Dir(consts.LocalGitReposPath, true))
	return e
}

func StartServer() {
	conf := configs.Get()
	utils.SetGinMode()
	e := GetRouter()
	logger.Infof("starting server on %v", conf.Listen)
	// API 接口总是使用 http 协议，ssl 证书由 nginx 管理
	if err := e.Run(conf.Listen); err != nil {
		logger.Fatalln(err)
	}
}
