package expense

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *Handler) GetExpensesById(c echo.Context) error {
	expenseId := c.Param("id")
	stmt, err := h.db.Prepare("SELECT * FROM expenses WHERE id = $1")
	defer stmt.Close()

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Couldn't prepare GetExpenseById statement, parameter might be invalid: " + err.Error()})
	}

	e := new(Expense)
	row := stmt.QueryRow(expenseId)
	err = row.Scan(&e.Id, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error getting expenses by id on DB: " + err.Error()})
	}
}
