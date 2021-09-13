module github.com/poniteru/go-coin-watcher

go 1.15

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.10.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/google/uuid v1.2.0
	github.com/huobirdcenter/huobi_golang v0.0.0
	github.com/jmoiron/sqlx v1.3.3
	github.com/kr/text v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/viper v1.8.1
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/unrolled/secure v1.0.7
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/huobirdcenter/huobi_golang => ./huobi_golang
