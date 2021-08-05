package fs

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"
)

func SafeCopy(src, dst string) (err error) {
	dest := dst
	num := 0
	for {
		if FileExist(dest) {
			num++
			dest = fmt.Sprintf("%s/%s(%d)%s", path.Dir(dst), strings.TrimSuffix(filepath.Base(dst), path.Ext(dst)), num, path.Ext(dst))
			logrus.Info(dest)
		} else {
			logrus.Info("here ", dest)
			break
		}
	}

	return Copy(src, dest)
}

// Copy ...
func Copy(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

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
		lines = append(lines, line)
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
	}
	return
}

// ReadBytes ...
func ReadBytes(file string) (bytes []byte, err error) {
	fd, err := os.Open(file)
	if err != nil {
		return
	}
	defer fd.Close()

	reader := bufio.NewReader(fd)
	buf := make([]byte, 1024*1024*16)

	for {
		size, err := reader.Read(buf)
		bytes = append(bytes, buf[:size]...)
		if errors.Is(err, io.EOF) || (err == nil && size == 0) {
			return bytes, nil
		}
		if err != nil {
			return bytes, err
		}
	}
}

// Save ...
func Save(data []byte, fullpath string) (err error) {
	dir := filepath.Dir(fullpath)
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		logrus.WithError(err).WithField("fullpath", fullpath).Error()
		return
	}
	var f *os.File
	if f, err = os.Create(fullpath); err != nil {
		logrus.WithError(err).WithField("fullpath", fullpath).Error()
		return
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		logrus.WithError(err).WithField("fullpath", fullpath).Error()
		return
	}
	return nil
}

// Stat ...
func Stat(fullpath string) (info os.FileInfo, err error) {
	if f, err := os.Open(fullpath); err != nil {
		logrus.WithError(err).WithField("fullpath", fullpath).Error()
		return nil, err
	} else {
		return f.Stat()
	}
}

// MIME ...
func MIME(fullpath string) (string, error) {
	mime, err := mimetype.DetectFile(fullpath)
	if err != nil {
		return ``, err
	}
	return mime.String(), err
}

// MD5 ...
func MD5(fullpath string) (string, error) {
	f, err := os.Open(fullpath)
	if err != nil {
		return ``, err
	}
	defer f.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return ``, err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
