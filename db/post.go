package db

import (
	"context"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/oleggator/tp-db/models"
	"log"

	//"log"
	"strconv"
	"time"
)

func CreatePosts(srcPosts []models.Post, threadSlug string) (posts []models.Post, status int) {
	var (
		forumId   int
		forumSlug string
		threadId  int32
		err       error
	)

	if threadId64, parseErr := strconv.ParseInt(threadSlug, 0, 32); parseErr == nil {
		threadId = int32(threadId64)

		err = conn.QueryRow(`thread_by_id`, threadId).Scan(&forumId, &forumSlug)
	} else {
		err = conn.QueryRow(`thread_by_slug`, threadSlug).Scan(&forumId, &forumSlug, &threadId)
	}

	if err != nil {
		return nil, 404
	}

	if len(srcPosts) == 0 {
		return nil, 201
	}

	batch := conn.BeginBatch()
	for i, _ := range srcPosts {
		batch.Queue(
			`select check_parent($1, $2)`,
			[]interface{}{threadId, srcPosts[i].Parent},
			[]pgtype.OID{pgtype.Int4OID, pgtype.Int8OID},
			[]int16{pgx.BinaryFormatCode},
		)
	}

	err = batch.Send(context.Background(), &pgx.TxOptions{IsoLevel: pgx.ReadUncommitted, AccessMode: pgx.ReadOnly})

	for _, _ = range srcPosts {
		var status int
		err = batch.QueryRowResults().Scan(&status)
		if err != nil || status != 201 {
			//log.Println("CreatePosts: batch scan error:", err)
			batch.Close()
			return nil, 409
		}
	}
	batch.Close()

	postsIds := make([]int64, 0, len(srcPosts))
	err = conn.QueryRow(`get_ids`, len(srcPosts)).Scan(&postsIds)

	batch = conn.BeginBatch()
	for i, _ := range srcPosts {
		if srcPosts[i].Parent != 0 {
			batch.Queue(
				`get_parents`,
				[]interface{}{srcPosts[i].Parent},
				[]pgtype.OID{pgtype.Int8OID},
				[]int16{pgx.BinaryFormatCode},
			)
		}
	}

	err = batch.Send(context.Background(), &pgx.TxOptions{IsoLevel: pgx.ReadUncommitted, AccessMode: pgx.ReadOnly})

	parents := make([][]int64, len(srcPosts))
	for i, _ := range srcPosts {
		if srcPosts[i].Parent != 0 {
			err = batch.QueryRowResults().Scan(&parents[i])
		}

		parents[i] = append(parents[i], postsIds[i])
	}
	batch.Close()

	tx, _ := conn.Begin()
	batch = tx.BeginBatch()
	for i, _ := range srcPosts {
		srcPosts[i].ID = postsIds[i]
		srcPosts[i].Thread = threadId
		srcPosts[i].Forum = forumSlug

		batch.Queue(
			`insert_post`,
			[]interface{}{srcPosts[i].Author, srcPosts[i].Message, threadId, srcPosts[i].IsEdited, forumId,
				srcPosts[i].Created, srcPosts[i].Parent, parents[i], parents[i][0], srcPosts[i].ID},
			nil,
			nil,
		)
	}

	err = batch.Send(context.Background(), nil)

	_, err = batch.ExecResults()
	if err != nil {
		log.Println("Batch:", err)
		batch.Close()
		tx.Rollback()
		return nil, 404
	}

	batch.Close()
	tx.Commit()

	return srcPosts, 201
}

func GetPost(postId int64, withAuthor bool, withThread bool, withForum bool) (postInfo *models.PostFull, status int) {
	postInfo = &models.PostFull{}

	post := models.Post{}
	postInfo.Post = &post
	post.ID = postId

	var (
		created time.Time
		forumId int32
		userId  int32
	)
	err := conn.QueryRow(`
        select "User".nickname, post.created, forum.slug, post.isEdited, post.message, thread.id, forum.id, "User".id, coalesce(post.parent, 0)
        from Post
        join "User" on "User".id = post.author
        join forum on forum.id = post.forum
        join thread on thread.id = post."thread"
        where post.id = $1
    `, postId).Scan(&post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Thread, &forumId, &userId, &post.Parent)

	if err != nil {
		//log.Println("GetPost:", err)
		return nil, 404
	}

	post.Created = (*strfmt.DateTime)(&created)

	if withForum {
		forum := models.Forum{}
		postInfo.Forum = &forum

		err := conn.QueryRow(`
            select forum.slug, forum.title, "User".nickname from forum
            join "User" on "User".id=forum.moderator
            where forum.id=$1
        `, forumId).Scan(&forum.Slug, &forum.Title, &forum.User)

		if err != nil {
			//log.Println("GetPost:", err)
			return nil, 404
		}

		conn.QueryRow(` select count(*) from thread where forum=$1 `, forumId).Scan(&forum.Threads)
		conn.QueryRow(` select count(*) from post where forum=$1 `, forumId).Scan(&forum.Posts)
	}

	if withThread {
		var created time.Time
		thread := models.Thread{}
		postInfo.Thread = &thread

		err = conn.QueryRow(`
            select thread.id, "User".nickname, thread.created, forum.slug, thread.message, coalesce(thread.slug, ''), thread.title, thread.votes
            from thread
            join "User" on "User".id = thread.author
            join forum on forum.id = thread.forum
            where thread.id = $1
        `, post.Thread).Scan(&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)

		if err != nil {
			//log.Println("GetPost:", err)
			return nil, 404
		}

		thread.Created = (*strfmt.DateTime)(&created)
	}

	if withAuthor {
		user := models.User{}
		postInfo.Author = &user

		err = conn.QueryRow(`
            select about, email, fullname, nickname from "User"
            where id=$1;
        `, userId).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

		if err != nil {
			//log.Println("GetPost:", err)
			return nil, 404
		}
	}

	return postInfo, 200
}

func ModifyPost(postUpdate *models.PostUpdate, postId int64) (post *models.Post, status int) {
	post = &models.Post{}
	post.ID = postId

	var created time.Time
	err := conn.QueryRow(`
        select "User".nickname, post.created, forum.slug, post.isEdited, post.message, thread.id, coalesce(post.parent, 0)
        from Post
        join "User" on "User".id = post.author
        join forum on forum.id = post.forum
        join thread on thread.id = post."thread"
        where post.id = $1
    `, postId).Scan(&post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Thread, &post.Parent)

	if err != nil {
		//log.Println("ModifyPost:", err)
		return nil, 404
	}

	post.Created = (*strfmt.DateTime)(&created)

	if postUpdate.Message == "" || post.Message == postUpdate.Message {
		return post, 200
	}

	tx, _ := conn.Begin()
	_, err = tx.Exec(`
        update post set message=$1, isEdited=TRUE where id=$2
    `, postUpdate.Message, postId)

	post.Message = postUpdate.Message
	post.IsEdited = true

	if err != nil {
		tx.Rollback()
		return nil, 404
	}
	tx.Commit()

	return post, 200
}

func GetPosts(threadSlug string, limit int32, since int, desc bool, sortString string) (posts []models.Post, status int) {
	var (
		threadId int32
		err      error
	)

	if threadId64, parseErr := strconv.ParseInt(threadSlug, 0, 32); parseErr == nil {
		threadId = int32(threadId64)

		err = conn.QueryRow(`
            select coalesce(thread.slug, '') from Thread
            where thread.id = $1
        `, threadId).Scan(&threadSlug)
	} else {
		err = conn.QueryRow(`
            select thread.id from Thread
            where thread.slug = $1
        `, threadSlug).Scan(&threadId)
	}

	if err != nil {
		//log.Println("GetPosts:", err)
		return nil, 404
	}

	var limitString string
	if limit != 0 {
		limitString = fmt.Sprintf("limit %d", limit)
	} else {
		limitString = ""
	}

	posts = make([]models.Post, 0)
	switch sortString {
	case "tree":
		var sorting string
		var compareString string
		if desc {
			sorting = "desc"
			if since != 0 {
				compareString = fmt.Sprintf(" and parents < (select parents from post where id = %d)", since)
			}
		} else {
			sorting = "asc"
			if since != 0 {
				compareString = fmt.Sprintf(" and parents > (select parents from post where id = %d)", since)
			}
		}

		query := fmt.Sprintf(`
            select post.id, "User".nickname, post.created, forum.slug, post.isEdited, post.message, thread.id, coalesce(post.parent, 0)
            from Post
            join "User" on "User".id = post.author
            join forum on forum.id = post.forum
            join thread on thread.id = post."thread"
            where post."thread" = $1 %s
            order by parents %s
            %s
        `, compareString, sorting, limitString)

		rows, err := conn.Query(query, threadId)
		if err != nil {
			//log.Println("GetPosts:", err)
			return nil, 404
		}

		for rows.Next() {
			post := models.Post{}

			var created time.Time
			rows.Scan(&post.ID, &post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Thread, &post.Parent)
			post.Created = (*strfmt.DateTime)(&created)

			posts = append(posts, post)
		}

		return posts, 200

	case "parent_tree":
		var sorting string
		var compareString string
		if desc {
			sorting = "desc"
			if since != 0 {
				compareString = fmt.Sprintf(" and root_parent < (select root_parent from post where id = %d)", since)
			}
		} else {
			sorting = "asc"
			if since != 0 {
				compareString = fmt.Sprintf(" and root_parent > (select root_parent from post where id = %d)", since)
			}
		}

		query := fmt.Sprintf(`
            select post.id, "User".nickname, post.created, forum.slug, post.isEdited, post.message, thread.id, coalesce(post.parent, 0)
            from Post
            join "User" on "User".id = post.author
            join forum on forum.id = post.forum
            join thread on thread.id = post."thread"
			join (
				select id from post
				where parent=0 and post."thread"=$1 %s
				order by id %s
				%s
			) selectedParents
			on root_parent=selectedParents.id
 			order by parents %s
		`, compareString, sorting, limitString, sorting)

		rows, err := conn.Query(query, threadId)
		if err != nil {
			//log.Println("GetPosts:", err)
			return nil, 404
		}

		for rows.Next() {
			post := models.Post{}

			var created time.Time
			rows.Scan(&post.ID, &post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Thread, &post.Parent)
			post.Created = (*strfmt.DateTime)(&created)

			posts = append(posts, post)
		}

		return posts, 200

	default:
		var sorting string
		var compareString string
		if desc {
			sorting = "desc"
			if since != 0 {
				compareString = fmt.Sprintf(" and post.id < %d", since)
			}
		} else {
			sorting = "asc"
			if since != 0 {
				compareString = fmt.Sprintf(" and post.id > %d", since)
			}
		}

		query := fmt.Sprintf(`
            select post.id, "User".nickname, post.created, forum.slug, post.isEdited, post.message, thread.id, coalesce(post.parent, 0)
            from Post
            join "User" on "User".id = post.author
            join forum on forum.id = post.forum
            join thread on thread.id = post."thread"
            where post."thread" = $1 %s
            order by post.id %s
            %s
        `, compareString, sorting, limitString)

		rows, err := conn.Query(query, threadId)
		if err != nil {
			//log.Println("GetPosts:", err)
			return nil, 404
		}

		for rows.Next() {
			post := models.Post{}

			var created time.Time
			rows.Scan(&post.ID, &post.Author, &created, &post.Forum, &post.IsEdited, &post.Message, &post.Thread, &post.Parent)
			post.Created = (*strfmt.DateTime)(&created)

			posts = append(posts, post)
		}

		return posts, 200
	}

	return posts, 200
}
