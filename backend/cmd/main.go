package main

import (
	"fmt"
	"log"
	"studyexam/config"
	"studyexam/database"
	"studyexam/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("正在启动学习平台后端服务...")

	if err := config.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	fmt.Println("配置加载成功")

	if err := database.InitDB(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	fmt.Println("数据库连接成功")

	gin.SetMode(config.AppConfig.Server.Mode)

	r := routes.SetupRouter()

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	fmt.Printf("服务启动成功，监听地址: http://localhost%s\n", addr)

	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
