package utils

import (
	"bufio"
	"os"
)

func ReadFileToList(sourceFilePath string) ([]string, error) {
	file, err := os.Open(sourceFilePath)

	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()
	return text, nil
}
