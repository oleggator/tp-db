package db

func CountForums() (count int32) {
	conn.QueryRow("select count(*) from Forum;").Scan(&count)
	return
}

func CountPosts() (count int64) {
	conn.QueryRow("select count(*) from Post;").Scan(&count)
	return
}

func CountThreads() (count int32) {
	conn.QueryRow("select count(*) from Thread;").Scan(&count)
	return
}

func CountUsers() (count int32) {
	conn.QueryRow(`select count(*) from "User";`).Scan(&count)
	return
}

func Clear() {
	conn.Exec(`TRUNCATE "User", post, thread, vote, forum, forumuser RESTART IDENTITY CASCADE`)
}
