package data

import (
	"fmt"
	"time"
)

type Thread struct {
	Id			int
	Uuid 		string
	Topic 		string
	UserId		int
	NumReplies	int
	CreatedAt	time.Time
}

type Post struct {
	Id			int
	Uuid		string
	Body		string
	UserId		int
	ThreadId	int
	CreatedAt	time.Time
}

func (th *Thread) Posts() (posts []*Post, err error) {
	statement := `SELECT id, uuid, body, user_id, thread_id, created_at 
FROM posts where thread_id=$1`
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(th.Id)
	if err != nil {
		fmt.Println("1,err:", err)
		return posts, err
	}
	defer rows.Close()

	for rows.Next() {
		post := &Post{}
		err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
		if err != nil {
		fmt.Println("2,err:", err)
			return
		}
		posts = append(posts, post)
	}
		fmt.Println("3,err:", err)
	return
}

func (th *Thread) Date() string {
	return fmt.Sprintf("%d.%d.%d %d:%2d", th.CreatedAt.Year(), th.CreatedAt.Month(), th.CreatedAt.Day(), 
	th.CreatedAt.Hour(), th.CreatedAt.Minute())
}

func (th *Thread) UpdateNumReplies(numReplies int) (err error) {
	statement := "UPDATE threads SET num_replies=$1 WHERE id=$2"
	
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(numReplies, th.Id)
	return	
}

func (th *Thread) User() *User {
	user, err := UserByID(th.UserId)
	if err != nil {
		return nil
	}
	return user
}

func (pt *Post) User() *User {
	user, err := UserByID(pt.UserId)
	if err != nil {
		return nil
	}
	return user
}

func GetAllThreads(num int) (threads []*Thread, err error) {
	statement := `SELECT id, uuid, topic, user_id, num_replies, created_at 
FROM threads ORDER BY created_at DESC`
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return threads, err
	}
	defer rows.Close()

	for rows.Next() {
		thread := &Thread{}
		err = rows.Scan(&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, &thread.NumReplies, &thread.CreatedAt)
		if err != nil {
			return
		}
		threads = append(threads, thread)
	}
	return
}

func ThreadByUUID(uuid string) (thread *Thread, err error) {
	statement := `SELECT id, uuid, topic, user_id, num_replies, created_at 
	FROM threads WHERE uuid=$1`
	
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	thread = &Thread{}
	err = stmt.QueryRow(uuid).Scan(
		&thread.Id, &thread.Uuid, &thread.Topic, &thread.UserId, 
		&thread.NumReplies, &thread.CreatedAt)
	return
}