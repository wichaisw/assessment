//go:build integration

package expense

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationGetExpenseById(t *testing.T) {
	ec := echo.New()
	serverPort := 3001
	connString := "postgresql://root:root@db/expensedb?sslmode=disable"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", connString)
		if err != nil {
			log.Fatal(err)
		}
		h := NewHandler(db)

		e.GET("/expenses/:id", h.GetExpenseById)
		e.Start(fmt.Sprintf(":%d", serverPort))
		defer db.Close()
	}(ec)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	expenseId := 1
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, expenseId), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}
	expectedRes := `{"id":1,"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food","beverage"]}`

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	// Assertion
	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	if assert.NoError(t, err) {
		a1 := assert.Equal(t, http.StatusOK, resp.StatusCode)
		a2 := assert.Equal(t, expectedRes, strings.TrimSpace(string(byteBody)))

		if a1 && a2 == false {
			ec.Shutdown(ctx)
		}
	}

	err = ec.Shutdown(ctx)
	t.Logf("Port:%d is shut down", serverPort)
	assert.NoError(t, err)
}
