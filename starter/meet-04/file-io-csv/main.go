package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const csvFile = "data/students.csv"

var reader = bufio.NewReader(os.Stdin)

type Student struct {
	ID    int
	Name  string
	Email string
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
	fmt.Println("ID | Nama | Email")
	for _, student := range students {
		fmt.Printf("%d | %s | %s\n", student.ID, student.Name, student.Email)
	}
}

func addStudent() {
	students := readStudents()

	name := input("Nama: ")
	email := input("Email: ")

	id := 1
	if len(students) > 0 {
		id = students[len(students)-1].ID + 1
	}

	students = append(students, Student{id, name, email})
	writeStudents(students)
	fmt.Println("Data ditambahkan")
}

func updateStudent() {
	students := readStudents()

	id := inputInt("ID yang diubah: ")
	name := input("Nama baru: ")
	email := input("Email baru: ")

	for i, student := range students {
		if student.ID == id {
			students[i] = Student{id, name, email}
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

	keyword := input("Cari nama/email: ")
	keyword = strings.ToLower(keyword)

	for _, student := range students {
		name := strings.ToLower(student.Name)
		email := strings.ToLower(student.Email)
		if strings.Contains(name, keyword) || strings.Contains(email, keyword) {
			fmt.Printf("%d | %s | %s\n", student.ID, student.Name, student.Email)
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
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var students []Student
	for _, row := range rows[1:] {
		id, _ := strconv.Atoi(row[0])
		students = append(students, Student{id, row[1], row[2]})
	}
	return students
}

func writeStudents(students []Student) {
	file, err := os.Create(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"id", "name", "email"})
	for _, student := range students {
		writer.Write([]string{
			strconv.Itoa(student.ID),
			student.Name,
			student.Email,
		})
	}
}
