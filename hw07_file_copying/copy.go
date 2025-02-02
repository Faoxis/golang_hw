package main

import (
	"errors"
	"io"
	"log"
	"os"
	"time"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем исходный файл
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		fileCloseErr := fromFile.Close()
		if fileCloseErr != nil {
			log.Printf("ERROR: error while closing file %v:%v", fromPath, fileCloseErr)
		}
	}()

	// Получаем информацию о файле
	fileInfo, err := fromFile.Stat()
	if err != nil {
		return err
	}

	if offset > fileInfo.Size() {
		return errors.New("offset больше размера файла")
	}

	// Устанавливаем смещение
	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// Создаем файл назначения
	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		fileCloseErr := toFile.Close()
		if fileCloseErr != nil {
			log.Printf("ERROR: error while closing file %v:%v", toFile, fileCloseErr)
		}
	}()

	// Копируем данные с ограничением
	buffer := make([]byte, 1024) // 1 KB
	var copied int64

	for {
		if limit > 0 && copied >= limit {
			break
		}

		chunkSize := int64(len(buffer))
		if limit > 0 && copied+chunkSize > limit {
			chunkSize = limit - copied
		}

		n, err := fromFile.Read(buffer[:chunkSize])
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if n == 0 {
			break
		}

		_, err = toFile.Write(buffer[:n])
		if err != nil {
			return err
		}

		copied += int64(n)
		time.Sleep(10 * time.Millisecond) // Задержка для ограничения скорости
	}

	return nil
}
