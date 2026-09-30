package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var jsonFile = "data/students.json"

var reader = bufio.NewReader(os.Stdin)

type Student struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Profile Profile  `json:"profile"`
	Courses []Course `json:"courses"`
}

type Profile struct {
	City  string `json:"city"`
	Phone string `json:"phone"`
}

type Course struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Grade string `json:"grade"`
}

func main() {
	for {
		fmt.Println("\nMenu Data Mahasiswa")
		fmt.Println("1. Tampilkan data")
		fmt.Println("2. Tambah data")
		fmt.Println("3. Ubah data")
		fmt.Println("4. Hapus data")
		fmt.Println("5. Cari data")
		fmt.Println("0. Keluar")

		menu := inputInt("Pilih menu: ")

		switch menu {
		case 1:
			showStudents()
		case 2:
			addStudent()
		case 3:
			updateStudent()
		case 4:
			deleteStudent()
		case 5:
			searchStudent()
		case 0:
			return
		default:
			fmt.Println("Menu tidak valid")
		}
	}
}

func showStudents() {
	students := readStudents()
	fmt.Println("ID | Nama | Email | Kota | Mata Kuliah")
	for _, student := range students {
		fmt.Printf("%d | %s | %s | %s | %s\n",
			student.ID,
			student.Name,
			student.Email,
			student.Profile.City,
			formatCourses(student.Courses),
		)
	}
}

func addStudent() {
	students := readStudents()

	name := input("Nama: ")
	email := input("Email: ")
	city := input("Kota: ")
	phone := input("Telepon: ")

	id := 1
	if len(students) > 0 {
		id = students[len(students)-1].ID + 1
	}

	students = append(students, Student{
		ID:    id,
		Name:  name,
		Email: email,
		Profile: Profile{
			City:  city,
			Phone: phone,
		},
		Courses: []Course{},
	})
	writeStudents(students)
	fmt.Println("Data ditambahkan")
}

func updateStudent() {
	students := readStudents()

	id := inputInt("ID yang diubah: ")
	name := input("Nama baru: ")
	email := input("Email baru: ")
	city := input("Kota baru: ")
	phone := input("Telepon baru: ")

	for i, student := range students {
		if student.ID == id {
			students[i].Name = name
			students[i].Email = email
			students[i].Profile = Profile{City: city, Phone: phone}
			writeStudents(students)
			fmt.Println("Data diubah")
			return
		}
	}
	fmt.Println("Data tidak ditemukan")
}

func deleteStudent() {
	students := readStudents()

	id := inputInt("ID yang dihapus: ")

	for i, student := range students {
		if student.ID == id {
			students = append(students[:i], students[i+1:]...)
			writeStudents(students)
			fmt.Println("Data dihapus")
			return
		}
	}
	fmt.Println("Data tidak ditemukan")
}

func searchStudent() {
	students := readStudents()

	keyword := input("Cari nama/email/kota/mata kuliah: ")
	keyword = strings.ToLower(keyword)

	for _, student := range students {
		name := strings.ToLower(student.Name)
		email := strings.ToLower(student.Email)
		city := strings.ToLower(student.Profile.City)
		courses := strings.ToLower(formatCourses(student.Courses))
		if strings.Contains(name, keyword) ||
			strings.Contains(email, keyword) ||
			strings.Contains(city, keyword) ||
			strings.Contains(courses, keyword) {
			fmt.Printf("%d | %s | %s | %s | %s\n",
				student.ID,
				student.Name,
				student.Email,
				student.Profile.City,
				formatCourses(student.Courses),
			)
		}
	}
}

func input(label string) string {
	// Catatan: helper ini dipakai supaya input seperti nama bisa memakai spasi.
	fmt.Print(label)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func inputInt(label string) int {
	value, _ := strconv.Atoi(input(label))
	return value
}

func readStudents() []Student {
	file, err := os.Open(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var students []Student
	if err := json.NewDecoder(file).Decode(&students); err != nil {
		log.Fatal(err)
	}
	return students
}

func writeStudents(students []Student) {
	file, err := os.Create(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(students); err != nil {
		log.Fatal(err)
	}
}

func formatCourses(courses []Course) string {
	var labels []string
	for _, course := range courses {
		labels = append(labels, fmt.Sprintf("%s %s (%s)", course.Code, course.Title, course.Grade))
	}
	return strings.Join(labels, ", ")
}
