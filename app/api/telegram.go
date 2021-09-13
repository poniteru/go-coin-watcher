package api

import (
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/poniteru/go-coin-watcher/telegram"
	"log"
	"net/http"
	"strings"
)

/*func TgWebhookListener(ctx *gin.Context) {
	var update tgbotapi.Update
	err := ctx.ShouldBindJSON(&update)
	if err != nil {
		//返回错误信息
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("webhook received: %+v\n", update)
	//ctxCopy := ctx.Copy()
	//执行异步操作
	go handleMessage(&update)
	//log.Println("Done! in path" + ctxCopy.Request.URL.Path)
	//返回一个空对象 var retData struct{}
	ctx.JSON(http.StatusOK, nil)
}*/

func TgWebhookListenerV2(botBean *telegram.BotBean) func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		var update tgbotapi.Update
		err := ctx.ShouldBindJSON(&update)
		if err != nil {
			//返回错误信息
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("webhook received: %+v\n", update)
		//ctxCopy := ctx.Copy()
		//执行异步操作
		go handleMessage(&update, botBean)
		//log.Println("Done! in path" + ctxCopy.Request.URL.Path)
		//返回一个空对象 var retData struct{}
		ctx.JSON(http.StatusOK, nil)
	}
}

func handleMessage(update *tgbotapi.Update, botBean *telegram.BotBean) {
	//log.Printf("%+v\n", update)
	//log.Printf("%+v\n", update.Message)
	//log.Printf("%+v\n", update.Message.From.ID)
	//log.Printf("%+v\n", update.Message.Chat)
	//log.Printf("%+v\n", update.EditedMessage)
	if update.Message != nil {
		if !strings.HasPrefix(update.Message.Text, "/") {
			return
		}
		cmdStr := strings.Fields(update.Message.Text)[0]
		botBean.BotMsgHandler.GetMsgHandler(cmdStr)(update, botBean)
	}

}
