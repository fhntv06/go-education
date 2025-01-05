package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type User struct {
	Username    string
	Email       string
	Password    string
	ConfirmPass string
}
type ErrorResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Type string `json:"type"`
}

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
		timeStr, user.Username, user.Email, user.Password,
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
		log.Fatalf("Не удалось записать данные в файл: %v", err)
	}

	fmt.Println("Добавлен новый пользователь!")
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Получаем данные из формы
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		// Проверяем, что пароли совпадают
		if password != confirmPassword {
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
		user := User{
			Username: username,
			Email:    email,
			Password: password, // В реальном приложении следует хранить только хэш пароля
		}

		// логирование новых пользователей
		loggingCreateUser(user)

		// Выводим сообщение об успешной регистрации
		successResponse := &ErrorResponse{Text: "Пользователь " + user.Username + " зарегистрирован!", Type: "success"}
		sendJSONResponse(w, successResponse, http.StatusCreated)
	} else {
		// Если метод не POST, показываем форму
		fmt.Fprintf(w, "<h1>Метод не POST!</h1>")
	}
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

func main() {
	// Обработчики маршрутов
	r := mux.NewRouter()

	// Обработчик для корня сайта
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(".", "frontend", "index.html")) // конкатенация с учетом OS
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

	r.HandleFunc("/register", handleRegistration) // Обработка формы регистрации

	// Запуск сервера на порту 8080
	fmt.Println("Запуск сервера на http://localhost:80")
	if err := http.ListenAndServe(":80", r); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}

// моменты
// http позволило просто сделать роутинг и легко работает frontend со своими файлами, благодаря обработки рута из папки frontend
// gorilla/mux позволило испольловать в роутинге slug но отвалилать статика хотя рут смотрит в frontend
