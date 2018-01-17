package db

import (
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx"
	"github.com/oleggator/tp-db/models"
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
		return nil, 404
	}

	var userId int32
	var nickname string
	err = conn.QueryRow(
		`get_user_nick`,
		srcThread.Author,
	).Scan(&userId, &nickname)

	if err != nil {
		return nil, 404
	}

	var (
		existingThread models.Thread
		created        time.Time
	)

	err = conn.QueryRow(
		`select id, authorNickname, created, forumSlug, Thread.message, slug, title, votes from Thread
		where slug=$1`,
		srcThread.Slug,
	).Scan(
		&existingThread.ID,
		&existingThread.Author,
		&created,
		&existingThread.Forum,
		&existingThread.Message,
		&existingThread.Slug,
		&existingThread.Title,
		&existingThread.Votes,
	)
	existingThread.Created = (*strfmt.DateTime)(&created)

	// Thread не существует
	if err != nil {
		tx, _ := conn.Begin()

		var threadId int32
		if srcThread.Created == nil {
			if srcThread.Slug != "" {
				err = tx.QueryRow(`
					insert into Thread (author, forum, message, title, slug, forumSlug, authorNickname)
					values ($1, $2, $3, $4, $5, $6, $7)
					returning id;`,
					userId, forumId, srcThread.Message, srcThread.Title, srcThread.Slug, forumSlug, nickname,
				).Scan(&threadId)
			} else {
				err = tx.QueryRow(`
					insert into Thread (author, forum, message, title, forumSlug, authorNickname)
					values ($1, $2, $3, $4, $5, $6)
					returning id;`,
					userId, forumId, srcThread.Message, srcThread.Title, forumSlug, nickname,
				).Scan(&threadId)
			}

		} else {

			if srcThread.Slug != "" {
				err = tx.QueryRow(`
					insert into Thread (author, created, forum, message, title, slug, forumSlug, authorNickname)
					values ($1, $2, $3, $4, $5, $6, $7, $8)
					returning id;`,
					userId, (*time.Time)(srcThread.Created), forumId, srcThread.Message, srcThread.Title, srcThread.Slug, forumSlug, nickname,
				).Scan(&threadId)
			} else {
				err = tx.QueryRow(`
					insert into Thread (author, created, forum, message, title, forumSlug, authorNickname)
					values ($1, $2, $3, $4, $5, $6, $7)
					returning id;`,
					userId, (*time.Time)(srcThread.Created), forumId, srcThread.Message, srcThread.Title, forumSlug, nickname,
				).Scan(&threadId)
			}

		}

		if err == nil {
			tx.Commit()
			srcThread.Forum = forumSlug
			srcThread.Author = nickname
			srcThread.ID = threadId

			conn.Exec(`
				with s as (
					select $1, about, email, fullname, $2 from "User"
					where id=$3
				)
				insert into ForumUser (slug, about, email, fullname, nickname)
				select * from s
				on conflict do nothing;
			`, forumSlug, nickname, userId)

			return srcThread, 201
		}
		tx.Rollback()
	}

	return &existingThread, 409
}

func GetThreads(slug string, limit int32, sinceString string, desc bool) (threads []models.Thread, status int) {
	var (
		forumId int32
	)
	err := conn.QueryRow(
		`select id from forum
        where slug=$1`,
		slug,
	).Scan(&forumId)

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
			select Thread.id, authorNickname, Thread.created, Thread.message, Thread.title, coalesce(Thread.slug, ''), votes, forumSlug
			from Thread
			where forum=$1 and created %s $2
			order by created %s
			%s
		`, compareString, sorting, limitString)

		rows, err = conn.Query(queryString, forumId, since)

	} else {
		queryString := fmt.Sprintf(`
			select Thread.id, authorNickname, Thread.created, Thread.message, Thread.title, coalesce(Thread.slug, ''), votes, forumSlug
			from Thread
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
		rows.Scan(&thread.ID, &thread.Author, &created, &thread.Message, &thread.Title, &thread.Slug, &thread.Votes, &thread.Forum)

		if slug != nil {
			thread.Slug = *slug
		}

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
			)
			update Thread
			set votes = votes - (select d from delta)
			where id = $1
			returning id, created, forumSlug, authorNickname, message, coalesce(slug, ''), title, votes
		`, threadId, vote.Nickname, vote.Voice).Scan(&thread.ID, &created, &thread.Forum, &thread.Author, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	} else {
		err = tx.QueryRow(`
			with delta as (
				INSERT INTO Vote (author, thread, voice)
					VALUES ((select id from "User" where nickname=$2), (select id from thread where slug = $1), $3)
				ON CONFLICT ON CONSTRAINT unique_author_and_thread
					DO UPDATE SET prevVoice = vote.voice, voice = EXCLUDED.voice
				RETURNING (prevVoice - voice) as d
			)
			update Thread
			set votes = votes - (select d from delta)
			where slug = $1
			returning id, created, forumSlug, authorNickname, message, coalesce(slug, ''), title, votes
		`, threadSlug, vote.Nickname, vote.Voice).Scan(&thread.ID, &created, &thread.Forum, &thread.Author, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	}

	if err != nil {
		tx.Rollback()
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
			select thread.id, authorNickname, thread.created, forumSlug, thread.message, coalesce(thread.slug, ''), thread.title, thread.votes
			from thread
			where thread.id = $1
		`, threadId).Scan(&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	} else {
		err = conn.QueryRow(`
			select thread.id, authorNickname, thread.created, forumSlug, thread.message, coalesce(thread.slug, ''), thread.title, thread.votes
			from thread
			where thread.slug = $1
		`, threadSlug).Scan(&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	}

	if err != nil {
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
			update thread set title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message)
			where thread.id = $3
			returning id, authorNickname, created, forumSlug, coalesce(slug, ''), votes, title, message
		`, threadUpdate.Title, threadUpdate.Message, threadId).Scan(
			&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Slug, &thread.Votes, &thread.Title, &thread.Message)
	} else {
		err = tx.QueryRow(`
			update thread set title=COALESCE(NULLIF($1, ''), title), message=COALESCE(NULLIF($2, ''), message)
			where thread.slug = $3
			returning id, authorNickname, created, forumSlug, coalesce(slug, ''), votes, title, message
		`, threadUpdate.Title, threadUpdate.Message, threadSlug).Scan(
			&thread.ID, &thread.Author, &created, &thread.Forum, &thread.Slug, &thread.Votes, &thread.Title, &thread.Message)
	}

	if err != nil {
		tx.Rollback()
		return nil, 404
	}

	tx.Commit()

	thread.Created = (*strfmt.DateTime)(&created)

	return thread, 200
}
