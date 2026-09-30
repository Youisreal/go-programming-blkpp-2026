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

func main() {
	for {
		fmt.Println("Menu:")
		fmt.Println("1. Show Users")
		fmt.Println("2. Exit")

		choice := InputInt("Pilih menu: ")
		switch choice {
		case 1:
			showUsers()
		case 2:
			return
		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}
