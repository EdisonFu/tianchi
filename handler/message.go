package handler

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
	"tianchi/models"
	"tianchi/services"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func SendMessage(ctx *fasthttp.RequestCtx) {
	username := JWTAuth(ctx)
	if len(username) == 0 {
		return
	}

	message := new(models.Message)
	args := ctx.Request.Body()
	err := json.Unmarshal(args, &message)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid input")
		l4g.Error("SendMessage Unmarshal args err:%v", err)
		return
	}

	if message == nil || message.Id == "" {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid input")
		l4g.Error("SendMessage message illegal:%v", message)
		return
	}

	message.TimeStamp = strconv.Itoa(int(time.Now().Unix()))
	err = services.SendMessage(username, message)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid input")
		l4g.Error("%s SendMessage err:%v", username, err)
		return
	}
	l4g.Debug("%s SendMessage:%s", username, message.Text)

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString("successful operation")
}

func RetrieveMessage(ctx *fasthttp.RequestCtx) {
	username := JWTAuth(ctx)
	if len(username) == 0 {
		return
	}

	messageControlData := new(models.MessageControlData)
	args := ctx.Request.Body()
	err := json.Unmarshal(args, &messageControlData)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("Invalid input")
		l4g.Error("RetrieveMessage Unmarshal args err:%v", err)
		return
	}

	if messageControlData.PageIndex >= 0 || messageControlData.PageSize <= 0 {
		ctx.SetStatusCode(http.StatusNotFound)
		ctx.SetBodyString("PageIndex/PageSize Error")
		return
	}

	messageList, err := services.RetrieveMessage(username, messageControlData.PageIndex, messageControlData.PageSize)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid input")
		l4g.Error("%s RetrieveMessage err:%v", username, err)
		return
	}
	l4g.Debug("%s RetrieveMessage len:%d", username, len(messageList))

	bindata, _ := json.Marshal(messageList)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(bindata)
}
