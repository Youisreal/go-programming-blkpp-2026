package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Menambahkan struct tags `json:""` agar format JSON rapi dengan huruf kecil
type Student struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Major  string `json:"major"`
	Score  int    `json:"score"`
	Active bool   `json:"active"`
}

const fileName = "students.json"

// Fungsi untuk membaca data dari file JSON
func loadStudents() ([]Student, int) {
	file, err := os.ReadFile(fileName)

	// Jika file belum ada, buat data default awal
	if os.IsNotExist(err) {
		initialData := []Student{
			{ID: 1, Name: "Budi", Major: "Backend", Score: 82, Active: true},
			{ID: 2, Name: "Siti", Major: "Frontend", Score: 91, Active: true},
			{ID: 3, Name: "Andi", Major: "Backend", Score: 74, Active: false},
		}
		saveStudents(initialData)
		return initialData, 4
	} else if err != nil {
		fmt.Println("Gagal membaca file JSON:", err)
		return []Student{}, 1
	}

	var students []Student
	err = json.Unmarshal(file, &students)
	if err != nil {
		fmt.Println("Gagal decode JSON:", err)
		return []Student{}, 1
	}

	// Cari ID tertinggi untuk menentukan nextID berikutnya
	maxID := 0
	for _, s := range students {
		if s.ID > maxID {
			maxID = s.ID
		}
	}

	return students, maxID + 1
}

// Fungsi untuk menyimpan/menulis data ke file JSON
func saveStudents(students []Student) {
	// Menggunakan MarshalIndent agar file JSON mudah dibaca oleh manusia (pretty print)
	jsonData, err := json.MarshalIndent(students, "", "  ")
	if err != nil {
		fmt.Println("Gagal konversi ke JSON:", err)
		return
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		fmt.Println("Gagal menyimpan ke file JSON:", err)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Load data langsung dari file JSON saat aplikasi dibuka
	students, nextID := loadStudents()

	for {
		fmt.Println("\n=== Menu Student (Terhubung dengan JSON) ===")
		fmt.Println("1. Lihat semua student")
		fmt.Println("2. Tambah student")
		fmt.Println("3. Update nilai")
		fmt.Println("4. Hapus student inactive")
		fmt.Println("5. Filter berdasarkan jurusan")
		fmt.Println("6. Cari student berdasarkan ID")
		fmt.Println("7. Lihat rata-rata nilai")
		fmt.Println("0. Keluar")

		choice := readInt(scanner, "Pilih menu: ")
		switch choice {
		case 1:
			printStudents(students)
		case 2:
			name := readString(scanner, "Nama: ")
			major := readString(scanner, "Jurusan: ")
			score := readInt(scanner, "Nilai: ")
			students = append(students, Student{
				ID:     nextID,
				Name:   name,
				Major:  major,
				Score:  score,
				Active: true,
			})
			nextID++

			// Simpan perubahan ke JSON setelah ditambah
			saveStudents(students)
			fmt.Println("Student berhasil ditambahkan.")
		case 3:
			id := readInt(scanner, "ID student: ")
			score := readInt(scanner, "Nilai baru: ")
			if updateScore(students, id, score) {
				// Simpan perubahan ke JSON setelah diupdate
				saveStudents(students)
				fmt.Println("Nilai berhasil diupdate.")
			} else {
				fmt.Println("Student tidak ditemukan.")
			}
		case 4:
			before := len(students)
			students = removeInactive(students)

			// Simpan perubahan ke JSON setelah dihapus
			saveStudents(students)
			fmt.Printf("%d student inactive dihapus.\n", before-len(students))
		case 5:
			major := readString(scanner, "Jurusan: ")
			printStudents(filterByMajor(students, major))
		case 6:
			id := readInt(scanner, "ID student: ")
			studentByID := indexByID(students)
			if student, ok := studentByID[id]; ok {
				printStudents([]Student{student})
			} else {
				fmt.Println("Student tidak ditemukan.")
			}
		case 7:
			fmt.Printf("Rata-rata nilai: %.2f\n", averageScore(students))
		case 0:
			fmt.Println("Selesai.")
			return
		default:
			fmt.Println("Menu tidak valid.")
		}
	}
}

func printStudents(students []Student) {
	if len(students) == 0 {
		fmt.Println("Data student kosong.")
		return
	}

	sort.Slice(students, func(i, j int) bool {
		return students[i].Score > students[j].Score
	})

	fmt.Println("Daftar student:")
	for _, student := range students {
		status := "inactive"
		if student.Active {
			status = "active"
		}
		fmt.Printf("- ID %d | %s | %s | nilai %d | grade %s | %s\n",
			student.ID, student.Name, student.Major, student.Score, grade(student.Score), status)
	}
}

func readString(scanner *bufio.Scanner, label string) string {
	for {
		fmt.Print(label)
		if !scanner.Scan() {
			return ""
		}

		value := strings.TrimSpace(scanner.Text())
		if value != "" {
			return value
		}
		fmt.Println("Input tidak boleh kosong.")
	}
}

func readInt(scanner *bufio.Scanner, label string) int {
	for {
		value := readString(scanner, label)
		number, err := strconv.Atoi(value)
		if err == nil {
			return number
		}
		fmt.Println("Input harus angka.")
	}
}

func updateScore(students []Student, id int, score int) bool {
	for i := range students {
		if students[i].ID == id {
			students[i].Score = score
			return true
		}
	}
	return false
}

func removeInactive(students []Student) []Student {
	activeStudents := students[:0]
	for _, student := range students {
		if student.Active {
			activeStudents = append(activeStudents, student)
		}
	}
	return activeStudents
}

func filterByMajor(students []Student, major string) []Student {
	var result []Student
	for _, student := range students {
		if strings.EqualFold(student.Major, major) {
			result = append(result, student)
		}
	}
	return result
}

func indexByID(students []Student) map[int]Student {
	studentByID := make(map[int]Student, len(students))
	for _, student := range students {
		studentByID[student.ID] = student
	}
	return studentByID
}

func averageScore(students []Student) float64 {
	if len(students) == 0 {
		return 0
	}

	total := 0
	for _, student := range students {
		total += student.Score
	}
	return float64(total) / float64(len(students))
}

func grade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	default:
		return "D"
	}
}
