package webRoute

import "github.com/gin-gonic/gin"

type Router struct{}

func (Router) Register(engine *gin.Engine) {
	HomeRouter{}.Load(engine)   // 欢迎页
	WsTestRouter{}.Load(engine) // web-socket-test
}
