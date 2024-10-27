package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Mgeorg1/go-todo-list/db"
	sqlc "github.com/Mgeorg1/go-todo-list/db/sqlc"
	"html/template"
	"net/http"
	"strconv"
)

type Server struct {
	todoTemplate *template.Template
	store        *db.Store
}

func CreateServer(store *db.Store) (*Server, error) {
	tmpl, err := template.New("base.html").Funcs(template.FuncMap{
		"sub": sub,
		"add": add,
	}).ParseFiles("web/templates/base.html", "web/templates/task_view.html")
	if err != nil {
		return nil, err
	}
	return &Server{todoTemplate: tmpl, store: store}, nil
}

const pageSize = 5

func sub(a, b int) int {
	return a - b
}

func add(a, b int) int {
	return a + b
}

func (server *Server) Search(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("title")

	tasks, err := server.store.Search(context.Background(), searchText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var selectedTask sqlc.Task
	var selectedTaskPtr *sqlc.Task

	selectedTaskPtr = nil
	taskID := r.URL.Query().Get("taskID")
	if taskID != "" {
		id, err := strconv.Atoi(taskID)
		if err == nil {
			selectedTask, _ = server.store.GetTask(context.Background(), int64(id))
			selectedTaskPtr = &selectedTask
		}
	}

	data := struct {
		Tasks        []sqlc.SearchRow
		SelectedTask *sqlc.Task
		TaskSelected bool
		CurrentPage  int
	}{
		Tasks:        tasks,
		SelectedTask: selectedTaskPtr,
		TaskSelected: &selectedTaskPtr != nil,
		CurrentPage:  1,
	}

	err = server.todoTemplate.Execute(w, data)
	if err != nil {
		err = fmt.Errorf("template error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *Server) ListTasks(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			http.Error(w, "Incorrect page", http.StatusBadRequest)
			return
		}
	}

	offset := int32((page - 1) * pageSize)

	tasks, err := server.store.GetTaskTitles(context.Background(), sqlc.GetTaskTitlesParams{Limit: pageSize, Offset: offset})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var selectedTask sqlc.Task
	var selectedTaskPtr *sqlc.Task

	selectedTaskPtr = nil
	taskID := r.URL.Query().Get("taskID")
	if taskID != "" {
		id, err := strconv.Atoi(taskID)
		if err == nil {
			selectedTask, _ = server.store.GetTask(context.Background(), int64(id))
			selectedTaskPtr = &selectedTask
		}
	}

	data := struct {
		Tasks        []sqlc.GetTaskTitlesRow
		SelectedTask *sqlc.Task
		TaskSelected bool
		CurrentPage  int
	}{
		Tasks:        tasks,
		SelectedTask: selectedTaskPtr,
		TaskSelected: &selectedTaskPtr != nil,
		CurrentPage:  page,
	}

	err = server.todoTemplate.Execute(w, data)
	if err != nil {
		err = fmt.Errorf("template error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (server *Server) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.URL.Path[len("/delete/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		err = server.store.DeleteTask(context.Background(), int64(id))
		if err != nil {
			http.Error(w, "Can not delete task", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (server *Server) SetDone(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.URL.Path[len("/done/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		err = server.store.SetDone(context.Background(), sqlc.SetDoneParams{
			Done: sql.NullBool{Bool: true, Valid: true},
			ID:   int64(id),
		})
		if err != nil {
			http.Error(w, "Can not update task", http.StatusInternalServerError)
			return
		}
		url := "/?taskID=" + idStr
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (server *Server) AddTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		if len(title) == 0 {
			http.Error(w, "Title can not be empty", http.StatusBadRequest)
			return
		}
		text := r.FormValue("text")

		_, err := server.store.CreateTask(context.Background(),
			sqlc.CreateTaskParams{ //tmp
				Title: title,
				Text: sql.NullString{
					Valid:  true,
					String: text,
				},
			})
		if err != nil {
			http.Error(w, "Can not create task", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
