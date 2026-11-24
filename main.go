package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
<<<<<<< HEAD
	"mime"
	"net/http"
	"net/url"
=======
	"net/http"
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Хранилище метаданных о файлах
type FileMeta struct {
	Path      string
<<<<<<< HEAD
	Name      string // исходное имя файла (для Content-Disposition)
=======
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	ExpiresAt time.Time
}

var (
<<<<<<< HEAD
	files     = make(map[string]FileMeta)
	filesMux  sync.Mutex
=======
	files    = make(map[string]FileMeta)
	filesMux sync.Mutex
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	uploadDir = "./uploads"
	// Время жизни файла
	fileTTL = 1 * time.Hour
	// Максимальный размер файла (500 MB)
	maxUploadSize int64 = 500 << 20
)

func main() {
<<<<<<< HEAD
	uploadDir = "./uploads"
	if d := os.Getenv("UPLOAD_DIR"); d != "" {
		uploadDir = d
	}
=======
	// Создаём директорию для загрузок
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Fatal(err)
	}

<<<<<<< HEAD
	go cleanupLoop()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
=======
	// Запускаем очистку просроченных файлов
	go cleanupLoop()

	// Раздача статики
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// API
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	http.HandleFunc("/api/upload", uploadHandler)
	http.HandleFunc("/api/file/", downloadHandler)

	addr := ":8080"
<<<<<<< HEAD
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	fmt.Printf("Сервер запущен: %s\n", addr)
=======
	fmt.Printf("Сервер запущен: http://localhost%s\n", addr)
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
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

<<<<<<< HEAD
// extractFilename пытается вытащить корректное имя из Content-Disposition
func extractFilename(h map[string][]string) string {
	cd := ""
	if v, ok := h["Content-Disposition"]; ok && len(v) > 0 {
		cd = v[0]
	}
	if cd == "" {
		return ""
	}
	_, params, err := mime.ParseMediaType(cd)
	if err != nil {
		return ""
	}
	if fn, ok := params["filename"]; ok {
		return fn
	}
	return ""
}

=======
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
// uploadHandler принимает файл и сохраняет его
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	// Ограничиваем размер тела запроса
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
<<<<<<< HEAD

	log.Printf("UPLOAD: ct=%q content-length=%d", r.Header.Get("Content-Type"), r.ContentLength)

	// 32 MB в памяти, остальное — во временных файлах на диске
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("PARSE ERROR: %v", err)
		http.Error(w, "Ошибка разбора формы: "+err.Error(), http.StatusBadRequest)
		return
=======
	if err := r.ParseMultipartForm(32 << 20); err != nil {
    	http.Error(w, "Файл слишком большой", http.StatusBadRequest)
    	return
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	}

	file, header, err := r.FormFile("file")
	if err != nil {
<<<<<<< HEAD
		log.Printf("FORMFILE ERROR: %v", err)
=======
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
		http.Error(w, "Файл не найден в запросе", http.StatusBadRequest)
		return
	}
	defer file.Close()

<<<<<<< HEAD
	// --- Определяем исходное имя файла ---
	// Приоритет 1: отдельное поле формы originalName (самый надёжный способ)
	origName := r.FormValue("originalName")

	// Приоритет 2: декодируем Content-Disposition (RFC 5987 + fallback)
	if origName == "" {
		origName = extractFilename(header.Header)
	}

	// Приоритет 3: сырой header.Filename с percent-декодированием
	if origName == "" {
		raw := header.Filename
		if decoded, derr := url.QueryUnescape(raw); derr == nil && decoded != raw {
			origName = decoded
		} else {
			origName = raw
		}
	}

	if origName == "" {
		origName = "file"
	}

	// Отрезаем путь, если браузер его прислал (на всякий случай)
	origName = filepath.Base(origName)

	log.Printf("GOT FILE: header.Filename=%q origName=%q size=%d", header.Filename, origName, header.Size)

	// Генерируем уникальное имя на диске
	id := generateID()
	ext := filepath.Ext(origName)
=======
	// Генерируем уникальное имя
	id := generateID()
	ext := filepath.Ext(header.Filename)
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	savePath := filepath.Join(uploadDir, id+ext)

	dst, err := os.Create(savePath)
	if err != nil {
<<<<<<< HEAD
		log.Printf("CREATE ERROR: %v", err)
=======
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
		http.Error(w, "Не удалось сохранить файл", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

<<<<<<< HEAD
	written, err := io.Copy(dst, file)
	if err != nil {
		log.Printf("COPY ERROR: %v", err)
		os.Remove(savePath)
=======
	if _, err := io.Copy(dst, file); err != nil {
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
		http.Error(w, "Ошибка при записи файла", http.StatusInternalServerError)
		return
	}

	// Регистрируем файл
	filesMux.Lock()
	files[id] = FileMeta{
		Path:      savePath,
<<<<<<< HEAD
		Name:      origName,
=======
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
		ExpiresAt: time.Now().Add(fileTTL),
	}
	filesMux.Unlock()

<<<<<<< HEAD
	// Формируем JSON-ответ
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
=======
	// Формируем ответ
	w.Header().Set("Content-Type", "application/json")
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	fmt.Fprintf(w, `{
		"id": "%s",
		"name": "%s",
		"size": %d,
		"url": "/api/file/%s",
		"expires_in": %d
<<<<<<< HEAD
	}`, id, jsonEscape(origName), written, id, int(fileTTL.Seconds()))
}

// downloadHandler отдаёт файл по id с исходным именем
=======
	}`, id, header.Filename, header.Size, id, int(fileTTL.Seconds()))
}

// downloadHandler отдаёт файл по id
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

<<<<<<< HEAD
	id := filepath.Base(r.URL.Path)
	if id == "" || id == "/" || id == "." {
=======
	// Извлекаем id из URL: /api/file/{id}
	id := filepath.Base(r.URL.Path)
	if id == "" || id == "/" {
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
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

<<<<<<< HEAD
=======
	// Проверка, что файл существует
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
	if _, err := os.Stat(meta.Path); err != nil {
		http.Error(w, "Файл недоступен", http.StatusNotFound)
		return
	}

<<<<<<< HEAD
	// Отдаём с исходным именем
	w.Header().Set("Content-Disposition",
		"attachment; filename*=UTF-8''"+url.PathEscape(meta.Name))

	http.ServeFile(w, r, meta.Path)
}

// jsonEscape экранирует спецсимволы для безопасной вставки в JSON-строку
func jsonEscape(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '\\':
			out = append(out, '\\', '\\')
		case '"':
			out = append(out, '\\', '"')
		case '\n':
			out = append(out, '\\', 'n')
		case '\r':
			out = append(out, '\\', 'r')
		case '\t':
			out = append(out, '\\', 't')
		default:
			if c < 0x20 {
				out = append(out, []byte(fmt.Sprintf(`\u%04x`, c))...)
			} else {
				out = append(out, c)
			}
		}
	}
	return string(out)
}

=======
	http.ServeFile(w, r, meta.Path)
}

>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
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
<<<<<<< HEAD
				log.Printf("Удалён файл %s (%s, истёк срок)", id, meta.Name)
=======
				log.Printf("Удалён файл %s (истёк срок)", id)
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
			}
		}
		filesMux.Unlock()
	}
<<<<<<< HEAD
}
=======
}
>>>>>>> 9df3276ac66afacc22ef63b8bfe2173396a33b6c
