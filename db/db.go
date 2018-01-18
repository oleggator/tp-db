package db

import (
	"github.com/jackc/pgx"
	"log"
)

var conn *pgx.ConnPool

func InitDB(config pgx.ConnConfig) {
	var err error
	conn, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 50,
	})

	if err != nil {
		log.Panic(err)
	}

	// Forum
	conn.Prepare("get_user", `select id from "User" where nickname=$1`)

	conn.Prepare("get_user_nick", `select id, nickname from "User" where nickname=$1;`)

	conn.Prepare("insert_into_forum", `insert into Forum (slug, title, moderator, moderatorNickname) values ($1, $2, $3, $4);`)

	conn.Prepare("select_forum", `select slug, title from forum where slug=$1;`)

	conn.Prepare("threads_count", `select count(*) from thread where forum=$1`)

	conn.Prepare("posts_count", `select count(*) from post where forum=$1`)

	conn.Prepare("insert_post", `
		insert into Post (author, message, "thread", isEdited, forum, created, parent, parents, root_parent, id, authorNickname, forumSlug)
		values ((select id from "User" where nickname=$1), $2, $3, $4, $5, $6, $7, $8, $9, $10, $1, $11)
    `)
	conn.Prepare("get_forum_details", `
        select id, slug, title, moderatorNickname, threadsCount, postsCount from forum
		where slug=$1
    `)

	// Posts
	conn.Prepare("thread_by_id", `
        select forum, forumSlug from Thread
            where thread.id = $1
    `)
	conn.Prepare("thread_by_slug", `
        select forum, forumSlug, thread.id from Thread
            where thread.slug = $1
    `)
	conn.Prepare("get_ids", `select array_agg(nextval('post_id_seq')::bigint) from generate_series(1,$1)`)

	conn.Prepare("get_parents", `select parents from post where id = $1`)

	conn.Prepare("get_forum_without_id",
		`select forum.slug, forum.title, moderatorNickname from forum
            where forum.id=$1`)

	conn.Prepare("get_thread",
		`select thread.id, authorNickname, thread.created, forumSlug, thread.message, coalesce(thread.slug, ''), thread.title, thread.votes
            where thread.id = $1`)

	conn.Prepare("get_author",
		`select about, email, fullname, nickname from "User"
            where id=$1`)

	conn.Prepare("update_post",
		`update post
        set message=$1, isEdited=TRUE
        where id=$2`)

	conn.Prepare("get_thread_slug",
		`select coalesce(thread.slug, '') from Thread
            where thread.id = $1`)

	conn.Prepare("get_thread_id",
		`select thread.id from Thread
            where thread.slug = $1`)

	conn.Prepare("get_forum_id_slug",
		`select id, slug from forum
        where slug=$1`)
}

func Close() {
	conn.Close()
}
