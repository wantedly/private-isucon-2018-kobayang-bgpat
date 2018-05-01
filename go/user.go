package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/parnurzeal/gorequest"
)

var (
	users map[int]*User
)

// User model
type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	LastLogin string
}

func authenticate(email string, password string) (User, bool) {
	var u User
	err := db.QueryRow("SELECT * FROM users WHERE email = ? LIMIT 1", email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		return u, false
	}
	result := password == u.Password
	return u, result
}

func notAuthenticated(session sessions.Session) bool {
	uid := session.Get("uid")
	return !(uid.(int) > 0)
}

func getUser(uid int) User {
	u := User{}
	r := db.QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", uid)
	err := r.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		return u
	}

	return u
}

func currentUser(session sessions.Session) User {
	uid := session.Get("uid")
	u := User{}
	r := db.QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", uid)
	err := r.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
	if err != nil {
		return u
	}

	return u
}

// BuyingHistory : products which user had bought
func (u *User) BuyingHistory() []Product {
	result := make([]Product, 0)
	for _, h := range userHistories[u.ID] {
		var p Product
		if product, ok := products[h.ProductID]; ok {
			p = *product
		}
		tmp, _ := time.Parse("2006-01-02 15:04:05", h.CreatedAt)
		p.CreatedAt = (tmp.Add(9 * time.Hour)).Format("2006-01-02 15:04:05")
		result = append([]Product{p}, result...)
	}

	return result
}

// BuyProduct : buy product
func (u *User) BuyProduct(pid string) {
	productID, _ := strconv.Atoi(pid)
	createdAt := time.Now()

	res, _ := db.Exec(
		"INSERT INTO histories (product_id, user_id, created_at) VALUES (?, ?, ?)",
		pid, u.ID, createdAt)

	id, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	h := History{
		ID:        int(id),
		ProductID: productID,
		UserID:    u.ID,
		CreatedAt: createdAt.Format("2006-01-02 15:04:05"),
	}
	setHistory(h)

	go func() {
		for _, s := range os.Args[1:] {
			go gorequest.New().Post(s + "/push/histories").Send(h).End()
		}
	}()
}

// CreateComment : create comment to the product
func (u *User) CreateComment(pid string, content string) {
	productID, _ := strconv.Atoi(pid)
	createdAt := time.Now()

	res, _ := db.Exec(
		"INSERT INTO comments (product_id, user_id, content, created_at) VALUES (?, ?, ?, ?)",
		pid, u.ID, content, createdAt)

	id, err := res.LastInsertId()
	if err != nil {
		panic(err.Error())
	}
	c := Comment{
		ID:        int(id),
		ProductID: productID,
		UserID:    u.ID,
		Content:   content,
		CreatedAt: createdAt.Format("2006-01-02 15:04:05"),
	}
	setComment(c)

	go func() {
		for _, s := range os.Args[1:] {
			go gorequest.New().Post(s + "/push/comments").Send(c).End()
		}
	}()
}

func (u *User) UpdateLastLogin() {
	lastLogin := time.Now().Format("2006-01-02 15:04:05")
	db.Exec("UPDATE users SET last_login = ? WHERE id = ?", lastLogin, u.ID)
	users[u.ID].LastLogin = lastLogin

	go func() {
		for _, s := range os.Args[1:] {
			go gorequest.New().Post(fmt.Sprintf("%s/push/comments/%d/%s", s, u.ID, lastLogin)).End()
		}
	}()
}

func setUser(u User) {
	users[u.ID] = &u
}

func loadUsers() {
	users = make(map[int]*User)

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Fprintf(os.Stderr, "getUsers %v\n", err)
		return
	}
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.LastLogin)
		setUser(u)
	}
}
