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
	"github.com/wichaisw/assessment/health"

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

	// inject real

	h := expense.NewHandler(expense.GetDb())

	e.GET("/health", health.GetHealthHandler)
	expenseRoutes := e.Group("/expenses")
	expenseRoutes.GET("/:id", h.GetExpenseById)
	expenseRoutes.POST("", h.CreateExpenses)
	expenseRoutes.PUT("/:id", h.UpdateExpenseById)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
