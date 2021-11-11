//go mod init
//go mod tidy

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTask []task //Creating array

var tasks = allTask{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

// +++++++ CREATE A NEW TASK ++++++++
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTasks task
	reqBody, err := ioutil.ReadAll(r.Body) //Information that the customer will send to the server
	if err != nil {
		fmt.Println("Insert a valid task")
	}
	json.Unmarshal(reqBody, &newTasks)

	newTasks.ID = len(tasks) + 1
	tasks = append(tasks, newTasks)

	//Return with the new task created.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTasks)

}

//+++++++ GET ALL TASKS +++++++
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

//+++++++ GET A SINGLE TASK +++++++
func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}

}

//++++++++ UPDATE A TASK ++++++++
func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Please Enter Valid Data")
	}
	json.Unmarshal(reqBody, &updatedTask)

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			updatedTask.ID = taskID
			tasks = append(tasks, updatedTask)
			fmt.Fprintf(w, "The task with the ID %v, has been updated successfully", taskID)
		}
	}

}

//++++++++++ DELETE A TASK +++++++++++++
func deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "Invalid ID")
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with the ID %v has been deleted successfully", taskID)
		}
	}
}

//+++++++++ INDEX +++++++++++
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my rest API")
}

//++++++++ FUNCTION MAIN +++++++++
func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	log.Fatal(http.ListenAndServe(":3000", router))
}
