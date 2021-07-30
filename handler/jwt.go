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
		ctx.SetStatusCode(http.StatusUnauthorized)
		ctx.SetBodyString("No authorization information in the headers")
		return ""
	}

	username, ok := cache.UserToken.Load(token)
	if !ok {
		ctx.SetStatusCode(http.StatusUnauthorized)
		ctx.SetBodyString("authorization information ERR in the headers")
		return ""
	}

	return username.(string)

	//claims, err := services.VerifyToken(token)
	//if err != nil {
	//	if err.Error() == "TokenExpire" {
	//		ctx.SetStatusCode(http.StatusUnauthorized)
	//		ctx.Response.Header.Set("message", "Authorization has expired")
	//		return ""
	//	}
	//	ctx.SetStatusCode(http.StatusUnauthorized)
	//	ctx.SetBodyString(err.Error())
	//	return ""
	//}
	//
	//return claims.UserId
}
