//go:build integration

package expense

import (
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

var serverPort = 3001

func TestIntegrationCreateExpenses(t *testing.T) {
	ec := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/expensedb?sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		h := NewHandler(db)

		e.POST("/expenses", h.CreateExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
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
	reqBody := `{"title":"buy a new phone","amount":39000,"note":"buy a new phone","tags":["gadget","shopping"]}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	expectedRes := `{"id":2,"title":"buy a new phone","amount":39000,"note":"buy a new phone","tags":["gadget","shopping"]}`

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.Equal(t, expectedRes, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = ec.Shutdown(ctx)
	assert.NoError(t, err)
}
