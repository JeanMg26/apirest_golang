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
	Id      int    `json:id`
	Name    string `json:name`
	Content string `json:content`
}

type alltaks []task

var tasks = alltaks{
	{
		Id:      1,
		Name:    "task1",
		Content: "content content",
	},
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")

	fmt.Println("iniciando servidor...")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "desde la ruta inicio")
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask task

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Ingresar datos corectos")
		return
	}

	json.Unmarshal(reqBody, &newTask)

	newTask.Id = len(tasks) + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

func GetTask(w http.ResponseWriter, r *http.Request) {

	contenido := mux.Vars(r)

	taskId, err := strconv.Atoi(contenido["id"])
	if err != nil {
		fmt.Fprintf(w, "el id es incorrecto")
	}

	for _, task := range tasks {
		if task.Id == taskId {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

	contenido := mux.Vars(r)

	taskId, err := strconv.Atoi(contenido["id"])
	if err != nil {
		fmt.Fprintf(w, "el id es incorrecto")
	}

	for i, task := range tasks {
		if task.Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "La tarea con el id %v fue eliminada", taskId)
		}
	}
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	contenido := mux.Vars(r)
	var updateTask task

	taskId, err := strconv.Atoi(contenido["id"])
	if err != nil {
		fmt.Fprintf(w, "el id es incorrecto")
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Ingresar datos corectos")
		return
	}

	json.Unmarshal(reqBody, &updateTask)

	for i, task := range tasks {
		if task.Id == taskId {
			tasks = append(tasks[:i], tasks[i+1:]...)

			updateTask.Id = taskId
			tasks = append(tasks, updateTask)
			fmt.Fprintf(w, "La tarea con el id %v fue actualizada", taskId)
		}
	}
}
