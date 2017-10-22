package handlers

import (
	"github.com/kataras/iris"
	"github.com/oleggator/tp-db/models"
)

func ServiceClearPost(ctx iris.Context) {
	error := models.Error{Message: "Not implemented yet"}
	json, _ := error.MarshalBinary()
	ctx.Write(json)
}

func ServiceStatusGet(ctx iris.Context) {
	status := models.Status{
		Forum:  1,
		Post:   1,
		Thread: 1,
		User:   1,
	}

	json, _ := status.MarshalBinary()
	ctx.Write(json)
}
