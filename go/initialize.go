package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initialize(c *gin.Context) {
	db.Exec("DELETE FROM users WHERE id > 5000")
	db.Exec("DELETE FROM products WHERE id > 10000")
	db.Exec("DELETE FROM comments WHERE id > 200000")
	db.Exec("DELETE FROM histories WHERE id > 500000")

	loading()

	c.String(http.StatusOK, "Finish")
}

func loading() {
	loadComments()
	loadUsers()
	loadProducts()
	loadHistories()
}
