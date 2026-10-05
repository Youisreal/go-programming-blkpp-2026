package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var userFile = "data/userDb.json"
var taskFile = "data/taskDb.json"

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Task struct {
	Id        int          `json:"id"`
	UserId    User         `json:"userId"`
	Title     string       `json:"title"`
	Timeline  Tasktimeline `json:"tasktimeline"`
	Deskripsi string       `json:"deskripsi"`
	Progress  float64      `json:"progress"`
	Status    bool         `json:"status"`
}

type Tasktimeline struct {
	CreatedAt time.Time `json:"waktuBuat"`
	Deadline  time.Time `json:"deadline"`
}

// Request struct khusus untuk penambahan Task dari Postman
type CreateTaskInput struct {
	UserId    int    `json:"userId"`
	Title     string `json:"title"`
	Deadline  string `json:"deadline"` // Format: YYYY-MM-DD
	Deskripsi string `json:"deskripsi"`
}

// Helper untuk membaca dan menulis file JSON
func readUsersFromFile() []User {
	file, err := os.Open(userFile)
	if err != nil {
		return []User{}
	}
	defer file.Close()

	var users []User
	json.NewDecoder(file).Decode(&users)
	return users
}

func writeUsersToFile(users []User) {
	os.MkdirAll("data", os.ModePerm)
	file, err := os.Create(userFile)
	if err != nil {
		log.Println("Gagal menulis file user:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(users)
}

func readTasksFromFile() []Task {
	file, err := os.Open(taskFile)
	if err != nil {
		return []Task{}
	}
	defer file.Close()

	var tasks []Task
	json.NewDecoder(file).Decode(&tasks)
	return tasks
}

func writeTasksToFile(tasks []Task) {
	os.MkdirAll("data", os.ModePerm)
	file, err := os.Create(taskFile)
	if err != nil {
		log.Println("Gagal menulis file task:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(tasks)
}

// --- HANDLER USER ---

func userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	switch r.Method {
	case http.MethodGet:
		users := readUsersFromFile()
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var newUser User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "Payload JSON tidak valid", http.StatusBadRequest)
			return
		}

		users := readUsersFromFile()
		id := 1
		if len(users) > 0 {
			id = users[len(users)-1].Id + 1
		}
		newUser.Id = id

		users = append(users, newUser)
		writeUsersToFile(users)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)

	case http.MethodPut:
		if len(pathParts) < 2 {
			http.Error(w, "ID User diperlukan di URL (contoh: /users/1)", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, "ID User harus berupa angka", http.StatusBadRequest)
			return
		}

		var updatedUser User
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			http.Error(w, "Payload JSON tidak valid", http.StatusBadRequest)
			return
		}

		users := readUsersFromFile()
		found := false
		for i, u := range users {
			if u.Id == id {
				updatedUser.Id = id
				users[i] = updatedUser
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}

		writeUsersToFile(users)
		json.NewEncoder(w).Encode(updatedUser)

	case http.MethodDelete:
		if len(pathParts) < 2 {
			http.Error(w, "ID User diperlukan di URL (contoh: /users/1)", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, "ID User harus berupa angka", http.StatusBadRequest)
			return
		}

		users := readUsersFromFile()
		found := false
		for i, u := range users {
			if u.Id == id {
				users = append(users[:i], users[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "User tidak ditemukan", http.StatusNotFound)
			return
		}

		writeUsersToFile(users)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil dihapus"})

	default:
		http.Error(w, "Method tidak didukung", http.StatusMethodNotAllowed)
	}
}

// --- HANDLER TASK ---

func taskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	switch r.Method {
	case http.MethodGet:
		tasks := readTasksFromFile()
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var input CreateTaskInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Payload JSON tidak valid", http.StatusBadRequest)
			return
		}

		users := readUsersFromFile()
		var foundUser User
		for _, u := range users {
			if u.Id == input.UserId {
				foundUser = u
				break
			}
		}

		if foundUser.Id == 0 {
			http.Error(w, "User ID tidak ditemukan", http.StatusNotFound)
			return
		}

		deadline, err := time.Parse("2006-01-02", input.Deadline)
		if err != nil {
			http.Error(w, "Format deadline harus YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		tasks := readTasksFromFile()
		id := 1
		if len(tasks) > 0 {
			id = tasks[len(tasks)-1].Id + 1
		}

		newTask := Task{
			Id:     id,
			UserId: foundUser,
			Title:  input.Title,
			Timeline: Tasktimeline{
				CreatedAt: time.Now(),
				Deadline:  deadline,
			},
			Deskripsi: input.Deskripsi,
			Progress:  0.0,
			Status:    false,
		}

		tasks = append(tasks, newTask)
		writeTasksToFile(tasks)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	case http.MethodPut:
		if len(pathParts) < 2 {
			http.Error(w, "ID Task diperlukan di URL (contoh: /tasks/1)", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, "ID Task harus berupa angka", http.StatusBadRequest)
			return
		}

		var updateData Task
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, "Payload JSON tidak valid", http.StatusBadRequest)
			return
		}

		tasks := readTasksFromFile()
		found := false
		for i, t := range tasks {
			if t.Id == id {
				if updateData.Title != "" {
					tasks[i].Title = updateData.Title
				}
				if updateData.Deskripsi != "" {
					tasks[i].Deskripsi = updateData.Deskripsi
				}
				if updateData.Progress != 0 {
					tasks[i].Progress = updateData.Progress
				}
				tasks[i].Status = updateData.Status
				if !updateData.Timeline.Deadline.IsZero() {
					tasks[i].Timeline.Deadline = updateData.Timeline.Deadline
				}
				found = true
				updateData = tasks[i]
				break
			}
		}

		if !found {
			http.Error(w, "Task tidak ditemukan", http.StatusNotFound)
			return
		}

		writeTasksToFile(tasks)
		json.NewEncoder(w).Encode(updateData)

	case http.MethodDelete:
		if len(pathParts) < 2 {
			http.Error(w, "ID Task diperlukan di URL (contoh: /tasks/1)", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, "ID Task harus berupa angka", http.StatusBadRequest)
			return
		}

		tasks := readTasksFromFile()
		found := false
		for i, t := range tasks {
			if t.Id == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "Task tidak ditemukan", http.StatusNotFound)
			return
		}

		writeTasksToFile(tasks)
		json.NewEncoder(w).Encode(map[string]string{"message": "Task berhasil dihapus"})

	default:
		http.Error(w, "Method tidak didukung", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/users", userHandler)
	http.HandleFunc("/users/", userHandler)
	http.HandleFunc("/tasks", taskHandler)
	http.HandleFunc("/tasks/", taskHandler)

	fmt.Println("Server berjalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
