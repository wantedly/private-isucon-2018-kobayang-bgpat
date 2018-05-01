package main

import (
	"sync"
)

var (
	histories     map[int]*History
	userHistories map[int][]*History
	historiesMu   sync.Mutex
)

type History struct {
	ID        int
	ProductID int
	UserID    int
	CreatedAt string
}

func setHistory(h History) {
	histories[h.ID] = &h
	if _, ok := userHistories[h.UserID]; !ok {
		userHistories[h.UserID] = make([]*History, 0, 1)
	}
	historiesMu.Lock()
	userHistories[h.UserID] = append(userHistories[h.UserID], &h)
	historiesMu.Unlock()
}

func loadHistories() {
	histories = make(map[int]*History)
	userHistories = make(map[int][]*History)

	rows, err := db.Query("SELECT * FROM histories")
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		var h History
		rows.Scan(&h.ID, &h.ProductID, &h.UserID, &h.CreatedAt)
		setHistory(h)
	}
}
