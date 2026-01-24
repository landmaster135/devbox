package writer

import (
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type CSVStreamWriter struct {
	file     *os.File
	writer   *csv.Writer
	headers  []string
	rows     int
	closed   bool
	closeErr error
}

func NewCSVStreamWriter(fileWriter FileWriter, filePath string, headers []string) (*CSVStreamWriter, error) {
	file, err := fileWriter.Create(filePath)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)

	if len(headers) > 0 {
		if err := writer.Write(headers); err != nil {
			writer.Flush()
			_ = file.Close()
			return nil, err
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			_ = file.Close()
			return nil, err
		}
	}

	return &CSVStreamWriter{
		file:    file,
		writer:  writer,
		headers: headers,
	}, nil
}

func formatCSVValue(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		if len(v) == 0 {
			return ""
		}
		return base64.StdEncoding.EncodeToString(v)
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(value)
	}
}

func (w *CSVStreamWriter) WriteBatch(rows []map[string]any) error {
	if w.closed {
		return errors.New("既にクローズされたライターに書き込めません")
	}

	for _, row := range rows {
		values := make([]string, len(w.headers))
		for i, header := range w.headers {
			values[i] = formatCSVValue(row[header])
		}

		if err := w.writer.Write(values); err != nil {
			return err
		}
		w.rows++
	}

	w.writer.Flush()
	if err := w.writer.Error(); err != nil {
		return err
	}

	return nil
}

func (w *CSVStreamWriter) Close() error {
	if w.closed {
		return w.closeErr
	}
	w.closed = true

	if w.writer != nil {
		w.writer.Flush()
		if err := w.writer.Error(); err != nil {
			w.closeErr = err
			_ = w.file.Close()
			return err
		}
	}

	if w.file != nil {
		if err := w.file.Close(); err != nil {
			w.closeErr = err
			return err
		}
	}

	return nil
}

func (w *CSVStreamWriter) RowsWritten() int {
	return w.rows
}
