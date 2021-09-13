package router

import (
	"github.com/gin-gonic/gin"
	"github.com/poniteru/go-coin-watcher/app/api"
	"github.com/poniteru/go-coin-watcher/app/config"
	"github.com/poniteru/go-coin-watcher/common/constant"
	"github.com/poniteru/go-coin-watcher/telegram"
)

// 初始化路由
func InitRouter(engine *gin.Engine) {
	baseGroup := engine.Group(config.BasePath)
	v1 := baseGroup.Group(constant.V1)
	{
		tgBotGroup := v1.Group(config.TgWebHookPath)
		for _, botBean := range telegram.BotMap {
			tgBotGroup.POST(botBean.BotHookPath, api.TgWebhookListenerV2(botBean))
		}
	}
}
