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
	startOffset int64
	size        int64
}

func SplitFile(file *os.File, fileSize int64, numParts int) ([]Chunk, error) {
	splitSize := fileSize / int64(numParts)
	buffer := make([]byte, maxLineLength)
	parts := make([]Chunk, 0, numParts)

	offset := int64(0)
	for part := 0; part < numParts; part++ {
		if part == numParts-1 {
			if offset < fileSize {
				parts = append(parts, Chunk{startOffset: offset, size: fileSize - offset})
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

		newlineIndex := bytes.LastIndexByte(chunk, '\n')
		if newlineIndex < 0 {
			return nil, fmt.Errorf("newline character not found at offset %d", offset+splitSize-maxLineLength)
		}

		remaining := len(chunk) - newlineIndex - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)
		parts = append(parts, Chunk{startOffset: offset, size: nextOffset - offset})
		offset = nextOffset
	}
	return parts, nil
}
