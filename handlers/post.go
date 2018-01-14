package handlers

import (
	"bytes"
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/models"
	"github.com/valyala/fasthttp"
	"strconv"
)

func PostIdDetailsGet(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	postId, _ := strconv.ParseInt(ctx.UserValue("id").(string), 0, 32)

	withForum := false
	withThread := false
	withAuthor := false

	if ctx.QueryArgs().Has("related") {
		for _, i := range bytes.Split(ctx.QueryArgs().Peek("related"), []byte(",")) {
			switch string(i) {
			case "user":
				withAuthor = true
			case "thread":
				withThread = true
			case "forum":
				withForum = true
			}
		}
	}

	switch postFull, status := db.GetPost(postId, withAuthor, withThread, withForum); status {
	case 200:
		json, _ := postFull.MarshalBinary()

		ctx.SetStatusCode(200)
		ctx.Write(json)
	case 404:
		json, _ := (&models.Error{Message: "Can't find thread\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}

func PostIdDetailsPost(ctx *fasthttp.RequestCtx) {
	postUpdate := models.PostUpdate{}
	postUpdate.UnmarshalBinary(ctx.PostBody())

	postId, _ := strconv.ParseInt(ctx.UserValue("id").(string), 0, 64)

	ctx.SetContentType("application/json")

	switch post, status := db.ModifyPost(postUpdate, postId); status {
	case 200:
		json, _ := post.MarshalBinary()

		ctx.SetStatusCode(200)
		ctx.Write(json)
	case 404:
		json, _ := (&models.Error{Message: "Can't find post\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}
