package parser

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrintableResult(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.3\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	printableResult := result.PrintableResult()

	assert.Equal(t,
		"{Canberra:6.1/6.1/6.1;Mogadishu:5.9/6.35/6.8;Odesa:10.2/10.2/10.2;Tirana:9.3/12.2/15.1;}",
		printableResult,
	)
}

func TestPrintableResultWithRepetitionOfSameStations(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.3\nMogadishu;6.8\nCanberra;6.1\nOdesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.3\nMogadishu;6.8\nCanberra;6.1\nOdesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.3\nMogadishu;6.8\nCanberra;6.1\nOdesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.3\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	printableResult := result.PrintableResult()

	assert.Equal(t,
		"{Canberra:6.1/6.1/6.1;Mogadishu:5.9/6.35/6.8;Odesa:10.2/10.2/10.2;Tirana:9.3/12.2/15.1;}",
		printableResult,
	)
}
