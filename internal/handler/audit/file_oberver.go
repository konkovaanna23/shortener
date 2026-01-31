package audit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type FileObserver struct {
	file *os.File
	buf  *bufio.Writer
	mu   sync.Mutex
}

// NewFileObserver создаёт наблюдатель, который записывает события в конец файла.
func NewFileObserver(filename string) (*FileObserver, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &FileObserver{
		file: file,
		buf:  bufio.NewWriterSize(file, 64*1024), // 64 KB buffer
	}, nil
}

// Notify записывает событие в файл в формате JSON, добавляя перевод строки.
func (f *FileObserver) Notify(event Event) {
	fmt.Println("ЗАпрос на запись в файл")
	f.mu.Lock()
	defer f.mu.Unlock()
	fmt.Println("ЗАпрос на запись в файл")
	data, err := json.Marshal(event)
	if err != nil {
		logrus.Error("Не удалось сериализовать событие аудита", err)
		return
	}

	if _, err := f.buf.Write(data); err != nil {
		logrus.Error("Ошибка записи данных", err)
		return
	}

	if err := f.buf.WriteByte('\n'); err != nil {
		logrus.Error("Ошибка записи символа новой строки", err)
		return
	}

	if err := f.buf.Flush(); err != nil {
		logrus.Error("Ошибка сброса буфера на диск", err)
		return
	}
}

// Close закрывает файл и сбрасывает буфер.
func (f *FileObserver) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.file.Close()
}
