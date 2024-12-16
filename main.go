package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getTasks(w http.ResponseWriter, j *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func postTask(p http.ResponseWriter, u *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(u.Body)
	if err != nil {
		http.Error(p, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(p, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := tasks[task.ID]
	if ok {
		http.Error(p, "Задача уже существует", http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task
	p.Header().Set("Content-Type", "application/json")
	p.WriteHeader(http.StatusCreated)
}

func getTask(q http.ResponseWriter, t *http.Request) {
	taskID := chi.URLParam(t, "id")
	task, ok := tasks[taskID]
	if !ok {
		http.Error(q, "Задача не найдена", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(q, err.Error(), http.StatusBadRequest)
		return
	}
	q.Header().Set("Content-Type", "application/json")
	q.WriteHeader(http.StatusOK)
	q.Write(resp)
}

func deleteTask(d http.ResponseWriter, e *http.Request) {
	deleteTaskId := chi.URLParam(e, "id")
	_, ok := tasks[deleteTaskId]
	if !ok {
		http.Error(d, "Задача не найдена", http.StatusBadRequest)
		return
	}
	delete(tasks, deleteTaskId)
	d.Header().Set("Content-Type", "application/json")
	d.WriteHeader(http.StatusOK)
	
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)


	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
