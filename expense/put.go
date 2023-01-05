package expense

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *Handler) UpdateExpenseById(c echo.Context) error {
	ex := new(Expense)
	expenseId := c.Param("id")
	err := c.Bind(&ex)
	query := `UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id`

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Request is invalid" + err.Error()})
	}

	row := h.db.QueryRow(query, ex.Title, ex.Amount, ex.Note, pq.Array(&ex.Tags), expenseId)
	if err != nil {
		log.Printf("Error updating expenses row id %s: %s", expenseId, err.Error())
	}

	err = row.Scan(&ex.Id)

	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "Updating target not found"})
	case nil:
		return c.JSON(http.StatusNoContent, nil)
	default:
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error updating expenses by id on DB: " + err.Error()})
	}
}
