package service_test

import (
	"testing"

	"github.com/martengine/reception/service"
)

const data = `{
	"items": [
		"id",
		"title",
		"description",
		"image": {
			"src",
			"thumbnail",
			"price"
		}
	],
	"data",
	"go"
}`

func BenchmarkDecode(b *testing.B) {
	dataBytes := []byte(data)
	for i := 0; i < b.N; i++ {
		service.Decode(dataBytes)
	}
}
