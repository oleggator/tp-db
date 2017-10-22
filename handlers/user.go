package handlers

import (
	"github.com/kataras/iris"
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/models"
	"io/ioutil"
)

func UserNicknameCreatePost(ctx iris.Context) {
	body, _ := ioutil.ReadAll(ctx.Request().Body)

	srcUser := models.User{}
	srcUser.UnmarshalBinary(body)
	srcUser.Nickname = ctx.Params().Get("nickname")

	users, ok := db.CreateUser(srcUser)

	ctx.ContentType("application/json")

	if ok {
		json, _ := srcUser.MarshalBinary()

		ctx.StatusCode(201)
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

		ctx.StatusCode(409)
		ctx.Write(body)
	}
}

func UserNicknameProfileGet(ctx iris.Context) {
	user, ok := db.GetUser(ctx.Params().Get("nickname"))

	ctx.ContentType("application/json")

	if ok {
		json, _ := user.MarshalBinary()

		ctx.StatusCode(200)
		ctx.Write(json)
	} else {
		error := models.Error{Message: "Can't find user\n"}
		json, _ := error.MarshalBinary()

		ctx.StatusCode(404)
		ctx.Write(json)
	}
}

func UserNicknameProfilePost(ctx iris.Context) {
	body, _ := ioutil.ReadAll(ctx.Request().Body)

	srcUser := models.User{}
	srcUser.UnmarshalBinary(body)
	srcUser.Nickname = ctx.Params().Get("nickname")

	ctx.ContentType("application/json")

	switch status := db.UpdateUser(srcUser); status {
	case 200:
		json, _ := srcUser.MarshalBinary()

		ctx.StatusCode(200)
		ctx.Write(json)

	case 404:
		error := models.Error{Message: "Can't find user\n"}
		json, _ := error.MarshalBinary()

		ctx.StatusCode(404)
		ctx.Write(json)

	case 409:
		error := models.Error{Message: "New data conflicts with old\n"}
		json, _ := error.MarshalBinary()

		ctx.StatusCode(409)
		ctx.Write(json)

	default:
		ctx.StatusCode(500)
	}

}
