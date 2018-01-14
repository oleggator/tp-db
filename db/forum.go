package db

import (
	"github.com/oleggator/tp-db/models"
	"log"
)

func CreateForum(srcForum models.Forum) (forum models.Forum, status int) {
	var userId int32
	var nickname string
	err := conn.QueryRow(
		`select id, nickname from "User"
        where lower(nickname)=lower($1::text);`,
		srcForum.User,
	).Scan(&userId, &nickname)

	if err != nil {
		log.Println("CreateForum", err)
		return models.Forum{}, 404
	}

	ct, _ := conn.Exec(`
           insert into Forum (slug, title, moderator)
           values ($1, $2, $3);`,
		srcForum.Slug, srcForum.Title, userId,
	)

	if ct.RowsAffected() > 0 {
		forum = srcForum
		forum.User = nickname
		return forum, 201
	}

	conn.QueryRow(
		`select slug, title from forum
		where slug=$1;`,
		srcForum.Slug,
	).Scan(&forum.Slug, &forum.Title)
	forum.User = nickname

	return forum, 409
}

func GetForumDetails(slug string) (forum models.Forum, status int) {
	forum = models.Forum{}

	var forumId int32
	err := conn.QueryRow(`
		select forum.id, forum.slug, forum.title, "User".nickname from forum
		join "User" on "User".id=forum.moderator
		where forum.slug=$1
	`, slug).Scan(&forumId, &forum.Slug, &forum.Title, &forum.User)

	if err != nil {
		log.Println("GetForumDetails:", err)
		return forum, 404
	}

	conn.QueryRow(` select count(*) from thread where forum=$1 `, forumId).Scan(&forum.Threads)
	conn.QueryRow(` select count(*) from post where forum=$1 `, forumId).Scan(&forum.Posts)

	return forum, 200
}
