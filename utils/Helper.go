package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func Dump2Json(v any) string {
	bytes, err := json.MarshalIndent(v, " ", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func Dump2Gob(v any) []byte {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(v)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func Gob2Obj[T any](buffer []byte) T {
	reader := bytes.NewReader(buffer)
	dec := gob.NewDecoder(reader)
	var obj T
	err := dec.Decode(&obj)
	if err != nil {
		panic(err)
	}
	return obj
}
