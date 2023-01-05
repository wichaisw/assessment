package expense

import (
	"context"
	"database/sql"
	"encoding/json"
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

func TestIntegrationCreateExpenses(t *testing.T) {
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

		e.POST("/expenses", h.CreateExpenses)
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

	ex := new(Expense)
	if err := json.Unmarshal(byteBody, &ex); err != nil {
		t.Error("Error unmarshalling response")
	}

	initialRows := 2
	if assert.NoError(t, err) {
		a1 := assert.Equal(t, http.StatusCreated, resp.StatusCode)
		a2 := assert.GreaterOrEqual(t, ex.Id, initialRows)

		if a1 && a2 == false {
			ec.Shutdown(ctx)
		}
	}

	err = ec.Shutdown(ctx)
	t.Logf("Port:%d is shut down", serverPort)
	assert.NoError(t, err)
}
