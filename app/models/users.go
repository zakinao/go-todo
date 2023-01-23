package models

import (
	"log"
	"time"
)

type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
	Todos     []Todo
}

type Session struct {
	ID        int
	UUID      string
	UserID    int
	Email     string
	CreatedAt time.Time
}

func (u *User) CreateUser() (err error) {
	cmd := `insert into users (
                   uuid,
                   name,
                   email,
                   password,
                   created_at) values (?, ?, ?, ?, ?)`

	_, err = Db.Exec(cmd, createUUID(), u.Name, u.Email, Encrypt(u.PassWord), time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUser(id int) (user User, err error) {
	user = User{}
	cmd := `select id, uuid, name, email, password, created_at from users where id = ?`
	err = Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt,
	)
	return user, err
}

func (u *User) UpdateUser() (err error) {
	cmd := `update users set name = ?, email = ? where id = ?`
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (u *User) DeleteUser() (err error) {
	cmd := `delete from users where id = ?`
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func GetUserByEmail(email string) (user User, err error) {
	user = User{}
	cmd := `select id, password, uuid, email, name, created_at from users where email = ?`
	Db.QueryRow(cmd, email).Scan(&user.ID, &user.PassWord, &user.UUID, &user.Email, &user.Name, &user.CreatedAt)
	return user, err
}

func (u *User) CreateSession() (session Session, err error) {
	session = Session{}
	cmd1 := `insert into sessions (uuid, email, user_id, created_at) values (?, ?, ?, ?)`

	_, err = Db.Exec(cmd1, createUUID(), u.Email, u.ID, time.Now())
	if err != nil {
		log.Println(err)
	}

	cmd2 := `select * from sessions where user_id = ? and email = ?`
	err = Db.QueryRow(cmd2, u.ID, u.Email).Scan(&session.ID, &session.UUID, &session.Email, &session.UserID, &session.CreatedAt)

	return session, err
}

func (sess *Session) CheckSession() (valid bool, err error) {
	cmd := `select * from sessions where uuid = ?`

	err = Db.QueryRow(cmd, sess.UUID).Scan(&sess.ID, &sess.UUID, &sess.Email, &sess.UserID, &sess.CreatedAt)

	if err != nil {
		valid = false
		return
	}

	if sess.ID != 0 {
		valid = true
	}
	return valid, err
}

func (sess *Session) DeleteSessionByUUID() (err error) {
	cmd := `delete from sessions where uuid = ?`
	Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

func (sess *Session) GetUserBySession() (user User, err error) {
	user = User{}
	cmd := `select id, uuid, name, email, created_at from users where id = ?`
	Db.QueryRow(cmd, sess.UserID).Scan(&user.ID, &user.UUID, &user.Name, &user.Email, &user.CreatedAt)

	return user, err
}
