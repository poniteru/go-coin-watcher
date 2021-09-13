package model

type (
	/**
	定义信息结构体
	该 struct 可以绑定在 Form 和 JSON 中
	binding:"required" 意思是必要参数。如果未提供，Bind 会返回 error
	uri:"" xml:""
	*/
	Msg struct {
		MessageBody   string `form:"messageBody" json:"messageBody"`
		SenderAddress string `form:"senderAddress" json:"senderAddress"`
		ReceivedTime  string `form:"receivedTime" json:"receivedTime"`
		Content       string `form:"content" json:"content" binding:"required"`
		Token         string `form:"token" json:"token"`
	}
)
