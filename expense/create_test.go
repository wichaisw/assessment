package expense

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	// mockData = map[int]*Expense{
	// 	0: {1, "buy a new phone", 39000, "buy a new phone", []string{"gadget", "shopping"}},
	// 	1: {2, "apple smoothie", 89, "no discount", []string{"beverage"}},
	// 	2: {3, "iPhone 14 Pro Max 1TB", 66900, "birthday gift from my love", []string{"gadget"}},
	// 	3: {4, "strawberry smoothie", 79, "night market promotion discount 10 bath", []string{"food", "beverage"}},
	// }
	expenseJson = `{"title":"buy a new phone","amount":39000,"note":"buy a new phone","tags":["gadget","shopping"]}`
	expectedRes = `{"id":1,"title":"buy a new phone","amount":39000,"note":"buy a new phone","tags":["gadget","shopping"]}`
)

func TestCreateExpensesHandler(t *testing.T) {

	e := echo.New()
	defer e.Close()

	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(expenseJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockDb, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while opening stub database connection: %s", err)
	}
	defer mockDb.Close()
	newMockRows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query := `INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id`
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("buy a new phone", 39000.0, "buy a new phone", pq.Array([]string{"gadget", "shopping"})).WillReturnRows(newMockRows)
	mockH := InjectHandler(mockDb)

	// Assertions
	if assert.NoError(t, mockH.CreateExpensesHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedRes+"\n", rec.Body.String())
	}
}
