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

func TestGetExpenseById(t *testing.T) {
	// ARRANGE
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while opening stub database connection: %s", err)
	}
	defer mockDb.Close()
	query := `SELECT id, title, amount, note, tags FROM expenses WHERE id = $1`
	mockH := NewHandler(mockDb)

	t.Run("should return expected expenses row", func(t *testing.T) {
		// ARRANGE
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")
		expectedRes := `{"id":2,"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]}`

		newMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(2, "apple smoothie", 89, "no discount", pq.Array(&[]string{"beverage"}))
		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs("2").WillReturnRows(newMockRows)

		// ACT
		err = mockH.GetExpenseById(c)

		// ASSERTION
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expectedRes, strings.TrimSpace(rec.Body.String()))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("should return error message if row not found", func(t *testing.T) {
		// ARRANGE
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")
		expectedErr := `{"message":"expense not found"}`

		newMockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"})
		mock.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs("2").WillReturnRows(newMockRows)

		// ACT
		err = mockH.GetExpenseById(c)

		// ASSERTION
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusNotFound, rec.Code)
			assert.Equal(t, expectedErr, strings.TrimSpace(rec.Body.String()))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
