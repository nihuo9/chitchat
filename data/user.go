package data

import (
	"fmt"
	"time"
)

type User struct {
	Id			int
	Uuid		string
	Name		string
	Email		string
	Password	string
	CreatedAt	time.Time
}

type Session struct {
	Id			int
	Uuid		string
	Email		string
	UserId		int
	CreatedAt	time.Time
}

func (user *User) CreateSession() (session *Session, err error) {
	statement := `insert into sessions 
(uuid, email, user_id, created_at) values ($1, $2, $3, $4) 
returning id, uuid, email, user_id, created_at`

	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	session = &Session{}
	err = stmt.QueryRow(CreateUUID(), user.Email, user.Id, time.Now()).Scan(
		&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

func (user *User) Session() (session *Session, err error) {
	statement := `SELECT id, uuid, email, user_id, created_at 
FROM sessions WHERE user_id=$1`

	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	session = &Session{}
	err = stmt.QueryRow(user.Id).Scan(
		&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return	
}

func (user *User) Export() (err error) {
	statement := `insert into users(uuid, name, email, password, created_at)
 values ($1, $2, $3, $4, $5) returning id, created_at`

	stmt, err := DB.Prepare(statement)
	if err != nil {
		fmt.Println("Prepare:", err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.Uuid, user.Name, user.Email, user.Password, time.Now()).Scan(&user.Id, &user.CreatedAt)
	return
}

func (user *User) Delete() (err error) {
	statement := "delete from users where id = $1"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return
}

func (user *User) Update() (err error) {
	statement := "update users set name=$2, email=$3 where id=$1"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Name, user.Email)
	return
}

func (user *User) CreateThread(topic string) (thread *Thread, err error) {
	statement := `insert into threads (uuid, topic, user_id, num_replies, created_at) values 
($1, $2, $3, $4, $5) returning id`
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	thread = &Thread{
		Uuid: CreateUUID(),
		Topic: topic,
		UserId: user.Id,
		NumReplies: 0,
		CreatedAt: time.Now(),
	}
	//_, err = stmt.QueryRow()
	err = stmt.QueryRow(thread.Uuid, thread.Topic, thread.UserId, thread.NumReplies, 
		thread.CreatedAt).Scan(&thread.Id)
	return
}

func (user *User) CreatePost(thread *Thread, body string) (post *Post, err error) {
	statement := `insert into posts (uuid, body, user_id, thread_id, created_at) values 
($1, $2, $3, $4, $5) returning id`
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	post = &Post{
		Uuid: CreateUUID(),
		Body: body,
		UserId: user.Id,
		ThreadId: thread.Id,
		CreatedAt: time.Now(),
	}
	err = stmt.QueryRow(post.Uuid, post.Body, post.UserId, post.ThreadId,
		 post.CreatedAt).Scan(&post.Id)
	return
}

func (sess *Session) CheckValid() (valid bool) {
	valid = true
	if sess.Id == 0 {
		valid = false
	}
	return
}

func (sess *Session) Delete() (err error) {
	statement := "delete from sessions where uuid=$1"
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(sess.Uuid)
	return
}

func (sess *Session) User() (user *User, err error) {
	fmt.Println("userid:", sess.UserId)
	return UserByID(sess.UserId)
}

//
func Users() (users []*User, err error) {
	rows, err := DB.Query("SELECT id, uuid, name, email, password, create_at FROM users")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	return
}

func UserByEmail(email string) (user *User, err error) {
	statement := `SELECT id, uuid, name, email, password, created_at 
FROM users WHERE email=$1`

	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	user = &User{}
	err = stmt.QueryRow(email).Scan(
		&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func UserByUUID(uuid string) (user *User, err error) {
	statement := `SELECT id, uuid, name, email, password, created_at 
FROM users WHERE uuid=$1`

	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	user = &User{}
	err = stmt.QueryRow(uuid).Scan(
		&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func UserByID(id int) (user *User, err error) {
	statement := `SELECT id, uuid, name, email, password, created_at 
FROM users WHERE id=$1`

	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	user = &User{}
	err = stmt.QueryRow(id).Scan(
		&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

func SessionByUUID(uuid string) (session *Session, err error) {
	statement := `SELECT id, uuid, email, user_id, created_at 
	FROM sessions WHERE uuid=$1`
	
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	session = &Session{}
	err = stmt.QueryRow(uuid).Scan(
		&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return	
}

func DeleteAllSessions() (err error) {
	statement := "delete from sessions"
	_, err = DB.Exec(statement)
	return
}

func DeleteAllUser() (err error) {
	statement := "delete from users"
	_, err = DB.Exec(statement)
	return
}