package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *Handler) GetExpenseById(c echo.Context) error {
	ex := new(Expense)
	expenseId := c.Param("id")
	query := `SELECT id, title, amount, note, tags FROM expenses WHERE id = $1`
	stmt, err := h.db.Prepare(query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Couldn't prepare GetExpenseById statement, parameter might be invalid: " + err.Error()})
	}

	row := stmt.QueryRow(expenseId)
	err = row.Scan(&ex.Id, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))
	defer stmt.Close()

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, ex)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error getting expenses by id on DB: " + err.Error()})
	}
}

func (h *Handler) GetAllExpenses(c echo.Context) error {
	query := `SELECT id, title, amount, note, tags FROM expenses`

	rows, err := h.db.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error querying all expense"})
	}

	expenses := []Expense{}
	for rows.Next() {
		var ex Expense
		err = rows.Scan(&ex.Id, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))
		expenses = append(expenses, ex)
	}

	return c.JSON(http.StatusOK, expenses)
}
