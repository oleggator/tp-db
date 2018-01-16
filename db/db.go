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

	conn.Prepare("insert_into_forum", `insert into Forum (slug, title, moderator) values ($1, $2, $3);`)

	conn.Prepare("select_forum", `select slug, title from forum where slug=$1;`)

	conn.Prepare("threads_count", `select count(*) from thread where forum=$1`)

	conn.Prepare("posts_count", `select count(*) from post where forum=$1`)

	conn.Prepare("insert_post", `
        insert into Post (author, message, "thread", isEdited, forum, created, parent, parents, root_parent, id)
        values ((select id from "User" where nickname=$1), $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `)
	conn.Prepare("get_forum_details", `
        select forum.id, forum.slug, forum.title, "User".nickname from forum
		join "User" on "User".id=forum.moderator
		where forum.slug=$1
    `)

	// Posts
	conn.Prepare("thread_by_id", `
        select forum.id, forum.slug from Thread
            join Forum on thread.forum = forum.id
            where thread.id = $1
    `)
	conn.Prepare("thread_by_slug", `
        select forum.id, forum.slug, thread.id from Thread
            join Forum on thread.forum = forum.id
            where thread.slug = $1
    `)
	conn.Prepare("get_ids", `select array_agg(nextval('post_id_seq')::bigint) from generate_series(1,$1)`)

	conn.Prepare("get_parents", `select parents from post where id = $1`)

	conn.Prepare("get_post",
		`select "User".nickname, post.created, forum.slug, post.isEdited, post.message, thread.id
        from Post
        join "User" on "User".id = post.author
        join forum on forum.id = post.forum
        join thread on thread.id = post."thread"
        where post.id = $1`)

	conn.Prepare("get_forum_without_id",
		`select forum.slug, forum.title, "User".nickname from forum
            join "User" on "User".id=forum.moderator
            where forum.id=$1`)

	conn.Prepare("get_thread",
		`select thread.id, "User".nickname, thread.created, forum.slug, thread.message, coalesce(thread.slug, ''), thread.title, thread.votes
            from thread
            join "User" on "User".id = thread.author
            join forum on forum.id = thread.forum
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
