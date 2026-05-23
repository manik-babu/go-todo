package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var db *pgx.Conn

type Todo struct {
	Id          int    `json:"id"`
	Message     string `json:"message"`
	IsCompleted bool   `json:"is_completed"`
}

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		panic(envErr)
	}
	connectDb()
	defer db.Close(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", rootController)
	mux.HandleFunc("POST /todos", createTodo)
	mux.HandleFunc("GET /todos", getTodos)
	mux.HandleFunc("PATCH /todos/{id}", updateTodo)
	mux.HandleFunc("DELETE /todos/{id}", deleteTodo)

	fmt.Println("Server is running on http://localhost:3000")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		panic(err)
	}
}
func connectDb() {
	connectStr := os.Getenv("DB_CONNECT")

	conn, err := pgx.Connect(context.Background(), connectStr)
	if err != nil {
		panic(err)
	}
	db = conn
	fmt.Println("Database connected successfully")
}
