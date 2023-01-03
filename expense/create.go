package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func CreateExpensesHandler(c echo.Context) error {
	var e Expenses
	err := c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Request is invalid"})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))
	err = row.Scan(&e.Id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: "Error creating expenses on DB:"})
	}

	return c.JSON(http.StatusCreated, e)
}
