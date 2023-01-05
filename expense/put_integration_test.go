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

func TestIntegrationUpdateExpenseById(t *testing.T) {
	ec := echo.New()
	serverPort := 3001
	connString := "postgresql://root:root@db/expensedb?sslmode=disable"

	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", connString)
		if err != nil {
			log.Fatal(err)
		}
		h := NewHandler(db)

		e.PUT("/expenses/:id", h.UpdateExpenseById)
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
	expenseId := 2
	reqBody := `{"title":"buy a new phone","amount":32000,"note":"discounted","tags":["gadget","shopping"]}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, expenseId), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Equal(t, "", strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = ec.Shutdown(ctx)
	t.Logf("Port:%d is shut down", serverPort)
	assert.NoError(t, err)
}
