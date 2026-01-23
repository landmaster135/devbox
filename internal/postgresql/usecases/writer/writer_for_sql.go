package dump

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	pq "github.com/lib/pq"

	sql "github.com/landmaster135/devbox/internal/postgresql/domain/sql"
)

type SQLStreamWriter struct {
	file          *os.File
	tableName     string
	quotedTable   string
	columns       []string
	quotedColumns []string
	headerWritten bool
	rows          int
	closed        bool
	closeErr      error
	generatedAt   time.Time
}

func NewSQLStreamWriter(fileWriter FileWriter, filePath, tableName string, columns []string) (*SQLStreamWriter, error) {
	file, err := fileWriter.Create(filePath)
	if err != nil {
		return nil, err
	}

	quotedTable, _, err := sql.QuoteQualifiedTableName(tableName)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = pq.QuoteIdentifier(col)
	}

	return &SQLStreamWriter{
		file:          file,
		tableName:     tableName,
		quotedTable:   quotedTable,
		columns:       columns,
		quotedColumns: quotedColumns,
		generatedAt:   time.Now(),
	}, nil
}

func formatSQLValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "NULL"
	case string:
		return pq.QuoteLiteral(v)
	case []byte:
		if len(v) == 0 {
			return "decode('', 'hex')"
		}
		hexStr := hex.EncodeToString(v)
		return fmt.Sprintf("decode('%s','hex')", hexStr)
	case time.Time:
		return pq.QuoteLiteral(v.Format(time.RFC3339Nano))
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
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
		return pq.QuoteLiteral(v.String())
	default:
		return pq.QuoteLiteral(fmt.Sprint(value))
	}
}

func (w *SQLStreamWriter) WriteBatch(rows []map[string]any) error {
	if w.closed {
		return errors.New("既にクローズされたライターに書き込めません")
	}

	if len(rows) == 0 {
		return nil
	}

	if !w.headerWritten {
		if _, err := fmt.Fprintf(w.file, "-- Table dump for %s\n", w.tableName); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w.file, "-- Generated at %s\n\n", w.generatedAt.Format("2006-01-02 15:04:05")); err != nil {
			return err
		}
		w.headerWritten = true
	}

	for _, row := range rows {
		if _, err := fmt.Fprintf(w.file, "INSERT INTO %s (%s) VALUES (", w.quotedTable, strings.Join(w.quotedColumns, ", ")); err != nil {
			return err
		}

		values := make([]string, len(w.columns))
		for i, col := range w.columns {
			values[i] = formatSQLValue(row[col])
		}

		if _, err := w.file.WriteString(strings.Join(values, ", ")); err != nil {
			return err
		}
		if _, err := w.file.WriteString(");\n"); err != nil {
			return err
		}

		w.rows++
	}

	return nil
}

func (w *SQLStreamWriter) Close() error {
	if w.closed {
		return w.closeErr
	}
	w.closed = true

	if w.file == nil {
		return nil
	}

	if w.rows == 0 {
		if _, err := w.file.WriteString("-- No data to export\n"); err != nil {
			w.closeErr = err
			_ = w.file.Close()
			return err
		}
	}

	if err := w.file.Close(); err != nil {
		w.closeErr = err
		return err
	}

	return nil
}

func (w *SQLStreamWriter) RowsWritten() int {
	return w.rows
}
