package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var jwtKey = []byte("Dbp4cR8Jwcsk2sTPT6cW") // Change this to a secure key

// User model
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

// Tour model
type Tour struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var db *gorm.DB

func initDB() {
	connStr := "host=localhost port=5433 user=myuser password=makakaonelove dbname=mydb sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Auto-migrate the updated User model
	db.AutoMigrate(&User{})

	log.Println("База данных успешно подключена и миграции выполнены")
}

// Hash Password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err

}

// Compare Hashed Password
// Compare Hashed Password
func checkPassword(hashed, password string) bool {
	log.Println(" Проверяем пароль...")
	log.Println(" Хеш из базы:", hashed)
	log.Println(" Введённый пароль:", password)

	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		log.Println(" Ошибка сравнения пароля:", err)
		return false
	}
	log.Println(" Пароль верный!")
	return true
}

// Generate JWT Token
func generateToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// Send Email
func sendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "baltabayev.adil@bk.ru") // Your email
	m.SetHeader("To", to)                        // Receiver's email
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.mail.ru", 587, "baltabayev.adil@bk.ru", "Dbp4cR8Jwcsk2sTPT6cW")
	return d.DialAndSend(m)
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"status": "fail", "message": "Некорректный запрос"}`, http.StatusBadRequest)
		return
	}

	// Хешируем пароль
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, `{"status": "fail", "message": "Ошибка при хешировании пароля"}`, http.StatusInternalServerError)
		return
	}

	log.Println("Исходный пароль:", user.Password)      // Выводим вводимый пароль
	log.Println("Хешированный пароль:", hashedPassword) // Выводим хеш пароля

	user.Password = hashedPassword

	// Проверяем, есть ли такой email
	var existingUser User
	db.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID != 0 {
		http.Error(w, `{"status": "fail", "message": "Пользователь с таким email уже существует"}`, http.StatusConflict)
		return
	}

	// Сохраняем пользователя
	db.Create(&user)

	log.Println("Пользователь успешно зарегистрирован:", user.Email)

	// Отправляем письмо
	emailErr := sendEmail(user.Email, "Добро пожаловать!", "Спасибо за регистрацию!")
	if emailErr != nil {
		log.Println("Ошибка отправки email:", emailErr)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Регистрация успешна! Проверьте почту.",
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Разбираем JSON-запрос
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"status": "fail", "message": "Некорректный запрос"}`, http.StatusBadRequest)
		return
	}

	// Ищем пользователя по email
	var user User
	db.Where("email = ?", creds.Email).First(&user)

	// Если пользователя нет
	if user.ID == 0 {
		log.Println("Ошибка: Пользователь не найден")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "fail",
			"message": "Такого пользователя не существует",
		})
		return
	}

	// Проверяем пароль
	if !checkPassword(user.Password, creds.Password) {
		log.Println("Ошибка: Неверный пароль для", creds.Email)
		log.Println("Хранимый хеш:", user.Password)
		log.Println("Введенный пароль:", creds.Password)

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "fail",
			"message": "Неверный пароль",
		})
		return
	}

	// Генерируем JWT-токен
	token, err := generateToken(user)
	if err != nil {
		http.Error(w, `{"status": "fail", "message": "Ошибка при генерации токена"}`, http.StatusInternalServerError)
		return
	}

	log.Println("Успешный вход в систему для:", creds.Email)

	// Успешный ответ
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Вы успешно вошли в аккаунт!",
		"token":   token,
	})
}

// Get all tours
func getTours(w http.ResponseWriter, r *http.Request) {
	var tours []Tour
	db.Find(&tours)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tours)
}

// Get a single tour by ID
func getTourByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/tours/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tour ID", http.StatusBadRequest)
		return
	}

	var tour Tour
	db.First(&tour, id)
	if tour.ID == 0 {
		http.Error(w, "Tour not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tour)
}

func main() {
	log.Println("Starting server...")
	initDB()

	// Раздаём ВСЕ файлы из public/
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// API маршруты
	http.HandleFunc("/api/register", registerHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/tours", getTours)
	http.HandleFunc("/api/tours/", getTourByID)

	srv := &http.Server{Addr: ":8080"}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exited gracefully")
}
