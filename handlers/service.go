package handlers

import (
	"fmt"
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/models"
	"github.com/valyala/fasthttp"
)

func ServiceClearPost(ctx *fasthttp.RequestCtx) {
	fmt.Println("ServiceClearPost")
	db.Clear()
}

func ServiceStatusGet(ctx *fasthttp.RequestCtx) {
	fmt.Println("ServiceStatusGet")

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
