package handler

import (
	l4g "github.com/alecthomas/log4go"
	"github.com/valyala/fasthttp"
	"net/http"
	"tianchi/models"
	"tianchi/services"
)

func CreateUser(ctx *fasthttp.RequestCtx) {
	user := new(models.User)
	args := ctx.Request.Body()
	err := json.Unmarshal(args, &user)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid input")
		l4g.Error("CreateUser Unmarshal args err:%v", err)
		return
	}

	err = services.CreateUser(user)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid input")
		l4g.Error("CreateUser err:%v", err)
		return
	}
	l4g.Debug("CreateUser:%v", user)

	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString("successful operation")
}

func UserLogin(ctx *fasthttp.RequestCtx) {
	user := new(models.User)
	args := ctx.QueryArgs()
	user.Username = string(args.Peek("username"))
	user.Password = string(args.Peek("password"))

	ok := services.Login(user)
	if !ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid username or password.")
		return
	}
	l4g.Debug("UserLogin:%v success!", user.Username)

	token, err := services.CreateToken(user)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid username supplied")
		l4g.Error("CreateToken err:%v", err)
		return
	}
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBodyString(token)
}

func GetUser(ctx *fasthttp.RequestCtx) {
	if len(ctx.Request.URI().Path()) < len("/user/") {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid username supplied")
		return
	}

	username := string(ctx.Request.URI().Path())[len("/user/"):]
	userResponse, err := services.GetUserInfo(username)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("Invalid username supplied")
		l4g.Error("GetUserInfo about %s err:%v", username, err)
		return
	}
	l4g.Debug("GetUser:%v by name:%v", userResponse, username)

	binData, _ := json.Marshal(userResponse)
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetBody(binData)
}
