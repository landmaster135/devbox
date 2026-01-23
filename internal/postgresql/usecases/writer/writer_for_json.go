package dump

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type JSONStreamWriter struct {
	file     *os.File
	rows     int
	closed   bool
	closeErr error
}

func NewJSONStreamWriter(fileWriter FileWriter, filePath string) (*JSONStreamWriter, error) {
	file, err := fileWriter.Create(filePath)
	if err != nil {
		return nil, err
	}

	if _, err := file.WriteString("["); err != nil {
		_ = file.Close()
		return nil, err
	}

	return &JSONStreamWriter{file: file}, nil
}

func (w *JSONStreamWriter) WriteBatch(rows []map[string]any) error {
	if w.closed {
		return errors.New("既にクローズされたライターに書き込めません")
	}

	for _, row := range rows {
		rowJSON, err := json.MarshalIndent(row, "", "  ")
		if err != nil {
			return err
		}

		prefix := "\n"
		if w.rows > 0 {
			prefix = ",\n"
		}

		if _, err := w.file.WriteString(prefix); err != nil {
			return err
		}
		if _, err := w.file.WriteString("  "); err != nil {
			return err
		}
		if _, err := w.file.WriteString(strings.ReplaceAll(string(rowJSON), "\n", "\n  ")); err != nil {
			return err
		}

		w.rows++
	}

	return nil
}

func (w *JSONStreamWriter) Close() error {
	if w.closed {
		return w.closeErr
	}
	w.closed = true

	if w.file == nil {
		return nil
	}

	var err error
	if w.rows > 0 {
		_, err = w.file.WriteString("\n]")
	} else {
		_, err = w.file.WriteString("]")
	}

	if err != nil {
		w.closeErr = err
		_ = w.file.Close()
		return err
	}

	if closeErr := w.file.Close(); closeErr != nil {
		err = closeErr
	}

	w.closeErr = err
	return err
}

func (w *JSONStreamWriter) RowsWritten() int {
	return w.rows
}
