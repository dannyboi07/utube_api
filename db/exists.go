package db

func ActorExistsByEmail(email string) (bool, error) {
	var (
		exists bool
		err    error
	)
	row := db.QueryRow(dbContext, "SELECT EXISTS(SELECT id FROM actor WHERE email = $1)", email)
	err = row.Scan(&exists)
	// if err != nil {
	// 	utils.Log.Println("Failed to check for existing actor by email, err:", err)
	// }

	return exists, err
}
