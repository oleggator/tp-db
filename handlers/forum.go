package handlers

import (
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/models"
	"github.com/valyala/fasthttp"
)

func ForumCreatePost(ctx *fasthttp.RequestCtx) {
	if ctx.UserValue("slug").(string) != "create" {
		ctx.SetStatusCode(405)
		return
	}

	body := ctx.PostBody()

	srcForum := models.Forum{}
	srcForum.UnmarshalBinary(body)

	ctx.SetContentType("application/json")

	switch forum, status := db.CreateForum(&srcForum); status {
	case 201:
		json, _ := srcForum.MarshalBinary()

		ctx.SetStatusCode(201)
		ctx.Write(json)

	case 404:
		json, _ := (&models.Error{Message: "Can't find user\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)

	case 409:
		json, _ := forum.MarshalBinary()

		ctx.SetStatusCode(409)
		ctx.Write(json)

	default:
		ctx.SetStatusCode(500)
	}

}

func ForumSlugCreatePost(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	srcThread := models.Thread{}
	srcThread.UnmarshalBinary(body)
	srcThread.Forum = ctx.UserValue("slug").(string)

	ctx.SetContentType("application/json")

	switch thread, status := db.CreateThread(srcThread); status {
	case 201:
		json, _ := thread.MarshalBinary()

		ctx.SetStatusCode(201)
		ctx.Write(json)

	case 404:
		json, _ := (&models.Error{Message: "Can't find user\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)

	case 409:
		json, _ := thread.MarshalBinary()

		ctx.SetStatusCode(409)
		ctx.Write(json)

	default:
		ctx.SetStatusCode(500)

	}

}

func ForumSlugDetailsGet(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	switch user, status := db.GetForumDetails(ctx.UserValue("slug").(string)); status {
	case 200:
		json, _ := (*user).MarshalBinary()

		ctx.SetStatusCode(200)
		ctx.Write(json)
	case 404:
		json, _ := (&models.Error{Message: "Can't find forum\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}

func ForumSlugThreadsGet(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	limit := int32(ctx.QueryArgs().GetUintOrZero("limit"))
	sinceString := string(ctx.QueryArgs().Peek("since"))
	desc := string(ctx.QueryArgs().Peek("desc")) == "true"

	ctx.SetContentType("application/json")

	switch threads, status := db.GetThreads(slug, limit, sinceString, desc); status {
	case 200:
		body := []byte("[")
		for i, thread := range threads {
			json, _ := thread.MarshalBinary()
			body = append(body, json...)

			if i != len(threads)-1 {
				body = append(body, byte(','))
			}
		}
		body = append(body, byte(']'))

		ctx.SetStatusCode(200)
		ctx.Write(body)

	case 404:
		json, _ := (&models.Error{Message: "Can't find forum\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}

func ForumSlugUsersGet(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	limit := int32(ctx.QueryArgs().GetUintOrZero("limit"))
	sinceString := string(ctx.QueryArgs().Peek("since"))
	desc := string(ctx.QueryArgs().Peek("desc")) == "true"

	ctx.SetContentType("application/json")

	switch users, status := db.GetForumUsers(slug, limit, sinceString, desc); status {
	case 200:
		body := []byte("[")
		for i, user := range users {
			json, _ := user.MarshalBinary()
			body = append(body, json...)

			if i != len(users)-1 {
				body = append(body, byte(','))
			}
		}
		body = append(body, byte(']'))

		ctx.SetStatusCode(200)
		ctx.Write(body)
	case 404:
		json, _ := (&models.Error{Message: "Can't find forum\n"}).MarshalBinary()

		ctx.SetStatusCode(404)
		ctx.Write(json)
	}
}
