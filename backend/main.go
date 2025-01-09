package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type User struct {
	username    string
	email       string
	password    string
	confirmpass string
}
type ErrorResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Type string `json:"type"`
}
type DB struct {
	DB *sql.DB
}

var globalDB *sql.DB

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

// Функция для отправки JSON-ответа с ошибкой
func sendJSONResponse(w http.ResponseWriter, response *ErrorResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func loggingCreateUser(user User) {
	// Открываем файл для дозаписи (если файл не существует, он будет создан)
	/*

		Используется функция os.OpenFile с флагами:
		os.O_APPEND: открывает файл для добавления данных в конец;
		os.O_CREATE: создает файл, если он еще не существует;
		os.O_WRONLY: открывает файл только для записи.

	*/
	// Формируем строку для записи в лог-файл

	// Получаем текущее время
	now := time.Now()

	// Форматируем время в нужный формат
	timeStr := fmt.Sprintf("%s %s", now.Format("15:04:05"), now.Format("02-01-2006"))
	userLog := fmt.Sprintf(
		"--> Создан новый пользователь:\n\tDate: %s,\n\tUsername: %s,\n\tEmail: %s,\n\tPassword: %s\n--------------------------\n",
		timeStr, user.username, user.email, user.password,
	)

	// Получаем текущую рабочую директорию (проекта)
	currentDir, err := os.Getwd()

	// Формируем полный путь до файла лога
	logPath := filepath.Join(currentDir, "log.txt")

	// Открываем файл для записи (или создаем новый, если его нет) и добавляем данные в конец
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %v", err)
	}
	defer file.Close()

	// Записываем сформированную строку в файл
	_, err = file.WriteString(userLog)
	if err != nil {
		log.Fatalf("Не удалось записать данные в файл логов: %v", err)
	}

	fmt.Println("Добавлен новый пользователь!")
}
func loggingError(err error) {
	now := time.Now()

	// Форматируем время в нужный формат
	timeStr := fmt.Sprintf("%s %s", now.Format("15:04:05"), now.Format("02-01-2006"))
	userLog := fmt.Sprintf("--> Новая ошибка:\n\tDate: %s\n\tError: %s\n--------------------------\n", timeStr, err)

	// Получаем текущую рабочую директорию (проекта)
	currentDir, err := os.Getwd()

	// Формируем полный путь до файла лога
	logPath := filepath.Join(currentDir, "log.txt")

	// Открываем файл для записи (или создаем новый, если его нет) и добавляем данные в конец
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл: %v", err)
	}
	defer file.Close()

	// Записываем сформированную строку в файл
	_, err = file.WriteString(userLog)
	if err != nil {
		log.Fatalf("Не удалось записать данные в файл логов ошибки: %v", err)
	}

	fmt.Println("Возникла ошибка в backend! Проверьте логи!")
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	// Получаем данные из формы
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	hashPassword, errPassword := HashPassword(password)
	hashConfirmPassword, errConfirmPassword := HashPassword(confirmPassword)

	if errPassword != nil || errConfirmPassword != nil {
		if errConfirmPassword != nil {
			loggingError(errConfirmPassword)
		} else {
			loggingError(errPassword)
		}

		errorResponse := &ErrorResponse{ID: "confirm_password", Text: "Ошибка проверки паролей на сервере!", Type: "error"}
		sendJSONResponse(w, errorResponse, http.StatusInternalServerError)
		return
	}

	// Проверяем, что пароли совпадают
	if hashPassword != hashConfirmPassword {
		errorResponse := &ErrorResponse{ID: "confirm_password", Text: "Пароли не совпадают", Type: "error"}
		sendJSONResponse(w, errorResponse, http.StatusBadRequest)
		return
	}

	// Простая валидация
	// Имя пользователя
	if len(username) == 0 {
		errorResponse := &ErrorResponse{ID: "username", Text: "Некорректные данные. Проверьте поле username.", Type: "error"}
		sendJSONResponse(w, errorResponse, http.StatusBadRequest)
		return
	}
	// Пароль
	if len(password) < 6 {
		errorResponse := &ErrorResponse{ID: "password", Text: "Пароль слишком короткий. Минимум 6 символов.", Type: "error"}
		sendJSONResponse(w, errorResponse, http.StatusBadRequest)
		return
	}
	// Email
	if !strings.Contains(email, "@") {
		errorResponse := &ErrorResponse{ID: "email", Text: "Некорректные данные. Проверьте поле Email.", Type: "error"}
		sendJSONResponse(w, errorResponse, http.StatusBadRequest)
		return
	}

	// Создаем нового пользователя
	user := User{username, email, hashPassword, hashConfirmPassword}

	// логирование новых пользователей
	loggingCreateUser(user)

	// Insert a new user
	CreateUser(username, email, hashPassword, time.Now())

	// Выводим сообщение об успешной регистрации
	successResponse := &ErrorResponse{Text: "Пользователь " + username + " зарегистрирован!", Type: "success"}
	sendJSONResponse(w, successResponse, http.StatusCreated)
}

func GetEnvParam(key string) string {
	// получаем СИСТЕМНЫЕ ПЕРЕМЕННЫЕ
	// если key определен в системе ранее, то его значение будет не корректо относительно вашего файла .env
	param := os.Getenv(key)

	if key != "PASSWORD" && param == "" {
		log.Fatalf("Missing required environment variable: %s!", key)
	}

	return param
}

// Функция для определения Content-Type по расширению файла
func getContentTypeByExtension(filePath string) string {
	switch filepath.Ext(filePath) {
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".html", ".htm":
		return "text/html"
	default:
		return "application/octet-stream"
	}
}
func CreateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	successResponse := &ErrorResponse{Text: "Creating Book with " + title + "!", Type: "success"}
	sendJSONResponse(w, successResponse, http.StatusCreated)
	//fmt.Fprintf(w, "Creating Book")
}
func ReadBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	title := vars["title"]

	fmt.Fprintf(w, "ReadBook Book is %s", title)
}
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "UpdateBook Book")
}
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DeleteBook Book")
}
func builderPathStaticFileHandler(dirName string) http.Handler {
	return http.StripPrefix(
		"/"+dirName+"/",
		http.FileServer(http.Dir(filepath.Join(".", "frontend", dirName))),
	)
}
func accountHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		w.WriteHeader(http.StatusForbidden)
		http.ServeFile(w, r, filepath.Join("..", "frontend", "forbidden.html")) // конкатенация с учетом OS
	} else {
		fmt.Println("User is logged in!")
		http.ServeFile(w, r, filepath.Join("..", "frontend", "account.html")) // конкатенация с учетом OS
	}
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...
	fmt.Println("Authentication is true! Is logged in!")

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)

	http.ServeFile(w, r, filepath.Join("..", "frontend", "login.html")) // конкатенация с учетом OS
}
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	fmt.Println("Authentication is false! Is logout!")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)

	http.ServeFile(w, r, filepath.Join("..", "frontend", "logout.html")) // конкатенация с учетом OS
}

// TODO: дописать процессы
// ++ 1) Пользователь регистрируется -> register + заносим в базу данные пользователя
// 2) * Пользователь после регистрации попадает в свой аккаунт => редирект после получения данных с сервера об успешной регистрации
// 3) * В аккаунте видны данных пользователя => на страницу добавить прелоадер + по загрузке данных (REST API + slug) вывести на страницу
// 4) * Пользователь выходит из аккаунта (кнопка "Выйти") и
//	=> при нажатии на "Выйти" идет запрос => убирается сессия => приходит ответ success
//	=> редирект на страницу logout
// 5) * Перенаправление пользователя после таймера 5 сек с logout на страницу login (переход реализован на js)
// 6) * Пользователь авторизуется
//	=> при нажатии на "Войти" идет запрос => создается сессия => приходит ответ success
//	=> редирект на страницу account

// * прежде починить роутинг чтобы подгрузился build

// ConstructorDB

func ConstructorDB(paramsDB string) {
	dbLocal, err := sql.Open("mysql", paramsDB)
	if err != nil {
		log.Fatal("Error in Open: ", err)
	}
	if err := dbLocal.Ping(); err != nil {
		log.Fatal("Error in Ping: ", err)
	}

	globalDB = dbLocal
}
func CreateUser(username string, email string, password string, createdAt time.Time) sql.Result {
	if globalDB == nil {
		log.Fatal("Ошибка: globalDB пуста!")
	}

	result, err := globalDB.Exec(
		`INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)`,
		username, email, password, createdAt,
	)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
func CreateNewTables(name string) {
	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s (
			id INT AUTO_INCREMENT,
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME,
			PRIMARY KEY (id)
		);`,
		name,
	)

	if _, err := globalDB.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func initialRouting() *mux.Router {
	// Обработчики маршрутов
	r := mux.NewRouter()

	// Обработчик для корня сайта
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("..", "frontend", "index.html")) // конкатенация с учетом OS
	})

	// Маршрут для статических файлов из папки frontend/build
	r.PathPrefix("/build/").Handler(builderPathStaticFileHandler("build"))
	// Маршрут для статических файлов из папки frontend/assets
	r.PathPrefix("/assets/").Handler(builderPathStaticFileHandler("assets"))

	// Разграничение по методам
	r.HandleFunc("/books/{title}", CreateBook).Methods("POST")
	r.HandleFunc("/books/{title}", ReadBook).Methods("GET")
	r.HandleFunc("/books/{title}", UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{title}", DeleteBook).Methods("DELETE")

	// Обработка запроса по slug'ам
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "R: %s", r)

		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	r.HandleFunc("/register", handleRegistration).Methods("POST") // Обработка формы регистрации
	r.HandleFunc("/account", accountHandler)                      // Обработка перехода на страницу account
	r.HandleFunc("/login", loginHandler)                          // Обработка перехода на страницу login
	r.HandleFunc("/logout", logoutHandler)                        // Обработка перехода на страницу logout

	return r
}
func initialServer() {
	r := initialRouting()

	// Запуск сервера на порту 8080
	fmt.Println("Запуск сервера на http://localhost:80")
	if err := http.ListenAndServe(":80", r); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
func initialDatabase() {
	// Загружаем переменные из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connectionParamsDB := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?%s",
		GetEnvParam("USER"),
		GetEnvParam("PASSWORD"),
		GetEnvParam("IP"),
		GetEnvParam("PORT"),
		GetEnvParam("DATABASE"),
		GetEnvParam("PARAMS"),
	)

	// ConstructorDB
	ConstructorDB(connectionParamsDB)
	CreateNewTables("users")
}

// Hashing password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	initialDatabase()
	initialServer()
}

// моменты
// http позволило просто сделать роутинг и легко работает frontend со своими файлами, благодаря обработки рута из папки frontend
// gorilla/mux позволило испольловать в роутинге slug но отвалилать статика хотя рут смотрит в frontend
