package db

import (
	"fmt"
	"github.com/jackc/pgx"
	"github.com/oleggator/tp-db/models"
)

func CreateUser(profile *models.User) (users []models.User, ok bool) {
	tx, _ := conn.Begin()
	ct, _ := tx.Exec(`
           insert into "User" (about, email, fullname, nickname)
           values ($1, $2, $3, $4);`,
		profile.About, string(profile.Email), profile.Fullname, profile.Nickname,
	)

	if ct.RowsAffected() > 0 {
		tx.Commit()
		return nil, true
	}
	tx.Rollback()

	users = make([]models.User, 0)
	rows, _ := conn.Query(
		`select about, email, fullname, nickname from "User"
        where nickname=$1 or email=$2`,
		profile.Nickname, string(profile.Email),
	)
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

		users = append(users, user)
	}

	return users, false
}

func GetUser(nickname string) (user *models.User, ok bool) {
	user = &models.User{}
	err := conn.QueryRow(
		`select about, email, fullname, nickname from "User"
        where nickname=$1`,
		nickname,
	).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	return user, err == nil
}

func UpdateUser(srcUser *models.User) (user *models.User, status int) {
	user = &models.User{Nickname: srcUser.Nickname}
	var id int32
	err := conn.QueryRow(
		`select id, about, email, fullname from "User"
        where nickname=$1`,
		srcUser.Nickname,
	).Scan(&id, &user.About, &user.Email, &user.Fullname)

	if err != nil {
		return nil, 404
	}

	if srcUser.About == "" && srcUser.Email == "" && srcUser.Fullname == "" {
		return user, 200
	}

	tx, _ := conn.Begin()
	ct, _ := tx.Exec(
		`update "User" set
			about=COALESCE(NULLIF($1, ''), about),
			email=COALESCE(NULLIF($2::text, ''), email),
			fullname=COALESCE(NULLIF($3, ''), fullname),
			nickname=COALESCE(NULLIF($4::text, ''), nickname)
		where id=$5;`,
		srcUser.About,
		srcUser.Email,
		srcUser.Fullname,
		srcUser.Nickname,
		id,
	)

	if ct.RowsAffected() == 0 {
		tx.Rollback()
		return nil, 409
	}
	tx.Commit()

	conn.QueryRow(
		`select about, email, fullname from "User"
        where id=$1;`,
		id,
	).Scan(&user.About, &user.Email, &user.Fullname)

	return user, 200
}

func GetForumUsers(slug string, limit int32, sinceNickname string, desc bool) (users []models.User, status int) {
	var forumId int32
	err := conn.QueryRow(`select id from forum where slug=$1`, slug).Scan(&forumId)

	if err != nil {
		return nil, 404
	}

	var sorting string
	var compareString string
	if desc {
		sorting = "desc"
		compareString = "<"
	} else {
		sorting = "asc"
		compareString = ">"
	}

	var limitString string
	if limit != 0 {
		limitString = fmt.Sprintf("limit %d", limit)
	} else {
		limitString = ""
	}

	var rows *pgx.Rows
	if sinceNickname != "" {
		queryString := fmt.Sprintf(`
			select about, email, fullname, nickname from ForumUser
			where slug=$1 and nickname %s $2
			order by lower(nickname) %s
			%s
		`, compareString, sorting, limitString)

		rows, err = conn.Query(queryString, slug, sinceNickname)
		defer rows.Close()
	} else {
		queryString := fmt.Sprintf(`
			select about, email, fullname, nickname from ForumUser
			where slug=$1
			order by lower(nickname) %s
			%s
		`, sorting, limitString)

		rows, err = conn.Query(queryString, slug)
		defer rows.Close()
	}

	if err != nil {
		return nil, 404
	}

	users = make([]models.User, 0)
	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

		users = append(users, user)
	}

	return users, 200
}
