package expense

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *Handler) CreateExpensesHandler(c echo.Context) error {
	e := new(Expense)
	err := c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Request is invalid" + err.Error()})
	}

	row := h.db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))
	err = row.Scan(&e.Id)

	log.Print(e.Amount)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error creating expenses on DB: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
