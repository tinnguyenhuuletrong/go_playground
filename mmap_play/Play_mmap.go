package mmap_play

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/edsrzf/mmap-go"
)

type MMapReader struct {
	mmap   mmap.MMap
	cursor int
}

func newMMapRead(m mmap.MMap) *MMapReader {
	return &MMapReader{
		mmap:   m,
		cursor: 0,
	}
}

func (ins *MMapReader) Read(p []byte) (n int, err error) {
	dataLen := len(ins.mmap)
	pLen := len(p)

	if ins.cursor >= dataLen {
		return 0, nil
	}

	nByte := dataLen
	if nByte > pLen {
		nByte = pLen
	}

	for i := 0; i < nByte; i++ {
		p[i] = ins.mmap[ins.cursor+i]
	}
	ins.cursor += nByte

	return nByte, nil
}

func Play_mmap() {

	f, err := os.OpenFile("./tmp/mmap_file", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	mmap, _ := mmap.Map(f, mmap.RDWR, 0)
	defer mmap.Unmap()

	wIo := newMMapRead(mmap)
	reader := bufio.NewReader(wIo)

	println("mmap setup: ./tmp/mmap_file ")
	println(`Try to run: echo "$(date)" > ./tmp/mmap_file `)

	ticker := time.NewTicker(1 * time.Second)

	for _ = range ticker.C {
		wIo.cursor = 0
		reader.Reset(wIo)
		tmp, _ := reader.ReadString(0)
		println(time.Now().Local().String(), tmp)
	}
}
