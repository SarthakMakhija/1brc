package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test1BrcStatistics(t *testing.T) {
	outputBuffer := &bytes.Buffer{}
	outputDevice = outputBuffer

	print1brcStatistics("../fixture/10_rows.txt")

	assert.Equal(t,
		"{Canberra:7/7/7;Halifax:10/10/10;Mogadishu:6/48/90;Odesa:6/27/66;Tirana:3/6/9;}",
		string(outputBuffer.Bytes()),
	)
}
