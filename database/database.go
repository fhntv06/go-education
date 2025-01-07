package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type user struct {
	id        int64
	username  string
	password  string
	createdAt time.Time
}

type DB struct {
	DB *sql.DB
}

func GetEnvParam(key string) string {
	param := os.Getenv(key)

	if key != "PASSWORD" && param == "" {
		log.Fatalf("Missing required environment variable: %s!", key)
	}

	return param
}
func GetConnectionParamsDB() string {
	// Загружаем переменные из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?%s",
		GetEnvParam("USER"),
		GetEnvParam("PASSWORD"),
		GetEnvParam("IP"),
		GetEnvParam("PORT"),
		GetEnvParam("DATABASE"),
		GetEnvParam("PARAMS"),
	)
}

// ConstructorDB

func ConstructorDB(paramsDB string) *DB {
	dbLocal, err := sql.Open("mysql", paramsDB)
	if err != nil {
		log.Fatal("Error in Open:", err)
	}
	if err := dbLocal.Ping(); err != nil {
		log.Fatal("Error in Ping:", err)
	}

	return &DB{DB: dbLocal}
}

func (dbLocal *DB) CreateUser(username string, password string, createdAt time.Time) sql.Result {
	result, err := dbLocal.DB.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, username, password, createdAt)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
func (dbLocal *DB) GetUserById(userId int64) *user {
	var (
		id        int64
		username  string
		password  string
		createdAt time.Time
	)

	query := "SELECT id, username, password, created_at FROM users WHERE id = ?"
	err := dbLocal.DB.QueryRow(query, userId).Scan(&id, &username, &password, &createdAt)

	if err != nil {
		log.Fatal("Error in GetUserById by id: ", userId, err)
	}

	return &user{id, username, password, createdAt}
}
func (dbLocal *DB) GetAllUsersFromTables(name string) []user {
	rows, err := dbLocal.DB.Query(fmt.Sprintf(`SELECT id, username, password, created_at FROM %s`, name))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []user
	for rows.Next() {
		var u user

		err := rows.Scan(&u.id, &u.username, &u.password, &u.createdAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return users
}
func (dbLocal *DB) CreateNewTables(name string) {
	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id INT AUTO_INCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME,
			PRIMARY KEY (id)
		);`,
		name,
	)

	if _, err := dbLocal.DB.Exec(query); err != nil {
		log.Fatal(err)
	}
}
func (dbLocal *DB) DeleteUserById(id int) {
	_, err := dbLocal.DB.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		log.Fatal("Error in DeleteUserById", err)
	}
}

func main() {
	connectionParamsDB := GetConnectionParamsDB()
	db := ConstructorDB(connectionParamsDB)

	//db.CreateNewTables("users")

	// Insert a new user
	result := db.CreateUser("johndoe23", "secret", time.Now())
	lastInsertId, _ := result.LastInsertId()
	fmt.Println("LastInsertId:", lastInsertId)

	// Query a single user
	lastUser := db.GetUserById(lastInsertId)
	fmt.Println(lastUser.id, lastUser.username, lastUser.password, lastUser.createdAt)

	// Query all users
	users := db.GetAllUsersFromTables("users")
	fmt.Printf("%#v", users)

	db.DeleteUserById(1)
}
