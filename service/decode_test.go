package service_test

import (
	"fmt"
	"testing"

	"github.com/martengine/reception/service"
)

const data = `{
	"items": "[
		\"id\",
		\"title\",
		\"description\",
		\"image\": {
			\"src\",
			\"thumbnail\",
			\"price\"
		}
	]",
	"data",
	"go"
}`

func BenchmarkDecode(b *testing.B) {
	dataBytes := []byte(data)
	for i := 0; i < b.N; i++ {
		service.Decode(dataBytes)
	}
}

func TestDecode(t *testing.T) {
	dataBytes := []byte(data)
	result, err := service.Decode(dataBytes)
	fmt.Println(result, err)
}
