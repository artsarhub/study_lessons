package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

// Config содержит настройки для генерации данных
type Config struct {
	UsersCount    int
	ChatsCount    int
	MessagesCount int
	DBHost        string
	DBPort        int
	DBName        string
	DBUser        string
	DBPassword    string
}

// User представляет структуру пользователя
type User struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

// Chat представляет структуру чата
type Chat struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

// Message представляет структуру сообщения
type Message struct {
	ID        int
	Content   string
	AuthorID  int
	ChatID    int
	CreatedAt time.Time
}

func main() {
	// Парсинг флагов командной строки
	config := parseFlags()

	// Подключение к базе данных
	db, err := connectDB(config)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка ping БД:", err)
	}

	fmt.Println("Подключение к БД успешно!")

	// Генерация данных
	fmt.Printf("Генерация данных:\n- Пользователи: %d\n- Чаты: %d\n- Сообщения: %d\n",
		config.UsersCount, config.ChatsCount, config.MessagesCount)

	err = generateData(db, config)
	if err != nil {
		log.Fatal("Ошибка генерации данных:", err)
	}

	fmt.Println("Данные успешно сгенерированы!")
}

func parseFlags() *Config {
	config := &Config{}

	flag.IntVar(&config.UsersCount, "users", 100, "Количество пользователей")
	flag.IntVar(&config.ChatsCount, "chats", 20, "Количество чатов")
	flag.IntVar(&config.MessagesCount, "messages", 500, "Количество сообщений")
	flag.StringVar(&config.DBHost, "host", "localhost", "Хост БД")
	flag.IntVar(&config.DBPort, "port", 5432, "Порт БД")
	flag.StringVar(&config.DBName, "db", "messenger", "Имя БД")
	flag.StringVar(&config.DBUser, "user", "postgres", "Пользователь БД")
	flag.StringVar(&config.DBPassword, "password", "postgres", "Пароль БД")

	flag.Parse()

	return config
}

func connectDB(config *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func generateData(db *sql.DB, config *Config) error {
	// Очистка существующих данных
	err := clearExistingData(db)
	if err != nil {
		return err
	}

	// Генерация пользователей
	users, err := generateUsers(db, config.UsersCount)
	if err != nil {
		return err
	}

	// Генерация чатов
	chats, err := generateChats(db, config.ChatsCount)
	if err != nil {
		return err
	}

	// Генерация связей пользователей с чатами
	err = generateUsersChats(db, users, chats)
	if err != nil {
		return err
	}

	// Генерация сообщений
	err = generateMessages(db, config.MessagesCount, users, chats)
	if err != nil {
		return err
	}

	return nil
}

func clearExistingData(db *sql.DB) error {
	tables := []string{"Messages", "Users_Chats", "Chats", "Users"}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("ошибка очистки таблицы %s: %v", table, err)
		}
	}

	fmt.Println("Существующие данные очищены")
	return nil
}

func generateUsers(db *sql.DB, count int) ([]int, error) {
	userIDs := make([]int, 0, count)

	for i := 1; i <= count; i++ {
		user := User{
			ID:        i,
			Name:      generateUserName(),
			CreatedAt: generateRandomTime(),
		}

		_, err := db.Exec(
			"INSERT INTO Users (id, name, created_at) VALUES ($1, $2, $3)",
			user.ID, user.Name, user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка вставки пользователя: %v", err)
		}

		userIDs = append(userIDs, user.ID)
	}

	fmt.Printf("Сгенерировано %d пользователей\n", count)
	return userIDs, nil
}

func generateChats(db *sql.DB, count int) ([]int, error) {
	chatIDs := make([]int, 0, count)

	for i := 1; i <= count; i++ {
		chat := Chat{
			ID:        i,
			Name:      generateChatName(),
			CreatedAt: generateRandomTime(),
		}

		_, err := db.Exec(
			"INSERT INTO Chats (id, name, created_at) VALUES ($1, $2, $3)",
			chat.ID, chat.Name, chat.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка вставки чата: %v", err)
		}

		chatIDs = append(chatIDs, chat.ID)
	}

	fmt.Printf("Сгенерировано %d чатов\n", count)
	return chatIDs, nil
}

func generateUsersChats(db *sql.DB, userIDs, chatIDs []int) error {
	// Для каждого чата добавляем случайных пользователей
	for _, chatID := range chatIDs {
		// Количество пользователей в чате (от 2 до 10)
		usersInChat := rand.Intn(9) + 2
		if usersInChat > len(userIDs) {
			usersInChat = len(userIDs)
		}

		// Выбираем случайных пользователей
		selectedUsers := make(map[int]bool)
		for len(selectedUsers) < usersInChat {
			userID := userIDs[rand.Intn(len(userIDs))]
			selectedUsers[userID] = true
		}

		// Вставляем связи
		for userID := range selectedUsers {
			_, err := db.Exec(
				"INSERT INTO Users_Chats (user_id, chat_id) VALUES ($1, $2)",
				userID, chatID,
			)
			if err != nil {
				return fmt.Errorf("ошибка вставки связи пользователь-чат: %v", err)
			}
		}
	}

	fmt.Printf("Сгенерированы связи пользователей с чатами\n")
	return nil
}

func generateMessages(db *sql.DB, count int, userIDs, chatIDs []int) error {
	for i := 1; i <= count; i++ {
		message := Message{
			ID:        i,
			Content:   generateMessageContent(),
			AuthorID:  userIDs[rand.Intn(len(userIDs))],
			ChatID:    chatIDs[rand.Intn(len(chatIDs))],
			CreatedAt: generateRandomTime(),
		}

		_, err := db.Exec(
			"INSERT INTO Messages (id, content, author_id, created_at, chat_id) VALUES ($1, $2, $3, $4, $5)",
			message.ID, message.Content, message.AuthorID, message.CreatedAt, message.ChatID,
		)
		if err != nil {
			return fmt.Errorf("ошибка вставки сообщения: %v", err)
		}
	}

	fmt.Printf("Сгенерировано %d сообщений\n", count)
	return nil
}

// Вспомогательные функции для генерации данных

func generateUserName() string {
	firstNames := []string{"Алексей", "Мария", "Иван", "Ольга", "Дмитрий", "Елена", "Сергей", "Анна", "Андрей", "Наталья"}
	lastNames := []string{"Иванов", "Петров", "Сидоров", "Кузнецов", "Смирнов", "Попов", "Васильев", "Фёдоров", "Михайлов", "Новиков"}

	return fmt.Sprintf("%s %s", firstNames[rand.Intn(len(firstNames))], lastNames[rand.Intn(len(lastNames))])
}

func generateChatName() string {
	prefixes := []string{"Общий", "Рабочий", "Семейный", "Друзья", "Проект", "Команда", "Поддержка", "Разработка"}
	topics := []string{"чат", "канал", "обсуждение", "группа", "сообщество"}

	return fmt.Sprintf("%s %s %d", prefixes[rand.Intn(len(prefixes))], topics[rand.Intn(len(topics))], rand.Intn(1000))
}

func generateMessageContent() string {
	messages := []string{
		"Привет всем!",
		"Как дела?",
		"Что нового?",
		"Отличная работа!",
		"Давайте обсудим этот вопрос",
		"У меня есть предложение",
		"Согласен с вами",
		"Интересная идея",
		"Когда сможем встретиться?",
		"Спасибо за помощь!",
		"Всем хорошего дня!",
		"Жду ваших комментариев",
		"Отправляю файл",
		"Проверил, всё работает",
		"Нужна дополнительная информация",
	}

	return messages[rand.Intn(len(messages))]
}

func generateRandomTime() time.Time {
	// Генерируем время в пределах последних 365 дней
	now := time.Now()
	daysAgo := rand.Intn(365)
	hoursAgo := rand.Intn(24)
	minutesAgo := rand.Intn(60)

	return now.Add(-time.Duration(daysAgo)*24*time.Hour -
		time.Duration(hoursAgo)*time.Hour -
		time.Duration(minutesAgo)*time.Minute)
}
