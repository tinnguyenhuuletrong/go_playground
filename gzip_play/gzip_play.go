package gzip_play

import (
	"bytes"
	"compress/gzip"
	"log"
)

func Gzip_Play() {

	msg := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
	log.Println("Msg:", msg)

	compressBuf := new(bytes.Buffer)
	w := gzip.NewWriter(compressBuf)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Panic(err)
	}
	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	// send <-> recv
	data := compressBuf.Bytes()[:]
	log.Println("compressed.", "len:", len(data), "data:", data)

	decompressBuf := new(bytes.Buffer)
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

	_, err = decompressBuf.ReadFrom(r)
	if err != nil {
		log.Panic(err)
	}

	log.Println("decompressed:", decompressBuf.String())
}
