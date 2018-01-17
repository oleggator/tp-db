package handlers

import (
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/models"
	"github.com/valyala/fasthttp"
)

func ServiceClearPost(ctx *fasthttp.RequestCtx) {
	db.Clear()
}

func ServiceStatusGet(ctx *fasthttp.RequestCtx) {
	status := models.Status{
		Forum:  db.CountForums(),
		Post:   db.CountPosts(),
		Thread: db.CountThreads(),
		User:   db.CountUsers(),
	}

	ctx.SetContentType("application/json")
	json, _ := status.MarshalBinary()
	ctx.Write(json)
}
