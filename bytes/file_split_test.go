package bytes

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"
)

func TestFileSplitInTwoParts(t *testing.T) {
	chunks, err := SplitFile("../fixture/10_weather_stations.txt", 2)
	assert.NoError(t, err)

	file, err := os.Open("../fixture/10_weather_stations.txt")
	assert.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	assert.Len(t, chunks, 2)
}

func TestFileSplitInTwoPartsWithEachPartContent(t *testing.T) {
	chunks, err := SplitFile("../fixture/10_weather_stations.txt", 2)
	assert.NoError(t, err)

	file, err := os.Open("../fixture/10_weather_stations.txt")
	assert.NoError(t, err)

	assert.Len(t, chunks, 2)

	content, _ := io.ReadAll(file)
	_ = file.Close()

	file, err = os.Open("../fixture/10_weather_stations.txt")
	assert.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	chunkContent := strings.Builder{}
	for _, chunk := range chunks {
		buffer := make([]byte, chunk.Size)
		_, _ = file.Read(buffer)
		chunkContent.WriteString(string(buffer))
	}

	assert.Equal(t, string(content), chunkContent.String())
}

func TestFileSplitInTwoPartsWithEachChunkTerminatedByNewLine(t *testing.T) {
	chunks, err := SplitFile("../fixture/10_weather_stations.txt", 2)
	assert.NoError(t, err)

	file, err := os.Open("../fixture/10_weather_stations.txt")
	assert.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	assert.Len(t, chunks, 2)

	_, _ = file.Seek(chunks[1].StartOffset-1, io.SeekStart)
	newLineBuffer := make([]byte, 1)
	_, _ = file.Read(newLineBuffer)

	assert.Equal(t, "\n", string(newLineBuffer))

	_, _ = file.Seek(chunks[0].Size+chunks[1].Size-1, io.SeekStart)
	newLineBuffer = make([]byte, 1)
	_, _ = file.Read(newLineBuffer)

	assert.Equal(t, "\n", string(newLineBuffer))
}

func TestFileSplitInThreePartsWithEachPartContent(t *testing.T) {
	chunks, err := SplitFile("../fixture/10_weather_stations.txt", 3)
	assert.NoError(t, err)

	file, err := os.Open("../fixture/10_weather_stations.txt")
	assert.NoError(t, err)

	assert.Len(t, chunks, 2)

	content, _ := io.ReadAll(file)
	_ = file.Close()

	file, err = os.Open("../fixture/10_weather_stations.txt")
	assert.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	chunkContent := strings.Builder{}
	for _, chunk := range chunks {
		buffer := make([]byte, chunk.Size)
		_, _ = file.Read(buffer)
		chunkContent.WriteString(string(buffer))
	}

	assert.Equal(t, string(content), chunkContent.String())
}
