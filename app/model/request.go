package model

type (
	/**
	定义请求结构体
	该 struct 可以绑定在 Form 和 JSON 中
	binding:"required" 意思是必要参数。如果未提供，Bind 会返回 error
	uri:"" xml:""
	*/
	Request struct {
		Sign string `form:"sign" json:"sign"`
		Data string `form:"data" json:"data"`
		Key  string `form:"key" json:"key"`
		//Message string `form:"message" json:"message" binding:"required"`
	}
)
