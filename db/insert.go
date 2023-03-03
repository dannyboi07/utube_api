package db

import (
	"utube/models"
)

func InsertActor(insertActor models.Actor) (models.Actor, error) {
	err := db.QueryRow(dbContext, "INSERT INTO actor (email, password, name) VALUES ($1, $2, $3) RETURNING id, email, password, name, created_at, updated_at", insertActor.Email, insertActor.Password, insertActor.Name).Scan(&insertActor.Id, &insertActor.Email, &insertActor.Password, &insertActor.Name, &insertActor.CreatedAt, &insertActor.UpdatedAt)
	// if err != nil {
	// 	utils.Log.Println("Err inserting actor into db, err:", err)
	// 	return insertActor, err
	// }

	return insertActor, err
}
