package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func rootController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server running...")
}
func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Ok:      false,
			Message: "Invalid data",
			Data:    nil,
		})
		return
	}
	query := `INSERT INTO todos(message) VALUES($1) RETURNING *`
	db.QueryRow(context.Background(), query, todo.Message).Scan(&todo.Id, &todo.Message, &todo.IsCompleted)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Ok:      true,
		Message: "Todo created successfully",
		Data:    todo,
	})
}
func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []Todo

	query := `SELECT * FROM todos`

	rows, err := db.Query(context.Background(), query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Ok:      false,
			Message: "Invalid data",
			Data:    nil,
		})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Message, &todo.IsCompleted)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Ok:      false,
				Message: "Unable to retrieved data",
				Data:    nil,
			})
			return
		}
		todos = append(todos, todo)
	}
	er := rows.Err()
	if er != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Ok:      false,
			Message: "Unable to retrieved data",
			Data:    nil,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Ok:      true,
		Message: "Todo retrieved successfully",
		Data:    todos,
	})

}
func updateTodo(w http.ResponseWriter, r *http.Request) {
	_id := r.PathValue("id")

	id, err := strconv.Atoi(_id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Ok:      false,
			Message: "Todo update failed! Invalid Id",
			Data:    nil,
		})
		return
	}
	var data Todo
	json.NewDecoder(r.Body).Decode(&data)

	query := `UPDATE todos SET message = $1, is_completed = $2 WHERE id = $3 RETURNING *`

	err = db.QueryRow(context.Background(), query, data.Message, data.IsCompleted, id).Scan(&data.Id, &data.Message, &data.IsCompleted)
	if err != nil {
		fmt.Println("Error:", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Ok:      false,
			Message: "Todo update failed",
			Data:    err,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(Response{
		Ok:      true,
		Message: "User retrieved successfully",
		Data:    data,
	})
}
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Ok:      false,
			Message: "Todo delate failed! Invalid Id",
			Data:    nil,
		})
		return
	}
	findQuery := `SELECT id FROM todos WHERE id = $1`
	var todoId int
	db.QueryRow(context.Background(), findQuery, idNum).Scan(&todoId)
	if todoId == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{
			Ok:      false,
			Message: "Todo not found",
			Data:    nil,
		})
		return
	}

	query := `DELETE FROM todos WHERE id = $1`
	db.QueryRow(context.Background(), query, idNum)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Ok:      true,
		Message: "Todo deleted successfully",
		Data:    nil,
	})

}
