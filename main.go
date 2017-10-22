package main

import (
	"github.com/jackc/pgx"
	"github.com/kataras/iris"
	"github.com/oleggator/tp-db/db"
	"github.com/oleggator/tp-db/handlers"
)

func main() {
	dbConfig := pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "docker",
		User:     "docker",
		Password: "docker",
	}
	db.InitDB(dbConfig)

	app := iris.New()

	app.Post("/api/forum/create", handlers.ForumCreatePost)
	app.Post("/api/forum/{slug:string}/create", handlers.ForumSlugCreatePost)
	app.Get("/api/forum/{slug:string}/details", handlers.ForumSlugDetailsGet)
	app.Get("/api/forum/{slug:string}/threads", handlers.ForumSlugThreadsGet)
	app.Get("/api/forum/{slug:string}/users", handlers.ForumSlugUsersGet)

	app.Get("/api/post/{id:int}/details", handlers.PostIdDetailsGet)
	app.Post("/api/post/{id:int}/details", handlers.PostIdDetailsPost)

	app.Get("/api/service/status", handlers.ServiceStatusGet)
	app.Post("/api/service/clear", handlers.ServiceClearPost)

	app.Post("/api/thread/{slug_or_id:string}/create", handlers.ThreadSlugOrIdCreatePost)
	app.Get("/api/thread/{slug_or_id:string}/details", handlers.ThreadSlugOrIdDetailsGet)
	app.Post("/api/thread/{slug_or_id:string}/details", handlers.ThreadSlugOrIdDetailsPost)
	app.Get("/api/thread/{slug_or_id:string}/posts", handlers.ThreadSlugOrIdPostsGet)
	app.Post("/api/thread/{slug_or_id:string}/vote", handlers.ThreadSlugOrIdVotePost)

	app.Post("/api/user/{nickname:string}/create", handlers.UserNicknameCreatePost)
	app.Get("/api/user/{nickname:string}/profile", handlers.UserNicknameProfileGet)
	app.Post("/api/user/{nickname:string}/profile", handlers.UserNicknameProfilePost)

	app.Run(iris.Addr(":5000"))
}
