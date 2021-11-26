package db

import (
	"context"
	"log"
)

type User struct {
	Name       string
	Email      string
	AvatarLink string
	Sex        int
	Bdate      string
	Unique     string
}

type PublicUser struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	AvatarLink string `json:"avatar_link"`
	Sex        int    `json:"sex"`
	Bdate      string `json:"bdate"`
}

func (db *PgxCon) AddUser(u User) (string, error) {
	connCtx, cancel := context.WithTimeout(context.Background(), waitTimeout)
	defer cancel()
	var id string

	_ = db.pgConn.QueryRow(connCtx, "SELECT id from users WHERE unique_identificator=$1", u.Unique).Scan(&id)
	if id != "" {
		return id, nil
	}

	tx, _ := db.pgConn.Begin(connCtx)
	err := tx.QueryRow(connCtx,
		"INSERT INTO Users (user_name,email,avatar_link,sex,bdate,unique_identificator) VALUES ($1,$2,$3,$4,$5,$6) returning id",
		u.Name, u.Email, u.AvatarLink, u.Sex, u.Bdate, u.Unique).Scan(&id)

	if err != nil {
		tx.Rollback(connCtx)
		return "", err
	}

	tx.Commit(connCtx)
	return id, nil
}

func (db *PgxCon) GetUser(uuid string) (*PublicUser, error) {
	var user PublicUser
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()
	err := db.pgConn.QueryRow(connCtx, "SELECT user_name,email,avatar_link,sex,bdate FROM users WHERE id=$1", uuid).
		Scan(&user.Name, &user.Email, &user.AvatarLink, &user.Sex, &user.Bdate)
	log.Println(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *PgxCon) UpdateUserUserName(userName string, uuid string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE users SET user_name=$1 WHERE id=$2",
		userName,
		uuid,
	)

	// Rollback transaction if any error happened on insertion
	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}

func (db *PgxCon) UpdateUserEmail(email string, uuid string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE users SET email=$1 WHERE id=$2",
		email,
		uuid,
	)

	// Rollback transaction if any error happened on insertion
	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}

func (db *PgxCon) UpdateUserAvatar(avatar string, uuid string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE users SET avatar_link=$1 WHERE id=$2",
		avatar,
		uuid,
	)

	// Rollback transaction if any error happened on insertion
	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}

func (db *PgxCon) UpdateUserSex(sex int, uuid string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE users SET sex=$1 WHERE id=$2",
		sex,
		uuid,
	)

	// Rollback transaction if any error happened on insertion
	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}

func (db *PgxCon) UpdateUserBdate(bdate string, uuid string) error {
	connCtx, cancel := context.WithTimeout(db.pgConnCtx, waitTimeout)
	defer cancel()

	tx, _ := db.pgConn.Begin(connCtx)
	_, err := tx.Exec(connCtx, "UPDATE users SET bdate=$1 WHERE id=$2",
		bdate,
		uuid,
	)

	// Rollback transaction if any error happened on insertion
	if err != nil {
		tx.Rollback(connCtx)
		return err
	}

	tx.Commit(connCtx)
	return nil
}
