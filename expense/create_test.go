//go:build unit

package expense

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateExpense(t *testing.T) {
	expenseJson := `{"title":"buy a new phone","amount":39000,"note":"buy a new phone","tags":["gadget","shopping"]}`
	expectedRes := `{"id":1,"title":"buy a new phone","amount":39000,"note":"buy a new phone","tags":["gadget","shopping"]}`

	// ARRANGE
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(expenseJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while opening stub database connection: %s", err)
	}
	defer mockDb.Close()

	query := `INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id`
	newMockRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("buy a new phone", 39000.0, "buy a new phone", pq.Array([]string{"gadget", "shopping"})).WillReturnRows(newMockRows)
	mockH := NewHandler(mockDb)

	// ACT
	err = mockH.CreateExpense(c)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	// ASSERTION
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedRes, strings.TrimSpace(rec.Body.String()))
	}
}
