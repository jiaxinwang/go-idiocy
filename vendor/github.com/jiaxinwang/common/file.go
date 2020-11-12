package common

import (
	"bufio"
	"io"
	"os"
)

// FileExist ...
func FileExist(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

// IsDir ...
func IsDir(name string) bool {
	info, err := os.Stat(name)
	return err == nil && info.IsDir()
}

// Readlines ...
func Readlines(file string) (lines []string, err error) {
	fd, err := os.Open(file)
	if err != nil {
		return
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
		lines = append(lines, line)
	}
	return
}
