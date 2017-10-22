package handlers

import (
	"github.com/kataras/iris"
	"github.com/oleggator/tp-db/models"
)

func PostIdDetailsGet(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func PostIdDetailsPost(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}
