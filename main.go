package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Хранилище метаданных о файлах
type FileMeta struct {
	Path      string
	ExpiresAt time.Time
}

var (
	files    = make(map[string]FileMeta)
	filesMux sync.Mutex
	uploadDir = "./uploads"
	// Время жизни файла
	fileTTL = 1 * time.Hour
	// Максимальный размер файла (500 MB)
	maxUploadSize int64 = 500 << 20
)

func main() {
	// Создаём директорию для загрузок
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Fatal(err)
	}

	// Запускаем очистку просроченных файлов
	go cleanupLoop()

	// Раздача статики
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/file/", downloadHandler)

	addr := ":8080"
	fmt.Printf("Сервер запущен: http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// generateID создаёт случайный идентификатор файла
func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// uploadHandler принимает файл и сохраняет его
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// Ограничиваем размер тела запроса
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
    	http.Error(w, "Файл слишком большой", http.StatusBadRequest)
    	return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Файл не найден в запросе", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Генерируем уникальное имя
	id := generateID()
	ext := filepath.Ext(header.Filename)
	savePath := filepath.Join(uploadDir, id+ext)

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Ошибка при записи файла", http.StatusInternalServerError)
		return
	}

	// Регистрируем файл
	filesMux.Lock()
	files[id] = FileMeta{
		Path:      savePath,
		ExpiresAt: time.Now().Add(fileTTL),
	}
	filesMux.Unlock()

	// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{
		"id": "%s",
		"name": "%s",
		"size": %d,
		"url": "/api/file/%s",
		"expires_in": %d
	}`, id, header.Filename, header.Size, id, int(fileTTL.Seconds()))
}

// downloadHandler отдаёт файл по id
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем id из URL: /api/file/{id}
	id := filepath.Base(r.URL.Path)
	if id == "" || id == "/" {
		http.Error(w, "Не указан id файла", http.StatusBadRequest)
		return
	}

	filesMux.Lock()
	meta, ok := files[id]
	if ok && time.Now().After(meta.ExpiresAt) {
		delete(files, id)
		os.Remove(meta.Path)
		ok = false
	}
	filesMux.Unlock()

	if !ok {
		http.Error(w, "Файл не найден или истёк срок хранения", http.StatusNotFound)
		return
	}

	// Проверка, что файл существует
	if _, err := os.Stat(meta.Path); err != nil {
		http.Error(w, "Файл недоступен", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, meta.Path)
}

// cleanupLoop периодически удаляет просроченные файлы
func cleanupLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		filesMux.Lock()
		for id, meta := range files {
			if now.After(meta.ExpiresAt) {
				os.Remove(meta.Path)
				delete(files, id)
				log.Printf("Удалён файл %s (истёк срок)", id)
			}
		}
		filesMux.Unlock()
	}
}