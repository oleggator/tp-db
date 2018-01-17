package db

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"github.com/oleggator/tp-db/models"
	"log"
	"strconv"
	"time"
)

func CreateThread(srcThread *models.Thread) (threadNew *models.Thread, status int) {
	var (
		forumSlug string
		forumId   int32
	)
	err := conn.QueryRow(
		`get_forum_id_slug`,
		srcThread.Forum,
	).Scan(&forumId, &forumSlug)

	if err != nil {
		//log.Println("CreateThread:", err)
		return nil, 404
	}

	var userId int32
	var nickname string
	err = conn.QueryRow(
		`get_user_nick`,
		srcThread.Author,
	).Scan(&userId, &nickname)

	if err != nil {
		//log.Println("CreateThread:", err)
		return nil, 404
	}

	var (
		existingThread  models.Thread
		authorId        int32
		existingForumId int32
		created         time.Time
	)

	err = conn.QueryRow(
		`select id, author, created, forum, Thread.message, slug, title, votes from Thread
		where slug=$1`,
		srcThread.Slug,
	).Scan(
		&existingThread.ID,
		&authorId,
		&created,
		&existingForumId,
		&existingThread.Message,
		&existingThread.Slug,
		&existingThread.Title,
		&existingThread.Votes,
	)
	existingThread.Created = (*strfmt.DateTime)(&created)

	tx, _ := conn.Begin()

	// Thread не существует
	if err != nil {
		var threadId int32
		if srcThread.Created == nil {
			if srcThread.Slug != "" {
				err = tx.QueryRow(`
					insert into Thread (author, forum, message, title, slug)
					values ($1, $2, $3, $4, $5)
					returning id;`,
					userId, forumId, srcThread.Message, srcThread.Title, srcThread.Slug,
				).Scan(&threadId)
			} else {
				err = tx.QueryRow(`
					insert into Thread (author, forum, message, title)
					values ($1, $2, $3, $4)
					returning id;`,
					userId, forumId, srcThread.Message, srcThread.Title,
				).Scan(&threadId)
			}

		} else {

			if srcThread.Slug != "" {
				err = tx.QueryRow(`
					insert into Thread (author, created, forum, message, title, slug)
					values ($1, $2, $3, $4, $5, $6)
					returning id;`,
					userId, (*time.Time)(srcThread.Created), forumId, srcThread.Message, srcThread.Title, srcThread.Slug,
				).Scan(&threadId)
			} else {
				err = tx.QueryRow(`
					insert into Thread (author, created, forum, message, title)
					values ($1, $2, $3, $4, $5)
					returning id;`,
					userId, (*time.Time)(srcThread.Created), forumId, srcThread.Message, srcThread.Title,
				).Scan(&threadId)
			}

		}

		if err == nil {
			tx.Commit()
			srcThread.Forum = forumSlug
			srcThread.Author = nickname
			srcThread.ID = threadId

			return srcThread, 201
		}
		tx.Rollback()
	}

	err = conn.QueryRow(
		`select slug from Forum
		where id=$1;`,
		existingForumId,
	).Scan(&existingThread.Forum)

	//log.Println(err)

	err = conn.QueryRow(
		`select nickname from "User"
		where id=$1;`,
		authorId,
	).Scan(&existingThread.Author)
	//if err != nil {
	//	log.Println("CreateThread:", err)
	//}

	return &existingThread, 409
}

func GetThreads(slug string, limit int32, sinceString string, desc bool) (threads []models.Thread, status int) {
	var (
		forumSlug string
		forumId   int32
	)
	err := conn.QueryRow(
		`select id, slug from forum
        where slug=$1`,
		slug,
	).Scan(&forumId, &forumSlug)

	if err != nil {
		return nil, 404
	}

	var sorting string
	var compareString string
	if desc {
		sorting = "desc"
		compareString = "<="
	} else {
		sorting = "asc"
		compareString = ">="
	}

	var limitString string
	if limit != 0 {
		limitString = fmt.Sprintf("limit %d", limit)
	} else {
		limitString = ""
	}

	var rows *pgx.Rows
	threads = make([]models.Thread, 0)
	if sinceString != "" {
		since, _ := time.Parse(time.RFC3339, sinceString)

		queryString := fmt.Sprintf(`
			select Thread.id, "User".nickname, Thread.created, Thread.message, Thread.title, coalesce(Thread.slug, '')
			from Thread
			join "User" on "User".id = Thread.author
			where forum=$1 and created %s $2
			order by created %s
			%s
		`, compareString, sorting, limitString)

		rows, err = conn.Query(queryString, forumId, since)

	} else {
		queryString := fmt.Sprintf(`
			select Thread.id, "User".nickname, Thread.created, Thread.message, Thread.title, coalesce(Thread.slug, '')
			from Thread
			join "User" on "User".id = Thread.author
			where forum=$1
			order by created %s
			%s
		`, sorting, limitString)
		rows, err = conn.Query(queryString, forumId)
	}

	for rows.Next() {
		thread := models.Thread{}

		var created time.Time
		var slug *string
		rows.Scan(&thread.ID, &thread.Author, &created, &thread.Message, &thread.Title, &thread.Slug)

		if slug != nil {
			thread.Slug = *slug
		}

		thread.Forum = forumSlug
		thread.Created = (*strfmt.DateTime)(&created)
		threads = append(threads, thread)
	}

	return threads, 200
}

func VoteThread(vote *models.Vote, threadSlug string) (thread *models.Thread, status int) {
	thread = &models.Thread{}
	var created time.Time

	tx, _ := conn.Begin()

	threadId, err := strconv.ParseInt(threadSlug, 0, 32)
	if err == nil {
		err = tx.QueryRow(`
			with delta as (
				INSERT INTO Vote (author, thread, voice)
					VALUES ((select id from "User" where nickname=$2), $1, $3)
				ON CONFLICT ON CONSTRAINT unique_author_and_thread
					DO UPDATE SET prevVoice = vote.voice, voice = EXCLUDED.voice
				RETURNING (prevVoice - voice) as d
			), s as (
				update Thread
				set votes = votes - (select d from delta)
				where id = $1
				returning id, created, forum, author, message, slug, title, votes
			)
			select s.id, "User".nickname, s.created, forum.slug, s.message, coalesce(s.slug, ''), s.title, s.votes
			from s
			join "User" on "User".id = s.author
			join forum on forum.id = s.forum
		`, threadId, vote.Nickname, vote.Voice).Scan(&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	} else {
		err = tx.QueryRow(`
			with delta as (
				INSERT INTO Vote (author, thread, voice)
					VALUES ((select id from "User" where nickname=$2), (select id from thread where slug = $1), $3)
				ON CONFLICT ON CONSTRAINT unique_author_and_thread
					DO UPDATE SET prevVoice = vote.voice, voice = EXCLUDED.voice
				RETURNING (prevVoice - voice) as d
			), s as (
				update Thread
				set votes = votes - (select d from delta)
				where slug = $1
				returning id, created, forum, author, message, slug, title, votes
			)
			select s.id, "User".nickname, s.created, forum.slug, s.message, coalesce(s.slug, ''), s.title, s.votes
			from s
			join "User" on "User".id = s.author
			join forum on forum.id = s.forum
		`, threadSlug, vote.Nickname, vote.Voice).Scan(&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	}

	if err != nil {
		tx.Rollback()
		log.Println("VoteThread: getThreadId:", err)
		return nil, 404
	}
	tx.Commit()

	thread.Created = (*strfmt.DateTime)(&created)

	return thread, 200
}

func GetThread(threadSlug string) (thread *models.Thread, status int) {
	thread = &models.Thread{}
	var created time.Time

	threadId, err := strconv.ParseInt(threadSlug, 0, 32)
	if err == nil {
		err = conn.QueryRow(`
			select thread.id, "User".nickname, thread.created, forum.slug, thread.message, coalesce(thread.slug, ''), thread.title, thread.votes
			from thread
			join "User" on "User".id = thread.author
			join forum on forum.id = thread.forum
			where thread.id = $1
		`, threadId).Scan(&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	} else {
		err = conn.QueryRow(`
			select thread.id, "User".nickname, thread.created, forum.slug, thread.message, coalesce(thread.slug, ''), thread.title, thread.votes
			from thread
			join "User" on "User".id = thread.author
			join forum on forum.id = thread.forum
			where thread.slug = $1
		`, threadSlug).Scan(&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	}

	if err != nil {
		//log.Println(err)
		return nil, 404
	}

	thread.Created = (*strfmt.DateTime)(&created)

	return thread, 200
}

func ModifyThread(threadUpdate *models.ThreadUpdate, threadSlug string) (thread *models.Thread, status int) {
	thread = &models.Thread{}
	var created time.Time

	tx, _ := conn.Begin()
	threadId, err := strconv.ParseInt(threadSlug, 0, 32)
	if err == nil {
		err = tx.QueryRow(`
			with updatedThread as (
				update thread set title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message)
				where thread.id = $3
				returning id, author, created, forum, slug, votes, title, message
			)
			select updatedThread.id, "User".nickname, updatedThread.created, forum.slug, coalesce(updatedThread.slug, ''), updatedThread.votes, updatedThread.title, updatedThread.message
			from updatedThread
			join "User" on "User".id = updatedThread.author
			join forum on forum.id = updatedThread.forum
		`, threadUpdate.Title, threadUpdate.Message, threadId).Scan(
			&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Slug, &thread.Votes, &thread.Title, &thread.Message)
	} else {
		err = tx.QueryRow(`
			with updatedThread as (
				update thread set title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message)
				where thread.slug = $3
				returning id, author, created, forum, slug, votes, title, message
			)
			select updatedThread.id, "User".nickname, updatedThread.created, forum.slug, coalesce(updatedThread.slug, ''), updatedThread.votes, updatedThread.title, updatedThread.message
			from updatedThread
			join "User" on "User".id = updatedThread.author
			join forum on forum.id = updatedThread.forum
		`, threadUpdate.Title, threadUpdate.Message, threadSlug).Scan(
			&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Slug, &thread.Votes, &thread.Title, &thread.Message)
	}

	if err != nil {
		tx.Rollback()
		//log.Println(err)
		return nil, 404
	}

	tx.Commit()

	thread.Created = (*strfmt.DateTime)(&created)

	return thread, 200
}
