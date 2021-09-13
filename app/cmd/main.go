package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/poniteru/go-coin-watcher/app/config"
	"github.com/poniteru/go-coin-watcher/app/router"
	"github.com/poniteru/go-coin-watcher/digitcoin"
	"github.com/unrolled/secure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	workingDirectory, _ := os.Getwd()
	fmt.Println("Working directory: ", workingDirectory)
}

func main() {
	var err error

	engine := gin.Default()
	if config.EnableSSL {
		engine.Use(TlsHandler())
	}

	router.InitRouter(engine)

	go digitcoin.RunAll()

	go func() {
		// 绑定端口，然后启动应用
		if config.EnableSSL {
			err = engine.RunTLS(config.ServeAddr, config.CertFile, config.KeyFile)
		} else {
			err = engine.Run(config.ServeAddr)
		}
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("run failed: ", err)
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//if err := engine.Shutdown(ctx); err != nil {
	//	log.Fatal("Server forced to shutdown: ", err)
	//}
	// todo 增加退出和取消订阅逻辑
	log.Println("Server exiting")
}

func TlsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			//SSLHost:     "localhost:88",
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}
