package bytes

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

const (
	maxStationLength     = 100
	maxTemperatureLength = 5
	separatorLength      = 1
	newLineLength        = 1
	maxLineLength        = maxStationLength + maxTemperatureLength + separatorLength + newLineLength
)

type Chunk struct {
	StartOffset int64
	Size        int64
}

func SplitFile(fileName string, numParts int) ([]Chunk, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileStat.Size()
	splitSize := fileSize / int64(numParts)
	buffer := make([]byte, maxLineLength)
	parts := make([]Chunk, 0, numParts)

	offset := int64(0)
	for part := 0; part < numParts; part++ {
		if part == numParts-1 {
			if offset < fileSize {
				parts = append(parts, Chunk{StartOffset: offset, Size: fileSize - offset})
			}
			break
		}

		seekOffset := max(offset+splitSize-maxLineLength, 0)
		_, err := file.Seek(seekOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		n, _ := io.ReadFull(file, buffer)
		chunk := buffer[:n]

		lastNewlineIndex := bytes.LastIndexByte(chunk, '\n')
		if lastNewlineIndex < 0 {
			return nil, fmt.Errorf("newline character not found at offset %d", offset+splitSize-maxLineLength)
		}

		remaining := len(chunk) - lastNewlineIndex - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)
		parts = append(parts, Chunk{StartOffset: offset, Size: nextOffset - offset})
		offset = nextOffset
	}
	return parts, nil
}
