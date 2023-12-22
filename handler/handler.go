package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Alperen10/Konzek-App/database"
	"github.com/Alperen10/Konzek-App/models"
)

var (
	dbMutex sync.Mutex
	logger  *log.Logger
)

// Initialize the logger in the init function
func InitLog() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}

	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}

// CreateTask 	godoc
// @Summary 	Create tasks
// @Description Save tasks data in Db.
// @Accept		application/json
// @Param 		task body models.CreateTask true "Task"
// @Produce 	application/json
// @Success 	200 {object} map[string]interface{} "success"
// @Router 		/newTask [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dbMutex.Lock()
	defer dbMutex.Unlock()

	db := database.Connection
	task := new(models.Task)

	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		logger.Printf("Error decoding request payload: %s", err)
		http.Error(w, fmt.Sprintf(`{"status": "error", "message": "Invalid request payload", "data": %s}`, err), http.StatusBadRequest)
		return
	}

	// Validate task fields
	if task.Title == "" || task.Description == "" || task.Status == "" {
		http.Error(w, `{"status": "error", "message": "There are empty spaces. Make sure the title, description and status fields are filled in", "data": null}`, http.StatusBadRequest)
		return
	}

	insertTask := `insert into tasks (title, description, status) values($1, $2, $3) RETURNING id`
	err := db.QueryRow(context.Background(), insertTask, task.Title, task.Description, task.Status).Scan(&task.Id)
	if err != nil {
		logger.Printf("Error inserting task into the database: %s", err)
		http.Error(w, fmt.Sprintf(`{"status": "error", "message": "Could not create task", "data": %s}`, err), http.StatusInternalServerError)
		return
	}

	logger.Printf("Task created successfully. ID: %d, Title: %s", task.Id, task.Title)

	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "Task has created", "data": task})
}

// GetAllTasks 		godoc
// @Summary			Get All tasks.
// @Description		Return list of tasks.
// @Tags			tasks
// @Success			200
// @Router			/task [get]
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Validate request method
	if r.Method != http.MethodGet {
		http.Error(w, `{"status": "error", "message": "Invalid request method", "data": null}`, http.StatusMethodNotAllowed)
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	db := database.Connection
	tasks := []models.Task{}
	var tasksMutex sync.Mutex // Mutex to protect concurrent access to tasks slice

	selectAllTasks := `select * from "tasks"`
	rows, err := db.Query(context.Background(), selectAllTasks)

	if err != nil {
		logger.Printf("Error querying tasks from the database: %s", err)
		http.Error(w, fmt.Sprintf(`{"status": "error", "message": "Something's wrong with your input", "data": %s}`, err), http.StatusInternalServerError)
		return
	}

	// Fetch all rows into a slice
	var allRows []models.Task
	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"status": "error", "message": "Could not read task", "data": %s}`, err), http.StatusInternalServerError)
			return
		}
		allRows = append(allRows, task)
	}

	var wg sync.WaitGroup

	for _, currentTask := range allRows {
		wg.Add(1)
		go func(task models.Task) {
			defer wg.Done()
			fmt.Printf("Id: %d - %s - %s - %s\n", task.Id, task.Title, task.Description, task.Status)

			// Protect access to tasks slice with a mutex
			tasksMutex.Lock()
			defer tasksMutex.Unlock()
			tasks = append(tasks, task)
		}(currentTask)
	}

	wg.Wait()

	// Create new instances of tasks to avoid ID duplication
	var newTasks []models.Task
	for _, t := range tasks {
		newTasks = append(newTasks, models.Task{Id: t.Id, Title: t.Title, Description: t.Description, Status: t.Status})
	}

	logger.Printf("All tasks retrieved successfully. Count: %d", len(newTasks))

	// Return tasks
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "All tasks reading", "data": newTasks})
}

// GetSingleTask 	godoc
// @Summary 		Get Single tasks by id.
// @Description 	The requested data is fetched according to the given id number.
// @Produce 		application/json
// @Tags 			tasks
// @Param 			id path int true "Task ID"
// @Success 		200 {object} map[string]interface{} "success"
// @Router 			/singleTask/{id} [get]
func GetSingleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract task ID from the URL path
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) != 4 {
		http.Error(w, `{"status": "error", "message": "Invalid URL path", "data": null}`, http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(pathSegments[3])
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid task ID", "data": null}`, http.StatusBadRequest)
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	db := database.Connection
	task := new(models.Task)

	// Validate task ID
	if taskID <= 0 {
		http.Error(w, `{"status": "error", "message": "Invalid task ID", "data": null}`, http.StatusBadRequest)
		return
	}

	selectTask := `SELECT id, title, description, status FROM tasks WHERE id = $1`
	err = db.QueryRow(context.Background(), selectTask, taskID).Scan(&task.Id, &task.Title, &task.Description, &task.Status)

	if err != nil {
		logger.Printf("Failed to get task: %s", err)
		http.Error(w, fmt.Sprintf(`{"status": "error", "message": "Failed to get task", "data": null}`), http.StatusNotFound)
		return
	}

	logger.Printf("Getting a task. ID: %d, Title: %s", task.Id, task.Title)

	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "Getting a task", "data": task})
}

// UpdateTask godoc
// @Summary Update task by id
// @Description Data is updated according to the given id.
// @Param id path int true "Task ID"
// @Param title path string true "Task Title"
// @Tags tasks
// @Produce application/json
// @Success 200 {object} map[string]interface{} "success"
// @Router /updateTask/{id}/{title} [put]
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract task ID and title from the URL path
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) != 5 {
		http.Error(w, `{"status": "error", "message": "Invalid URL path", "data": null}`, http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(pathSegments[3])
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid task ID", "data": null}`, http.StatusBadRequest)
		return
	}

	taskTitle := pathSegments[4]
	if taskTitle == "" {
		http.Error(w, `{"status": "error", "message": "Title cannot be empty", "data": null}`, http.StatusBadRequest)
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	db := database.Connection

	// Validate task ID
	if taskID <= 0 {
		http.Error(w, `{"status": "error", "message": "Invalid task ID", "data": null}`, http.StatusBadRequest)
		return
	}

	updateTasks := `UPDATE tasks SET title = $1 WHERE id = $2`
	res, err := db.Exec(context.Background(), updateTasks, taskTitle, taskID)

	if err != nil {
		logger.Printf("Failed to update task: %s", err)
		http.Error(w, fmt.Sprintf(`{"status": "error", "message": "Failed to update task", "data": %s}`, err), http.StatusInternalServerError)
		return
	}

	count := res.RowsAffected()
	if count == 0 {
		http.Error(w, `{"status": "error", "message": "Failed to update task", "data": null}`, http.StatusNotFound)
		return
	}

	logger.Printf("Task updated successfully. ID: %d, New Title: %s", taskID, taskTitle)

	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "Task updated"})
}

// DeleteTaskByID godoc
// @Summary Delete task by id
// @Description Data is deleted according to the given id.
// @Produce application/json
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]interface{} "success"
// @Router /deleteTask/{id} [delete]
func DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract task ID from the URL path
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) != 4 {
		http.Error(w, `{"status": "error", "message": "Invalid URL path", "data": null}`, http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(pathSegments[3])
	if err != nil {
		http.Error(w, `{"status": "error", "message": "Invalid task ID", "data": null}`, http.StatusBadRequest)
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	db := database.Connection

	// Validate task ID
	if taskID <= 0 {
		http.Error(w, `{"status": "error", "message": "Invalid task ID", "data": null}`, http.StatusBadRequest)
		return
	}

	deleteTasks := `DELETE FROM tasks WHERE id = $1`
	res, err := db.Exec(context.Background(), deleteTasks, taskID)

	if err != nil {
		logger.Printf("Failed to delete task: %s", err)
		http.Error(w, fmt.Sprintf(`{"status": "error", "message": "Failed to delete task", "data": %s}`, err), http.StatusInternalServerError)
		return
	}

	count := res.RowsAffected()
	if count == 0 {
		http.Error(w, `{"status": "error", "message": "Failed to delete task", "data": null}`, http.StatusNotFound)
		return
	}

	logger.Printf("Task deleted successfully. ID: %d", taskID)

	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "message": "Task deleted"})
}
