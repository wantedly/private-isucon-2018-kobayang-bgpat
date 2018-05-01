package main

import (
	"fmt"
	"os"
	"sync"
)

var (
	products   map[int]*Product
	productsMu sync.Mutex
)

// Product Model
type Product struct {
	ID          int
	Name        string
	Description string
	ImagePath   string
	Price       int
	CreatedAt   string
}

// ProductWithComments Model
type ProductWithComments struct {
	ID           int
	Name         string
	Description  string
	ImagePath    string
	Price        int
	CreatedAt    string
	CommentCount int
	Comments     []CommentWriter
}

// CommentWriter Model
type CommentWriter struct {
	Content string
	Writer  string
}

func getProduct(pid int) Product {
	p := Product{}
	row := db.QueryRow("SELECT * FROM products WHERE id = ? LIMIT 1", pid)
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)
	if err != nil {
		panic(err.Error())
	}

	return p
}

func getProductsWithCommentsAt(page int) []ProductWithComments {
	// select 50 products with offset page*50
	products := []ProductWithComments{}
	rows, err := db.Query("SELECT * FROM products ORDER BY id DESC LIMIT 50 OFFSET ?", page*50)
	if err != nil {
		return nil
	}

	defer rows.Close()
	for rows.Next() {
		p := ProductWithComments{}
		err = rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)

		// select comment count for the product
		cnt := len(productComments[p.ID])
		p.CommentCount = cnt

		if cnt > 0 {
			// select 5 comments and its writer for the product
			var cWriters []CommentWriter
			for i := 0; i < 5; i++ {
				if i >= cnt {
					break
				}
				c := productComments[p.ID][cnt-i-1]
				cWriters = append(cWriters, CommentWriter{
					Content: c.Content,
					Writer:  users[c.UserID].Name,
				})
			}

			p.Comments = cWriters
		}

		products = append(products, p)
	}

	return products
}

func (p *Product) isBought(uid int) bool {
	for _, h := range userHistories[uid] {
		if h.ProductID == p.ID {
			return true
		}
	}

	return false
}

func setProduct(p Product) {
	productsMu.Lock()
	products[p.ID] = &p
	productsMu.Unlock()
}

func loadProducts() {
	products = make(map[int]*Product)

	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		fmt.Fprintf(os.Stderr, "getProducts: %v\n", err)
		return
	}
	for rows.Next() {
		var p Product
		rows.Scan(&p.ID, &p.Name, &p.Description, &p.ImagePath, &p.Price, &p.CreatedAt)
		setProduct(p)
	}
}
