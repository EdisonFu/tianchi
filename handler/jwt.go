package handler

import (
	"github.com/valyala/fasthttp"
	"net/http"
	"strings"
	"tianchi/dao/cache"
)

func JWTAuth(ctx *fasthttp.RequestCtx) string {
	token := string(ctx.Request.Header.Peek("Authorization"))
	token = strings.Replace(token, "Bearer ", "", 1)
	if token == "" {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("No authorization information in the headers")
		return ""
	}

	username, ok := cache.UserToken.Load(token)
	if !ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		ctx.SetBodyString("authorization information ERR in the headers")
		return ""
	}

	return username.(string)
}
