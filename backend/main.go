package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
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

func main() {
	// Обработчики маршрутов

	// Обрабатываем статику из папки 'frontend'
	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/", fs) // теперь root смотрит в папку frontend

	// Перехватываем запросы и добавляем правильные заголовки Content-Type
	http.HandleFunc("/assets/javascripts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("assets/javascripts")
		// Определяем Content-Type по расширению файла
		w.Header().Set("Content-Type", "application/javascript")

		// Отдаем файл
		fs.ServeHTTP(w, r)
	})
	http.HandleFunc("/assets/javascripts/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("assets/javascripts/")
		// Определяем Content-Type по расширению файла
		w.Header().Set("Content-Type", "application/javascript")

		// Отдаем файл
		fs.ServeHTTP(w, r)
	})

	http.HandleFunc("/register", handleRegistration) // Обработка формы регистрации

	// Запуск сервера на порту 8080
	fmt.Println("Запуск сервера на http://localhost:80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
