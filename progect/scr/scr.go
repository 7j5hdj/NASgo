package scr

import (
	dat "nasgo/progect/data"

	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type FileNode struct {
	Name     string               `json:"name"`
	Children map[string]*FileNode `json:"children,omitempty"`
}

type Storage struct {
	mu    sync.Mutex
	files map[string][]byte
}

// 1
func NewStorage() *Storage {
	return &Storage{
		files: make(map[string][]byte),
	}
}

// Unzip распаковывает содержимое zip-архива из памяти в папку назначения
func Unzip(data []byte, targetDir string) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(targetDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		fileInZip, err := file.Open()
		if err != nil {
			return err
		}

		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileInZip.Close()
			return err
		}

		_, err = io.Copy(dstFile, fileInZip)
		dstFile.Close()
		fileInZip.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// 3
func (s *Storage) Save(key string, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Если загружается ZIP, распаковываем его содержимое
	if filepath.Ext(key) == ".zip" {
		log.Printf("Обнаружен архив %s, начинаю распаковку...", key)
		if err := Unzip(data, dat.STORAGE_DIR); err != nil {
			log.Printf("Ошибка распаковки: %v", err)
		}
	}

	// Сохраняем сам файл (или архив) в память и на диск
	s.files[key] = data
	fullPath := filepath.Join(dat.STORAGE_DIR, key)
	os.MkdirAll(filepath.Dir(fullPath), 0755)

	err := ioutil.WriteFile(fullPath, data, 0644)
	if err != nil {
		log.Printf("Ошибка при сохранении %s: %v", key, err)
	}
}

// 4
func (s *Storage) Load(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.files[key]
	if exists {
		return data, true
	}

	data, err := ioutil.ReadFile(filepath.Join(dat.STORAGE_DIR, key))
	if err != nil {
		return nil, false
	}

	s.files[key] = data
	return data, true
}

// 5
func RenderTree(w io.Writer, path string, indent string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for i, file := range files {
		isLast := i == len(files)-1
		prefix := "├── "
		if isLast {
			prefix = "└── "
		}

		// Выводим имя текущего файла/папки
		fmt.Fprintf(w, "%s%s%s\n", indent, prefix, file.Name())

		// Если это папка, ныряем глубже (рекурсия)
		if file.IsDir() {
			newIndent := indent + "│   "
			if isLast {
				newIndent = indent + "    "
			}
			RenderTree(w, filepath.Join(path, file.Name()), newIndent)
		}
	}
	return nil
}
