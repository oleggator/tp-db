package handlers

import (
	"github.com/kataras/iris"
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/models"
	"io/ioutil"
)

func ForumCreatePost(ctx iris.Context) {
	body, _ := ioutil.ReadAll(ctx.Request().Body)

	srcForum := models.Forum{}
	srcForum.UnmarshalBinary(body)

	ctx.ContentType("application/json")

	switch forum, status := db.CreateForum(srcForum); status {
	case 201:
		json, _ := forum.MarshalBinary()

		ctx.StatusCode(201)
		ctx.Write(json)
	case 404:
		error := models.Error{Message: "Can't find user\n"}
		json, _ := error.MarshalBinary()

		ctx.StatusCode(404)
		ctx.Write(json)
	default:
		ctx.StatusCode(500)
	}

}

func ForumSlugCreatePost(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func ForumSlugDetailsGet(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func ForumSlugThreadsGet(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func ForumSlugUsersGet(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}
