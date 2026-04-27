package handle

import (
	dat "nasgo/progect/data"
	"nasgo/progect/scr"

	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func HandleUpload(w http.ResponseWriter, r *http.Request, storage *scr.Storage) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Path[dat.UPLOAD_PREFIX_LEN:]
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}

	storage.Save(key, data)
	fmt.Fprintf(w, "Объект %s обработан", key)
}

func HandleDownload(w http.ResponseWriter, r *http.Request, storage *scr.Storage) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Path[dat.DOWNLOAD_PREFIX_LEN:]
	data, exists := storage.Load(key)
	if !exists {
		http.Error(w, "Объект не найден", http.StatusNotFound)
		return
	}
	w.Write(data)
}

func HandleList(w http.ResponseWriter, r *http.Request, storage *scr.Storage) {

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем тип контента как обычный текст
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintln(w, "storage") // Корень хранилища

	// Запускаем отрисовку дерева из папки STORAGE_DIR
	err := scr.RenderTree(w, dat.STORAGE_DIR, "")
	if err != nil {
		http.Error(w, "Ошибка при чтении файлов", http.StatusInternalServerError)
	}
}

func HandleCommands(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Читаем файл commands.txt из корня проекта
	content, err := os.ReadFile("commands.txt")
	if err != nil {
		// Если файла нет, выводим ошибку
		log.Printf("Ошибка чтения commands.txt: %v", err)
		http.Error(w, "Файл с командами не найден", http.StatusNotFound)
		return
	}

	// Указываем тип контента и отдаем содержимое файла
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
