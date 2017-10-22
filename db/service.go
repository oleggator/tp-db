package db

func CountForums() (uint, error) {
	var count uint
	err := conn.QueryRow("select count(*) from Forum;").Scan(&count)

	return count, err
}

func CountPosts() (uint, error) {
	var count uint
	err := conn.QueryRow("select count(*) from Post;").Scan(&count)

	return count, err
}

func CountThreads() (uint, error) {
	var count uint
	err := conn.QueryRow("select count(*) from Thread;").Scan(&count)

	return count, err
}

func CountUsers() (uint, error) {
	var count uint
	err := conn.QueryRow("select count(*) from User;").Scan(&count)

	return count, err
}
