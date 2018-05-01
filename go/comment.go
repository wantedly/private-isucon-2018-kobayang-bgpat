package main

import (
	"fmt"
	"os"
)

var (
	comments        map[int]*Comment
	productComments map[int][]*Comment
)

// Comment Model
type Comment struct {
	ID        int
	ProductID int
	UserID    int
	Content   string
	CreatedAt string
	User      *User
}

func getComments(pid int) []Comment {
	pComments := productComments[pid]
	obj := make([]Comment, 0)
	for _, c := range pComments {
		obj = append(obj, *c)
	}
	return obj
}

func setComment(c Comment) {
	comments[c.ID] = &c
	if _, ok := productComments[c.ProductID]; !ok {
		productComments[c.ProductID] = make([]*Comment, 0, 1)
	}
	productComments[c.ProductID] = append(productComments[c.ProductID], &c)
}

func loadComments() {
	comments = make(map[int]*Comment)
	productComments = make(map[int][]*Comment)

	rows, err := db.Query("SELECT * FROM comments")
	if err != nil {
		fmt.Fprintf(os.Stderr, "getComments: %v\n", err)
		return
	}
	for rows.Next() {
		var c Comment
		rows.Scan(&c.ID, &c.ProductID, &c.UserID, &c.Content, &c.CreatedAt)
		setComment(c)
	}
}
