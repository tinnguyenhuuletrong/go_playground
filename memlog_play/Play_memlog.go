package memlog_play

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/embano1/memlog"

	"github.com/tidwall/wal"
)

func stringGenerator(ml *memlog.Log, prefix string, count uint) {
	ctx := context.Background()
	go func() {
		for i := uint(0); i < count; i++ {
			data := fmt.Sprintf("%s-%d", prefix, i)
			ml.Write(ctx, []byte(data))
			time.Sleep(10 * time.Millisecond)
		}
	}()

}

func scanByOffset(ml *memlog.Log) {
	var offset memlog.Offset
	batchBuf := make([]memlog.Record, 10)
	ctx := context.Background()

	for {
		earliest, latest := ml.Range(ctx)
		fmt.Printf("offset (earliest,latest) %d,%d - cursor: %d \n", earliest, latest, offset)

		// scan data
		if offset < latest {
			num, err := ml.ReadBatch(ctx, offset, batchBuf)
			if err != nil {
				log.Panic(err)
			}
			offset += memlog.Offset(num)

			for _, v := range batchBuf {
				fmt.Printf("\t%d - %s\n", v.Metadata.Offset, v.Data)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func stream(ml *memlog.Log) {
	ctx := context.Background()

	stream := ml.Stream(ctx, 0)

	// write walLog
	walLog, err := wal.Open("tmp/wal", &wal.Options{
		NoSync:           false,    // Fsync after every write
		SegmentSize:      2048,     // 1 MB log segment files.
		LogFormat:        wal.JSON, // Binary format is small and fast.
		SegmentCacheSize: 2,        // Number of cached in-memory segments
		NoCopy:           false,    // Make a new copy of data for every Read call.
		DirPerms:         0750,     // Permissions for the created directories
		FilePerms:        0640,     // Permissions for the created data files
	})

	if err != nil {
		log.Panic(err)
	}

	for {
		if v, ok := stream.Next(); ok {
			fmt.Printf("\t%d - %s\n", v.Metadata.Offset, v.Data)
			index, err := walLog.LastIndex()
			if err != nil {
				log.Panic(err)
			}
			walLog.Write(index+1, v.Data)
		}
	}
}

func Play_MemLog() {
	ctx := context.Background()
	ml, err := memlog.New(ctx)
	if err != nil {
		log.Panicf("create log: %v", err)
	}

	go func() {
		for i := 0; i < 2; i++ {
			fmt.Printf("----------round %d--------\n", i)
			stringGenerator(ml, "itm", 100)
			time.Sleep(10 * time.Second)
		}
	}()

	// go scanByOffset(ml)
	go stream(ml)

	// Wait for Ctrl + C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
