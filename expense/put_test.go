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

func TestUpdateExpenseById(t *testing.T) {
	e := echo.New()

	t.Run("should update a row if existed", func(t *testing.T) {
		expenseJson := `{"title":"buy a new phone","amount":32000,"note":"discounted","tags":["gadget","shopping"]}`
		req := httptest.NewRequest(http.MethodPut, "/expenses/:id", strings.NewReader(expenseJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		// ARRANGE
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("2")

		mockDb, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error while opening stub database connection: %s", err)
		}
		defer mockDb.Close()

		query := `UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id`
		newMockRows := sqlmock.NewRows([]string{"id"}).AddRow(2)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("buy a new phone", 32000.0, "discounted", pq.Array(&[]string{"gadget", "shopping"}), "2").WillReturnRows(newMockRows)
		mockH := NewHandler(mockDb)

		// ACT
		err = mockH.UpdateExpenseById(c)

		// ASSERTION

		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "\"success\"\n", rec.Body.String())

		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("should return error if updating target not existed", func(t *testing.T) {
		// ARRANGE
		expenseJson := `{"title":"buy a new phone","amount":32000,"note":"discounted","tags":["gadget","shopping"]}`
		req := httptest.NewRequest(http.MethodPut, "/expenses/:id", strings.NewReader(expenseJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("10")
		expectedErr := `{"message":"Updating target not found"}`

		mockDb, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Error while opening stub database connection: %s", err)
		}
		defer mockDb.Close()

		query := `UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id`
		newMockRows := sqlmock.NewRows([]string{"id"})
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("buy a new phone", 32000.0, "discounted", pq.Array(&[]string{"gadget", "shopping"}), "10").WillReturnRows(newMockRows)
		mockH := NewHandler(mockDb)

		// ACT
		err = mockH.UpdateExpenseById(c)

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
