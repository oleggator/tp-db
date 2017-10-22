package db

import (
	"github.com/oleggator/tp-db/models"
)

func CreateUser(profile models.User) (users []models.User, ok bool) {
	ct, _ := conn.Exec(`
           insert into "User" (about, email, fullname, nickname)
           values ($1, $2::text, $3, $4::text);`,
		profile.About, profile.Email, profile.Fullname, profile.Nickname,
	)

	if ct.RowsAffected() > 0 {
		return nil, true
	}

	users = make([]models.User, 0)
	rows, _ := conn.Query(
		`select about, email, fullname, nickname from "User"
        where lower(nickname)=lower($1::text) or lower(email)=lower($2::text);`,
		profile.Nickname, profile.Email,
	)
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

		users = append(users, user)
	}

	return users, false
}

func GetUser(nickname string) (user models.User, ok bool) {
	user = models.User{}
	err := conn.QueryRow(
		`select about, email, fullname, nickname from "User"
        where lower(nickname)=lower($1::text);`,
		nickname,
	).Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)

	return user, err == nil
}

func UpdateUser(srcUser models.User) (status int) {
	user := models.User{}
	var id int
	err := conn.QueryRow(
		`select id, about, email, fullname from "User"
        where lower(nickname)=lower($1::text);`,
		srcUser.Nickname,
	).Scan(&id, &user.About, &user.Email, &user.Fullname)

	if err != nil {
		return 404
	}

	if user.About != srcUser.About && user.About != "" {
		ct, _ := conn.Exec(
			`update "User"
			set about=$1
			where id=$2;`,
			srcUser.About, id,
		)

		if ct.RowsAffected() == 0 {
			return 409
		}
	}

	if user.Email != srcUser.Email && user.Email != "" {
		ct, _ := conn.Exec(
			`update "User"
			set email=$1::text
			where id=$2;`,
			srcUser.Email, id,
		)

		if ct.RowsAffected() == 0 {
			return 409
		}
	}

	if user.Fullname != srcUser.Fullname && user.Fullname != "" {
		ct, _ := conn.Exec(
			`update "User"
			set fullname=$1
			where id=$2;`,
			srcUser.Fullname, id,
		)

		if ct.RowsAffected() == 0 {
			return 409
		}
	}

	return 200
}
