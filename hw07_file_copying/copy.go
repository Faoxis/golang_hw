package main

import (
	"errors"
	"io"
	"log"
	"os"

	//nolint:depguard
	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile          = errors.New("unsupported file")
	ErrOffsetExceedsFileSize    = errors.New("offset exceeds file size")
	ErrFromAndToPathsAreTheSame = errors.New("from path and to path are the same")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := checkSameFile(fromPath, toPath); err != nil {
		return err
	}
	if err := checkFileFinite(fromPath); err != nil {
		return err
	}

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
		return ErrOffsetExceedsFileSize
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

	totalSize := fileInfo.Size() - offset
	if limit > 0 && limit < totalSize {
		totalSize = limit
	}
	if limit <= 0 || limit > totalSize {
		limit = totalSize
	}
	bar := progressbar.DefaultBytes(totalSize, "Copying")

	if err := copyData(fromFile, toFile, limit, bar); err != nil {
		return err
	}
	return nil
}

func copyData(fromFile *os.File, toFile *os.File, limit int64, bar *progressbar.ProgressBar) error {
	_, err := io.CopyN(io.MultiWriter(toFile, bar), fromFile, limit)
	if err != nil {
		return err
	}
	return nil
}

func checkSameFile(from, to string) error {
	srcInfo, err := os.Stat(from)
	if err != nil {
		return err
	}

	dstInfo, err := os.Stat(to)
	if err == nil { // Если файл назначения уже существует, проверяем, не совпадает ли он с исходным
		if os.SameFile(srcInfo, dstInfo) {
			return ErrFromAndToPathsAreTheSame
		}
	}
	return nil
}

func checkFileFinite(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return ErrUnsupportedFile
	}

	// Проверяем, является ли файл символьным устройством (например, /dev/urandom).
	if info.Mode()&os.ModeCharDevice != 0 {
		return ErrUnsupportedFile
	}

	// Проверяем, является ли файл именованным каналом (pipe), что тоже бесконечный источник.
	if info.Mode()&os.ModeNamedPipe != 0 {
		return ErrUnsupportedFile
	}

	// Если файл имеет нулевой размер, он может быть либо пустым, либо бесконечным потоком.
	if info.Size() == 0 {
		// Попробуем открыть файл и считать 1 байт.
		f, err := os.Open(path)
		if err != nil {
			return ErrUnsupportedFile
		}
		defer func() {
			err := f.Close()
			log.Printf("ERROR: %v", err)
		}()

		buf := make([]byte, 1)
		readBytes, err := f.Read(buf)
		if err != nil {
			return err
		}

		// Если удалось что-то прочитать, значит, файл выдаёт данные бесконечно.
		if readBytes > 0 {
			return ErrUnsupportedFile
		}
	}

	// Если файл не символьное устройство, не pipe и имеет фиксированный размер — он нормальный.
	return nil
}
