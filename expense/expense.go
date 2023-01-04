package expense

import "database/sql"

type (
	Expense struct {
		Id     int      `json:"id"`
		Title  string   `json:"title"`
		Amount float64  `json:"amount"`
		Note   string   `json:"note"`
		Tags   []string `json:"tags"`
	}
	Handler struct {
		// db map[string]*Expense
		db *sql.DB
	}
	Err struct {
		Message string `json:"message"`
	}
)

func InjectHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}
