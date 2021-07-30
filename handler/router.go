package handler

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func InitRouter()  {
	// 创建路由
	router := fasthttprouter.New()

	router.GET("/userLogin", UserLogin)

	router.POST("/room", CreateRoom)

	router.PUT("/room/:roomid/enter", EnterRoom)

	router.PUT("/roomLeave", LeaveRoom)

	router.GET("/room/:roomid", GetRoom)

	router.GET("/room/:roomid/users", GetUserList)

	router.POST("/roomList", GetRoomList)

	router.POST("/user", CreateUser)

	router.GET("/user/:username", GetUser)

	router.POST("/message/send", SendMessage)

	router.POST("/message/retrieve", RetrieveMessage)

	fasthttp.ListenAndServe(":8080", router.Handler)
}