package handler

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/valyala/fasthttp"
	"net/http"
	"tianchi/models"
	"tianchi/services"
)

func CreateRoom(ctx *fasthttp.RequestCtx) {
	username := JWTAuth(ctx)
	if len(username) == 0 {
		l4g.Error("CreateRoom JWTAuth user err!")
		return
	}

	room := new(models.Room)
	args := ctx.Request.Body()
	err := json.Unmarshal(args, &room)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid input")
		l4g.Error("CreateRoom Unmarshal args err:%v", err)
		return
	}
	if room == nil || room.Name == "" {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid input")
		l4g.Error("CreateRoom room illegal:%v", room)
		return
	}

	id, err := services.CreateRoom(username, room)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid input")
		l4g.Error("%s CreateRoom err:%v", username, err)
		return
	}
	l4g.Debug("%s CreateRoom:%s", username, id)

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString(id)
}

func EnterRoom(ctx *fasthttp.RequestCtx) {
	username := JWTAuth(ctx)
	if len(username) == 0 {
		l4g.Error("EnterRoom JWTAuth user err!")
		return
	}

	if len(ctx.Request.URI().Path()) < len("/room//enter") {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid Room ID")
		return
	}

	roomId := string(ctx.Request.URI().Path())[len("/room/"):]
	roomId = roomId[:len(roomId)-len("/enter")]

	services.LeaveRoom(username)

	err := services.EnterRoom(username, roomId)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid Room ID")
		l4g.Error("%s EnterRoom room:%s err:%v", username, roomId, err)
		return
	}
	l4g.Debug("%s EnterRoom roomId:%v ok!", username, roomId)

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString("Enter the Room")
}

func LeaveRoom(ctx *fasthttp.RequestCtx) {
	username := JWTAuth(ctx)
	if len(username) == 0 {
		return
	}

	err := services.LeaveRoom(username)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid Room ID")
		l4g.Error("%s LeaveRoom err:%v", username, err)
		return
	}
	l4g.Debug("user:%s LeaveRoom ok!", username)

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString("Left the room")
}

func GetRoom(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.URI().Path()) < len("/room/") {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid Room ID")
		return
	}

	roomId := string(ctx.Request.URI().Path())[len("/room/"):]
	roomName, err := services.GetRoomInfo(roomId)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid Room ID")
		l4g.Error("GetRoomInfo %s err:%v", roomId, err)
		return
	}
	l4g.Debug("GetRoomInfo:%s,%s", roomId, roomName)

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString(roomName)
}

func GetUserList(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.URI().Path()) < len("/room//users") {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid Room ID")
		return
	}

	roomId := string(ctx.Request.URI().Path())[len("/room/"):]
	roomId = roomId[:len(roomId)-len("/users")]

	userList, err := services.GetUserList(roomId)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid Room ID")
		l4g.Error("GetUserList in room:%s err:%v", roomId, err)
		return
	}
	l4g.Debug("GetUserList from roomId:%s,%v", roomId, userList)

	bindata, _ := json.Marshal(userList)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(bindata)
}

func GetRoomList(ctx *fasthttp.RequestCtx) {
	roomControlData := new(models.RoomControlData)
	args := ctx.Request.Body()
	err := json.Unmarshal(args, &roomControlData)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid input")
		l4g.Error("GetRoomList Unmarshal args err:%v", err)
		return
	}

	if roomControlData.PageIndex < 0 || roomControlData.PageSize <= 0 {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("PageIndex/PageSize Error")
		return
	}

	roomList, err := services.GetRoomList(roomControlData.PageIndex, roomControlData.PageSize)
	if err != nil {
		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBodyString(err.Error())
		l4g.Error("GetRoomList err:%v", err)
		return
	}
	l4g.Debug("GetRoomList len:%d", len(roomList))

	bindata, _ := json.Marshal(roomList)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(bindata)
}
