package handlers

import (
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/models"
	"github.com/valyala/fasthttp"
)

func UserNicknameCreatePost(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	srcUser := models.User{}
	srcUser.UnmarshalBinary(body)
	srcUser.Nickname = ctx.UserValue("nickname").(string)

	users, ok := db.CreateUser(srcUser)

	ctx.SetContentType("application/json")

	if ok {
		json, _ := srcUser.MarshalBinary()

		ctx.SetStatusCode(201)
		ctx.Write(json)
	} else {
		body := []byte("[")
		for i, user := range users {
			json, _ := user.MarshalBinary()
			body = append(body, json...)

			if i != len(users)-1 {
				body = append(body, byte(','))
			}
		}
		body = append(body, byte(']'))

		ctx.SetStatusCode(409)
		ctx.Write(body)
	}
}

func UserNicknameProfileGet(ctx *fasthttp.RequestCtx) {
	user, ok := db.GetUser(ctx.UserValue("nickname").(string))

	ctx.SetContentType("application/json")

	if ok {
		json, _ := user.MarshalBinary()

		ctx.SetStatusCode(200)
		ctx.Write(json)
	} else {
		err := models.Error{Message: "Can't find user\n"}
		json, _ := err.MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}

func UserNicknameProfilePost(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	srcUser := models.User{}
	srcUser.UnmarshalBinary(body)
	srcUser.Nickname = ctx.UserValue("nickname").(string)

	ctx.SetContentType("application/json")

	switch user, status := db.UpdateUser(srcUser); status {
	case 200:
		json, _ := user.MarshalBinary()

		ctx.SetStatusCode(200)
		ctx.Write(json)

	case 404:
		err := models.Error{Message: "Can't find user\n"}
		json, _ := err.MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)

	case 409:
		error := models.Error{Message: "New data conflicts with old\n"}
		json, _ := error.MarshalBinary()

		ctx.SetStatusCode(409)
		ctx.Write(json)

	default:
		ctx.SetStatusCode(500)
	}

}
