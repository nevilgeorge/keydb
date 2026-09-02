package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Container struct for a record.
type Record struct {
	index int
	key   string
	value string
}

type RecordStore interface {
	Get(key string) (string, error)
	Put(key string, value string) error
	Close() error
}

type RecordWriter interface {
	WriteRecord(record Record) (offset int64, err error)
	Close() error
}

type RecordReader interface {
	ReadRecord(offset int64) (*Record, error)
	Close() error
}

type FileRecordStore struct {
	filename    string
	writer      RecordWriter
	reader      RecordReader
	index       int
	keyToOffset map[string]int64
}

func (rs *FileRecordStore) Get(key string) (string, error) {
	offset, ok := rs.keyToOffset[key]
	if !ok {
		return "", fmt.Errorf("Missing key: %s", key)
	}

	record, err := rs.reader.ReadRecord(offset)
	if err != nil {
		return "", err
	}

	return record.value, nil
}

func (rs *FileRecordStore) Put(key string, value string) error {
	record := Record{rs.index, key, value}
	offset, err := rs.writer.WriteRecord(record)
	if err != nil {
		return err
	}
	rs.keyToOffset[key] = offset
	rs.index += 1
	return nil
}

func (rs *FileRecordStore) Close() error {
	// Close file handles in reader and writer.
	err := rs.reader.Close()
	if err != nil {
		return err
	}
	err = rs.writer.Close()
	if err != nil {
		return err
	}
	return nil
}

type FileRecordReader struct {
	file *os.File
}

type FileRecordWriter struct {
	file *os.File
}

func (w *FileRecordWriter) WriteRecord(record Record) (int64, error) {
	line := fmt.Sprintf("%d\t%s\t%s\n", record.index, record.key, record.value)
	offset, err := w.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	_, err = w.file.WriteString(line)
	if err != nil {
		return 0, err
	}
	return offset, err
}

func (w *FileRecordWriter) Close() error {
	return w.file.Close()
}

func (r *FileRecordReader) ReadRecord(offset int64) (*Record, error) {
	if _, err := r.file.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	line, err := bufio.NewReader(r.file).ReadString('\n')
	record, err := r.parseLineIntoRecord(line)
	return record, err
}

func (r *FileRecordReader) parseLineIntoRecord(line string) (*Record, error) {
	fields := strings.Split(line, "\t")
	if len(fields) != 3 {
		return nil, fmt.Errorf("Expected 3 fields, got %d: %s", len(fields), line)
	}
	index, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("Expected integer index, got %s: %s", fields[0], line)
	}
	return &Record{index: index, key: fields[1], value: fields[2]}, nil
}

func (r *FileRecordReader) Close() error {
	return r.file.Close()
}

func newFileRecordWriter(filename string) (*FileRecordWriter, error) {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &FileRecordWriter{file: f}, nil
}

func newFileRecordReader(filename string) (*FileRecordReader, error) {
	f, err := os.Open(filename) // read-only
	if err != nil {
		return nil, err
	}
	return &FileRecordReader{file: f}, nil
}

func NewFileRecordStore(filename string) (RecordStore, error) {
	recordWriter, err := newFileRecordWriter(filename)
	if err != nil {
		return &FileRecordStore{}, err
	}
	recordReader, err := newFileRecordReader(filename)
	if err != nil {
		return &FileRecordStore{}, err
	}
	return &FileRecordStore{
		filename:    filename,
		writer:      recordWriter,
		reader:      recordReader,
		keyToOffset: make(map[string]int64),
	}, nil
}
