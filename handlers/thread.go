package handlers

import (
	"bufio"
	"encoding/json"
	"github.com/go-openapi/strfmt"
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/models"
	"github.com/valyala/fasthttp"
	"time"
)

func ThreadSlugOrIdCreatePost(ctx *fasthttp.RequestCtx) {
	var srcPostsStrings []json.RawMessage
	json.Unmarshal(ctx.PostBody(), &srcPostsStrings)

	now := time.Now()
	var srcPosts []models.Post
	for _, postJson := range srcPostsStrings {
		post := models.Post{}
		post.UnmarshalBinary(postJson)
		if post.Created == nil {
			post.Created = (*strfmt.DateTime)(&now)
		}

		srcPosts = append(srcPosts, post)
	}

	ctx.SetContentType("application/json")

	switch posts, status := db.CreatePosts(srcPosts, ctx.UserValue("slug_or_id").(string)); status {
	case 201:
		ctx.SetStatusCode(201)
		ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
			w.Write([]byte("["))
			for i, post := range posts {
				json, _ := post.MarshalBinary()
				w.Write(json)

				if i != len(posts)-1 {
					w.Write([]byte(","))
				}
			}
			w.Write([]byte("]"))
			w.Flush()
		})
	case 404:
		json, _ := (&models.Error{Message: "Can't find thread\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	case 409:
		json, _ := (&models.Error{Message: "Can't find thread\n"}).MarshalBinary()

		ctx.SetStatusCode(409)
		ctx.Write(json)
	}
}

func ThreadSlugOrIdDetailsGet(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	switch thread, status := db.GetThread(ctx.UserValue("slug_or_id").(string)); status {
	case 200:
		json, _ := thread.MarshalBinary()

		ctx.SetStatusCode(200)
		ctx.Write(json)
	case 404:
		json, _ := (&models.Error{Message: "Can't find thread\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}

func ThreadSlugOrIdDetailsPost(ctx *fasthttp.RequestCtx) {
	threadUpdate := models.ThreadUpdate{}
	threadUpdate.UnmarshalBinary(ctx.PostBody())

	ctx.SetContentType("application/json")
	switch thread, status := db.ModifyThread(&threadUpdate, ctx.UserValue("slug_or_id").(string)); status {
	case 200:
		json, _ := thread.MarshalBinary()

		ctx.SetStatusCode(200)
		ctx.Write(json)
	case 404:
		json, _ := (&models.Error{Message: "Can't find thread\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}

func ThreadSlugOrIdPostsGet(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	slug := ctx.UserValue("slug_or_id").(string)
	limit := int32(ctx.QueryArgs().GetUintOrZero("limit"))
	since := ctx.QueryArgs().GetUintOrZero("since")
	desc := string(ctx.QueryArgs().Peek("desc")) == "true"
	sortString := string(ctx.QueryArgs().Peek("sort"))

	switch posts, status := db.GetPosts(slug, limit, since, desc, sortString); status {
	case 200:
		ctx.SetStatusCode(200)
		ctx.SetBodyStreamWriter(func(w *bufio.Writer) {
			w.Write([]byte("["))
			for i, post := range posts {
				json, _ := post.MarshalBinary()
				w.Write(json)

				if i != len(posts)-1 {
					w.Write([]byte(","))
				}
			}
			w.Write([]byte("]"))
			w.Flush()
		})
	case 404:
		json, _ := (&models.Error{Message: "Can't find thread\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}

func ThreadSlugOrIdVotePost(ctx *fasthttp.RequestCtx) {
	vote := models.Vote{}
	vote.UnmarshalBinary(ctx.PostBody())

	ctx.SetContentType("application/json")

	switch thread, status := db.VoteThread(&vote, ctx.UserValue("slug_or_id").(string)); status {
	case 200:
		json, _ := thread.MarshalBinary()

		ctx.SetStatusCode(200)
		ctx.Write(json)
	case 404:
		json, _ := (&models.Error{Message: "Can't find thread\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}
