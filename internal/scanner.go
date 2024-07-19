package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/dustin/go-humanize"
)

type Entity struct {
	Checksum string
	FileSize int64
	Files    []string
}

type Option func(*Scanner)

type Scanner struct {
	index         map[string]Entity
	indexFilename string
}

func WithIndexFilename(filename string) Option {
	return func(s *Scanner) {
		s.indexFilename = filename
	}
}

func NewScanner(opts ...Option) *Scanner {
	index := make(map[string]Entity)
	scanner := &Scanner{
		index: index,
	}

	for _, fn := range opts {
		fn(scanner)
	}

	return scanner
}

func (s *Scanner) Scan(directories []string) error {
	file, err := os.Create(s.indexFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, directory := range directories {
		err = s.walk(file, directory)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scanner) walk(file *os.File, directory string) error {

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrInvalid) || errors.Is(err, fs.ErrPermission) || errors.Is(err, fs.ErrExist) || errors.Is(err, fs.ErrNotExist) || errors.Is(err, fs.ErrClosed) {
				fmt.Fprintln(os.Stderr, err.Error())
				return nil
			} else {
				return err
			}
		}

		if info.Mode().IsRegular() {
			hash := checksum(path)
			if hash != "" {
				size := info.Size()
				s.store(hash, size, path)
				fmt.Fprintf(file, "%040x %010d %s\n", hash, size, path)
			}
		}
		return nil
	})

	return err
}

func (s *Scanner) store(hash string, size int64, path string) {
	key := fmt.Sprintf("%v:%d", hash, size)
	entity, found := s.index[key]
	if !found {
		entity = Entity{
			Checksum: hash,
			FileSize: size,
			Files:    make([]string, 0),
		}
		//s.index[key] = entity
	}
	entity.Files = append(entity.Files, path)
	s.index[key] = entity
}

func (s *Scanner) PrintStats() {
	list := make([]Entity, 0)
	for _, value := range s.index {
		if len(value.Files) > 1 {
			list = append(list, value)
		}
	}

	sort.Slice(list, func(i, j int) bool {
		size1 := int64(len(list[i].Files)) * list[i].FileSize
		size2 := int64(len(list[j].Files)) * list[j].FileSize
		return size1 > size2
	})

	for _, entity := range list {
		count := len(entity.Files)
		size := uint64(count) * uint64(entity.FileSize)
		fmt.Printf("%d files occupy %v\n", count, humanize.Bytes(size))
		for _, file := range entity.Files {
			fmt.Printf("  %v\n", file)
		}
	}

	if len(list) == 0 {
		fmt.Println("No duplicates file found")
	}
}

func (s *Scanner) LoadIndex() error {
	file, err := os.Open(s.indexFilename)
	if err != nil {
		return fmt.Errorf("cannot load index: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		tokens := strings.SplitN(row, " ", 3)
		size, _ := strconv.ParseInt(tokens[1], 10, 64)

		s.store(tokens[0], size, tokens[2])
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error scanning file: %w", err)
	}
	return nil
}

func checksum(file string) string {
	in, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer in.Close()

	digest := xxhash.New()
	written, err := io.Copy(digest, in)
	if err != nil {
		return ""
	}
	if written == 0 {
		return ""
	}
	sum64 := digest.Sum64()

	checksum := strconv.FormatUint(sum64, 10)
	return checksum
}
