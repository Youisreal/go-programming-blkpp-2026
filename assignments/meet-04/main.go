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

const csvFile = "data/userDb.csv"

var reader = bufio.NewReader(os.Stdin)

type User struct {
	ID            int
	Name          string
	Username      string
	Password      string
	Authenticator string
}

func showUsers() {
	users := readUsers()
	fmt.Println("ID | Name | Username | Password | Authenticator")
	for _, user := range users {
		fmt.Printf("%d | %s | %s | %s | %s\n", user.ID, user.Name, user.Username, user.Password, user.Authenticator)
	}
}

func readUsers() []User {
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	for _, row := range rows[1:] {
		id, _ := strconv.Atoi(row[0])
		users = append(users, User{id, row[1], row[2], row[3], row[4]})
	}
	return users
}

func writeUsers(users []User) {
	file, err := os.Create(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"ID", "Name", "Username", "Password", "Authenticator"})

	// Write user data
	for _, user := range users {
		writer.Write([]string{
			strconv.Itoa(user.ID),
			user.Name,
			user.Username,
			user.Password,
			user.Authenticator,
		})
	}
}

func input(label string) string {
	// Catatan: helper ini dipakai supaya input seperti nama bisa memakai spasi.
	fmt.Print(label)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func InputInt(label string) int {
	value, _ := strconv.Atoi(input(label))
	return value
}

func addUser() {
	users := readUsers()

	name := input("Nama: ")
	username := input("Username: ")
	password := input("Password: ")
	authenticator := input("Authenticator: ")

	id := 1
	if len(users) > 0 {
		id = users[len(users)-1].ID + 1
	}

	users = append(users, User{id, name, username, password, authenticator})
	writeUsers(users)
	fmt.Println("User berhasil ditambahkan")
}

func updateUser() {
	users := readUsers()
	showUsers()
	id := InputInt("Masukkan ID user yang ingin diupdate: ")
	name := input("Nama baru: ")
	username := input("Username baru: ")
	password := input("Password baru: ")
	authenticator := input("Authenticator baru: ")

	for i, user := range users {
		if user.ID == id {
			users[i].Name = name
			users[i].Username = username
			users[i].Password = password
			users[i].Authenticator = authenticator
			writeUsers(users)
			fmt.Println("User berhasil diupdate")
			return
		}
	}
	fmt.Println("User dengan ID tersebut tidak ditemukan")
}

func deleteUser() {
	users := readUsers()
	showUsers()
	id := InputInt("Masukkan ID user yang ingin dihapus: ")
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			writeUsers(users)
			fmt.Println("User berhasil dihapus")
			return
		}
	}
	fmt.Println("User dengan ID tersebut tidak ditemukan")
}

func main() {
	for {
		fmt.Println("Menu:")
		fmt.Println("1. Tampilkan Users")
		fmt.Println("2. Tambah User")
		fmt.Println("3. Update User")
		fmt.Println("4. Hapus User")
		fmt.Println("5. Exit")

		choice := InputInt("Pilih menu: ")
		switch choice {
		case 1:
			showUsers()
		case 2:
			addUser()
		case 3:
			updateUser()
		case 4:
			deleteUser()
		case 5:
			fmt.Println("Terima kasih!")
			return
		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}
