package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/LuisanaMTDev/todo/database/gosql_queries"
	"github.com/joho/godotenv"
)

type config struct {
	DbConnection *gosql_queries.Queries
	Platform     string
}

func main() {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	var dbURL string

	if platform == "PROD" {
		dbURL = os.Getenv("DB_URL_PROD")
	} else {
		dbURL = os.Getenv("DB_URL_DEV")
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Printf("Error while opening db: %v", err)
		return
	}

	dbQueries := gosql_queries.New(db)
	handler := http.NewServeMux()
	serverConfig := config{DbConnection: dbQueries, Platform: platform}
	server := http.Server{Handler: handler, Addr: ":8081"}
	log.Printf("Running platfotm: %s", serverConfig.Platform)

	handler.Handle("GET /app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./frontend/assets/"))))

	handler.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

		subjects, err := serverConfig.DbConnection.GetAllSubjects(r.Context())

		if err != nil {
			log.Printf("Error while getting subjects: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		tasksBySubjects := make(map[string][]gosql_queries.Task)

		for _, subject := range subjects {
			subjectsTasks, err := serverConfig.DbConnection.GetAllTasksBySubject(r.Context(), subject.ID)

			if err != nil {
				log.Printf("Error while getting tasks by subjects: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			tasksBySubjects[subject.ID] = subjectsTasks
		}

		data := struct {
			Subjects       []gosql_queries.Subject
			TasksBySubject map[string][]gosql_queries.Task
		}{
			Subjects:       subjects,
			TasksBySubject: tasksBySubjects,
		}

		tmpl := template.Must(template.ParseFiles("./frontend/views/index.html"))
		tmpl.Execute(w, data)
	})

	handler.HandleFunc("GET /subject/add", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./frontend/templates/subjects.go.html"))
		tmpl.ExecuteTemplate(w, "AddSubject", "")
	})

	handler.HandleFunc("POST /subject/add", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		lastSubjectID, err := serverConfig.DbConnection.GetLastSubjectId(r.Context())

		if err != nil {
			log.Printf("Error getting last subject ID: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		newSubjectID, err := increaseID(lastSubjectID)

		if err != nil {
			log.Printf("Error generating new subject ID: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		newSubjectData := gosql_queries.AddSubjectParams{ID: newSubjectID, Name: r.FormValue("name")}

		err = serverConfig.DbConnection.AddSubject(r.Context(), newSubjectData)

		if err != nil {
			log.Printf("Error adding a new subject: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		w.Header().Add("Hx-Redirect", "/")
	})

	handler.HandleFunc("GET /task/add", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./frontend/templates/tasks.go.html"))
		tmpl.ExecuteTemplate(w, "AddTasks", "")
	})

	handler.HandleFunc("POST /task/add", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		lastTaskID, err := serverConfig.DbConnection.GetLastTaskID(r.Context())
		if err != nil {
			log.Printf("Error getting last task ID: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		newTaskID, err := increaseID(lastTaskID)
		if err != nil {
			log.Printf("Error generating new task ID: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		newTaskData := gosql_queries.AddTaskParams{
			ID:          newTaskID,
			Description: r.FormValue("description"),
			ToTimestamp: r.FormValue("due_date"),
			State:       r.FormValue("state"),
			Tags:        strings.Split(r.FormValue("tags"), ", "),
			SubjectID:   r.FormValue("subject_id"),
		}

		err = serverConfig.DbConnection.AddTask(r.Context(), newTaskData)
		if err != nil {
			log.Printf("Error adding a new task: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		w.Header().Add("Hx-Redirect", "/")
	})

	handler.HandleFunc("PUT /task/state", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		_, err = serverConfig.DbConnection.ModifyTaskState(r.Context(), gosql_queries.ModifyTaskStateParams{
			ID:    r.FormValue("task_id"),
			State: r.FormValue("state"),
		})
		if err != nil {
			log.Printf("Error modifing task's state: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	handler.HandleFunc("GET /task/description", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./frontend/templates/tasks.go.html"))
		tmpl.ExecuteTemplate(w, "ChangeTasksDescription", "")
	})

	handler.HandleFunc("PUT /task/description", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		_, err := serverConfig.DbConnection.ModifyTaskDescription(r.Context(), gosql_queries.ModifyTaskDescriptionParams{
			ID:          r.FormValue("task_id"),
			Description: r.FormValue("new_description"),
		})
		if err != nil {
			log.Printf("Error modifing task's description : %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		tmpl := template.Must(template.ParseFiles("./frontend/templates/tasks.go.html"))
		tmpl.ExecuteTemplate(w, "NewDescription", r.FormValue("new_description"))
	})

	handler.HandleFunc("GET /task/due-date", func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("./frontend/templates/tasks.go.html"))
		tmpl.ExecuteTemplate(w, "ChangeDueDate", "")
	})

	handler.HandleFunc("PUT /task/due-date", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		_, err := serverConfig.DbConnection.ModifyTaskDueDate(r.Context(), gosql_queries.ModifyTaskDueDateParams{
			ID:          r.FormValue("task_id"),
			ToTimestamp: r.FormValue("new_due_date"),
		})
		if err != nil {
			log.Printf("Error modifing task's due date: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		task, err := serverConfig.DbConnection.GetTaskByID(r.Context(), r.FormValue("task_id"))
		if err != nil {
			log.Printf("Error while getting task: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		tmpl := template.Must(template.ParseFiles("./frontend/templates/tasks.go.html"))
		tmpl.ExecuteTemplate(w, "NewDueDate", task.DueDate)
	})

	handler.HandleFunc("GET /subject/name", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./frontend/templates/subjects.go.html"))
		tmpl.ExecuteTemplate(w, "ChangeSubjectName", "")
	})

	handler.HandleFunc("PUT /subject/name", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		_, err := serverConfig.DbConnection.ModifySubjectName(r.Context(), gosql_queries.ModifySubjectNameParams{
			ID:   r.FormValue("subject_id"),
			Name: r.FormValue("new_name"),
		})
		if err != nil {
			log.Printf("Error modifing task's due date: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		tmpl := template.Must(template.ParseFiles("./frontend/templates/subjects.go.html"))
		tmpl.ExecuteTemplate(w, "NewName", r.FormValue("new_name"))
	})

	handler.HandleFunc("DELETE /task/delete", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		err = serverConfig.DbConnection.DeleteTask(r.Context(), r.FormValue("task_id"))
		if err != nil {
			log.Printf("Error while deleting task: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Hx-Redirect", "/")
	})

	handler.HandleFunc("DELETE /subject/delete", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		_, err = serverConfig.DbConnection.DeleteSubject(r.Context(), r.FormValue("subject_id"))
		if err != nil {
			log.Printf("Error while deleting subject: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Hx-Redirect", "/")
	})

	handler.HandleFunc("GET /tags/edit", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./frontend/templates/tags.go.html"))
		tmpl.ExecuteTemplate(w, "EditOptions", "")
	})

	handler.HandleFunc("DELETE /tags/delete", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		err = serverConfig.DbConnection.DeleteTag(r.Context(), gosql_queries.DeleteTagParams{
			ID:          r.FormValue("task_id"),
			ArrayRemove: r.FormValue("tag_to_delete"),
		})
		if err != nil {
			log.Printf("Error while deleting tag: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Hx-Redirect", "/")
	})

	handler.HandleFunc("PUT /tags/add", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("Error decoding request body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			// TODO: Return some html indicating the error that occurs and what the user can do next.
			return
		}

		err = serverConfig.DbConnection.AddTag(r.Context(), gosql_queries.AddTagParams{
			ID:          r.FormValue("task_id"),
			ArrayAppend: r.FormValue("tag_to_add"),
		})
		if err != nil {
			log.Printf("Error while adding tag: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Hx-Redirect", "/")
	})

	log.Fatal(server.ListenAndServe())
}
