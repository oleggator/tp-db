package main

import (
	"github.com/jackc/pgx"
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/handlers"
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	db.InitDB(pgx.ConnConfig{
		//Host:     "localhost",
		Host:     "/var/run/postgresql/",
		Port:     5432,
		Database: "docker",
		User:     "docker",
		Password: "docker",
	})

	defer db.Close()

	router := fasthttprouter.New()

	router.POST("/api/forum/:slug", handlers.ForumCreatePost) // "/api/forum/create"
	router.POST("/api/forum/:slug/create", handlers.ForumSlugCreatePost)
	router.GET("/api/forum/:slug/details", handlers.ForumSlugDetailsGet)
	router.GET("/api/forum/:slug/threads", handlers.ForumSlugThreadsGet)
	router.GET("/api/forum/:slug/users", handlers.ForumSlugUsersGet)

	router.GET("/api/post/:id/details", handlers.PostIdDetailsGet)
	router.POST("/api/post/:id/details", handlers.PostIdDetailsPost)

	router.GET("/api/service/status", handlers.ServiceStatusGet)
	router.POST("/api/service/clear", handlers.ServiceClearPost)

	router.POST("/api/thread/:slug_or_id/create", handlers.ThreadSlugOrIdCreatePost)
	router.GET("/api/thread/:slug_or_id/details", handlers.ThreadSlugOrIdDetailsGet)
	router.POST("/api/thread/:slug_or_id/details", handlers.ThreadSlugOrIdDetailsPost)
	router.GET("/api/thread/:slug_or_id/posts", handlers.ThreadSlugOrIdPostsGet)
	router.POST("/api/thread/:slug_or_id/vote", handlers.ThreadSlugOrIdVotePost)

	router.POST("/api/user/:nickname/create", handlers.UserNicknameCreatePost)
	router.GET("/api/user/:nickname/profile", handlers.UserNicknameProfileGet)
	router.POST("/api/user/:nickname/profile", handlers.UserNicknameProfilePost)

	err := fasthttp.ListenAndServe(":5000", router.Handler)
	if err != nil {
		log.Println(err)
	}
}
