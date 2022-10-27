package read

import (
	"bytes"
	"io"
	"os"
)

// slow
func ReadFileAll(filename string) (content string, err error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// fast
func ReadFilePointer(filename string) (content string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var chunk []byte
	buf := make([]byte, 2048)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		} else if n == 0 || err == io.EOF {
			break
		}
		chunk = append(chunk, buf[:n]...)
	}
	return string(chunk), nil
}

func ReadLineToChannel(filename string, ch chan []byte) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var chunk []byte
	buf := make([]byte, 2048)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		} else if n == 0 || err == io.EOF {
			if len(chunk) != 0 {
				ch <- chunk
			}
			break
		}
		chunk = append(chunk, buf[:n]...)

		lines := bytes.Split(chunk, []byte("\n"))
		if len(lines) <= 1 {
			continue
		}

		lastLineLen := len(lines[len(lines)-1])
		pos := len(chunk) - lastLineLen
		chunk = chunk[pos:]

		for _, val := range lines[:len(lines)-1] {
			ch <- val
		}
	}
	return nil
}
