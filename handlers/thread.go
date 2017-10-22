package handlers

import (
	"github.com/kataras/iris"
	"github.com/oleggator/tp-db/models"
)

func ThreadSlugOrIdCreatePost(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func ThreadSlugOrIdDetailsGet(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func ThreadSlugOrIdDetailsPost(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func ThreadSlugOrIdPostsGet(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func ThreadSlugOrIdVotePost(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}
