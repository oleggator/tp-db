package db

import (
	"github.com/oleggator/tp-db/models"
	"log"
)

func CreateForum(srcForum models.Forum) (forum models.Forum, status int) {
	var userId int
	conn.QueryRow(
		`select id from "User"
        where lower(nickname)=lower($1::text);`,
		srcForum.User,
	).Scan(&userId)

	log.Println(userId)

	if userId == 0 {
		return models.Forum{}, 404
	}

	ct, err := conn.Exec(`
           insert into Forum (slug, title, moderator)
           values ($1, $2, $3);`,
		srcForum.Slug, srcForum.Title, userId,
	)

	log.Println(err)

	if ct.RowsAffected() > 0 {
		return srcForum, 201
	}

	return models.Forum{}, 404
}
