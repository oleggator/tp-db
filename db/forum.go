package db

import (
	"github.com/oleggator/tp-db/models"
)

func CreateForum(srcForum *models.Forum) (forum *models.Forum, status int) {
	var userId int32
	err := conn.QueryRow(
		`get_user_nick`,
		srcForum.User,
	).Scan(&userId, &srcForum.User)

	if err != nil {
		//log.Println("CreateForum", err)
		return nil, 404
	}

	tx, _ := conn.Begin()
	_, err = tx.Exec(`insert_into_forum`,
		srcForum.Slug, srcForum.Title, userId,
	)

	if err == nil {
		tx.Commit()
		return srcForum, 201
	}
	//log.Println("CreateForum:", err)
	tx.Rollback()

	conn.QueryRow(
		`select_forum`,
		srcForum.Slug,
	).Scan(&srcForum.Slug, &srcForum.Title)

	return srcForum, 409
}

func GetForumDetails(slug string) (forum *models.Forum, status int) {
	forum = &models.Forum{}

	var forumId int32
	err := conn.QueryRow(`get_forum_details`, slug).Scan(&forumId, &forum.Slug, &forum.Title, &forum.User)

	if err != nil {
		return forum, 404
	}

	conn.QueryRow(`threads_count`, forumId).Scan(&forum.Threads)
	conn.QueryRow(`posts_count`, forumId).Scan(&forum.Posts)

	return forum, 200
}
