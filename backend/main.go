package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"

    _ "github.com/mattn/go-sqlite3"
)

type Todo struct {
    ID   int    `json:"id"`
    Task string `json:"task"`
}

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite3", "./todos.db")
    if err != nil {
        log.Fatal(err)
    }

    _, err = db.Exec("CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY AUTOINCREMENT, task TEXT)")
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/api/todos", todosHandler)

    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    switch r.Method {
    case "GET":
        rows, err := db.Query("SELECT id, task FROM todos")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var todos []Todo
        for rows.Next() {
            var t Todo
            rows.Scan(&t.ID, &t.Task)
            todos = append(todos, t)
        }
        json.NewEncoder(w).Encode(todos)

    case "POST":
        var t Todo
        json.NewDecoder(r.Body).Decode(&t)
        res, err := db.Exec("INSERT INTO todos (task) VALUES (?)", t.Task)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        id, _ := res.LastInsertId()
        t.ID = int(id)
        json.NewEncoder(w).Encode(t)

    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
