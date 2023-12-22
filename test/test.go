package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Alperen10/Konzek-App/handler"
	"github.com/Alperen10/Konzek-App/models"
)

// TestCreateTask function tests the CreateTask handler
func TestCreateTask(t *testing.T) {
	// Create a sample task for testing
	task := &models.CreateTask{
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "Pending",
	}

	// Convert task to JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/api/newTask", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.CreateTask(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Decode the response JSON
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response status
	if status, ok := response["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %s", status)
	}

	// Check the response message
	if message, ok := response["message"].(string); !ok || message != "Task has created" {
		t.Errorf("Expected message 'Task has created', got %s", message)
	}
}

// TestGetAllTasks function tests the GetAllTasks handler
func TestGetAllTasks(t *testing.T) {
	// Create a request
	req, err := http.NewRequest("GET", "/api/task", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.GetAllTasks(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Decode the response JSON
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response status
	if status, ok := response["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %s", status)
	}

	// Check the response message
	if message, ok := response["message"].(string); !ok || message != "All tasks reading" {
		t.Errorf("Expected message 'All tasks reading', got %s", message)
	}
}

// TestGetSingleTask function tests the GetSingleTask handler
func TestGetSingleTask(t *testing.T) {
	taskID := 1

	// Create a request
	req, err := http.NewRequest("GET", "/api/singleTask/"+strconv.Itoa(taskID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.GetSingleTask(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Decode the response JSON
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response status
	if status, ok := response["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %s", status)
	}

	// Check the response message
	if message, ok := response["message"].(string); !ok || message != "Getting a task" {
		t.Errorf("Expected message 'Getting a task', got %s", message)
	}

}

// TestUpdateTask function tests the UpdateTask handler
func TestUpdateTask(t *testing.T) {
	taskID := 1
	newTitle := "Updated Task Title"

	// Create a request
	req, err := http.NewRequest("PUT", "/api/updateTask/"+strconv.Itoa(taskID)+"/"+newTitle, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.UpdateTask(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Decode the response JSON
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response status
	if status, ok := response["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %s", status)
	}

	// Check the response message
	if message, ok := response["message"].(string); !ok || message != "Task updated" {
		t.Errorf("Expected message 'Task updated', got %s", message)
	}
}

// TestDeleteTaskByID function tests the DeleteTaskByID handler
func TestDeleteTaskByID(t *testing.T) {
	taskID := 1

	// Create a request
	req, err := http.NewRequest("DELETE", "/api/deleteTask/"+strconv.Itoa(taskID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	handler.DeleteTaskByID(rr, req)

	// Check the status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Decode the response JSON
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// Check the response status
	if status, ok := response["status"].(string); !ok || status != "success" {
		t.Errorf("Expected status 'success', got %s", status)
	}

	// Check the response message
	if message, ok := response["message"].(string); !ok || message != "Task deleted" {
		t.Errorf("Expected message 'Task deleted', got %s", message)
	}
}
