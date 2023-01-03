package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/wichaisw/assessment/expense"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	expense.InitDb()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	go func() {
		log.Println("Server is starting...")
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			log.Fatal("shutting down server")
		}
	}()

	fmt.Println("start at port:", os.Getenv("PORT"))

	expenseRoutes := e.Group("/expenses")
	expenseRoutes.POST("", expense.CreateExpensesHandler)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
