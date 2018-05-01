package main

var (
	histories     map[int]*History
	userHistories map[int][]*History
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
	userHistories[h.UserID] = append(userHistories[h.UserID], &h)
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
