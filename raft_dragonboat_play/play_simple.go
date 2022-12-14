package raft_dragonboat_play

import (
	"bufio"
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lni/dragonboat/v4"
	"github.com/lni/dragonboat/v4/config"
	"github.com/lni/dragonboat/v4/logger"
	sm "github.com/lni/dragonboat/v4/statemachine"

	"github.com/lni/goutils/syncutil"
)

const (
	exampleShardID uint64 = 128
)

var (
	// initial nodes count is fixed to three, their addresses are also fixed
	// these are the initial member nodes of the Raft cluster.
	addresses = []string{
		"localhost:63001",
		"localhost:63002",
		"localhost:63003",
	}
)

func Play_Simple() {
	replicaID := flag.Int("replicaid", 1, "ReplicaID to use")
	flag.Parse()
	if *replicaID != 1 && *replicaID != 2 && *replicaID != 3 {
		fmt.Fprintf(os.Stderr, "node id must be 1, 2 or 3 when address is not specified\n")
		os.Exit(1)
	}
	initialMembers := make(map[uint64]string)
	for idx, v := range addresses {
		// key is the ReplicaID, ReplicaID is not allowed to be 0
		// value is the raft address
		initialMembers[uint64(idx+1)] = v
	}
	nodeAddr := initialMembers[uint64(*replicaID)]
	fmt.Fprintf(os.Stdout, "node address: %s\n", nodeAddr)

	logger.GetLogger("raft").SetLevel(logger.ERROR)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)

	// config for raft node
	// See GoDoc for all available options
	rc := config.Config{
		// ShardID and ReplicaID of the raft node
		ReplicaID: uint64(*replicaID),
		ShardID:   exampleShardID,
		// In this example, we assume the end-to-end round trip time (RTT) between
		// NodeHost instances (on different machines, VMs or containers) are 200
		// millisecond, it is set in the RTTMillisecond field of the
		// config.NodeHostConfig instance below.
		// ElectionRTT is set to 10 in this example, it determines that the node
		// should start an election if there is no heartbeat from the leader for
		// 10 * RTT time intervals.
		ElectionRTT: 10,
		// HeartbeatRTT is set to 1 in this example, it determines that when the
		// node is a leader, it should broadcast heartbeat messages to its followers
		// every such 1 * RTT time interval.
		HeartbeatRTT: 1,
		CheckQuorum:  true,
		// SnapshotEntries determines how often should we take a snapshot of the
		// replicated state machine, it is set to 10 her which means a snapshot
		// will be captured for every 10 applied proposals (writes).
		// In your real world application, it should be set to much higher values
		// You need to determine a suitable value based on how much space you are
		// willing use on Raft Logs, how fast can you capture a snapshot of your
		// replicated state machine, how often such snapshot is going to be used
		// etc.
		SnapshotEntries: 10,
		// Once a snapshot is captured and saved, how many Raft entries already
		// covered by the new snapshot should be kept. This is useful when some
		// followers are just a little bit left behind, with such overhead Raft
		// entries, the leaders can send them regular entries rather than the full
		// snapshot image.
		CompactionOverhead: 5,
	}

	datadir := filepath.Join(
		"tmp",
		"raf-simple-data",
		fmt.Sprintf("node%d", *replicaID))

	// config for the nodehost
	// See GoDoc for all available options
	// by default, insecure transport is used, you can choose to use Mutual TLS
	// Authentication to authenticate both servers and clients. To use Mutual
	// TLS Authentication, set the MutualTLS field in NodeHostConfig to true, set
	// the CAFile, CertFile and KeyFile fields to point to the path of your CA
	// file, certificate and key files.
	nhc := config.NodeHostConfig{
		// WALDir is the directory to store the WAL of all Raft Logs. It is
		// recommended to use Enterprise SSDs with good fsync() performance
		// to get the best performance. A few SSDs we tested or known to work very
		// well
		// Recommended SATA SSDs -
		// Intel S3700, Intel S3710, Micron 500DC
		// Other SATA enterprise class SSDs with power loss protection
		// Recommended NVME SSDs -
		// Most enterprise NVME currently available on the market.
		// SSD to avoid -
		// Consumer class SSDs, no matter whether they are SATA or NVME based, as
		// they usually have very poor fsync() performance.
		//
		// You can use the pg_test_fsync tool shipped with PostgreSQL to test the
		// fsync performance of your WAL disk. It is recommended to use SSDs with
		// fsync latency of well below 1 millisecond.
		//
		// Note that this is only for storing the WAL of Raft Logs, it is size is
		// usually pretty small, 64GB per NodeHost is usually more than enough.
		//
		// If you just have one disk in your system, just set WALDir and NodeHostDir
		// to the same location.
		WALDir: datadir,
		// NodeHostDir is where everything else is stored.
		NodeHostDir: datadir,
		// RTTMillisecond is the average round trip time between NodeHosts (usually
		// on two machines/vms), it is in millisecond. Such RTT includes the
		// processing delays caused by NodeHosts, not just the network delay between
		// two NodeHost instances.
		RTTMillisecond: 200,
		// RaftAddress is used to identify the NodeHost instance
		RaftAddress: nodeAddr,
	}

	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		panic(err)
	}

	if err := nh.StartReplica(initialMembers, false, newSimpleStateMachine, rc); err != nil {
		fmt.Fprintf(os.Stderr, "failed to add cluster, %v\n", err)
		os.Exit(1)
	}

	raftStopper := syncutil.NewStopper()
	consoleStopper := syncutil.NewStopper()
	ch := make(chan string, 16)

	consoleStopper.RunWorker(func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			if s == "exit\n" {
				raftStopper.Stop()
				// no data will be lost/corrupted if nodehost.Stop() is not called
				nh.Close()
				return
			}
			ch <- s
		}
	})

	raftStopper.RunWorker(func() {
		// this goroutine makes a linearizable read every 10 second. it returns the
		// Count value maintained in IStateMachine. see datastore.go for details.
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				result, err := nh.SyncRead(ctx, exampleShardID, []byte{})
				cancel()
				if err == nil {
					var count uint64
					count = binary.LittleEndian.Uint64(result.([]byte))
					fmt.Fprintf(os.Stdout, "count: %d\n", count)
				}
			case <-raftStopper.ShouldStop():
				return
			}
		}
	})

	raftStopper.RunWorker(func() {
		// use a NO-OP client session here
		// check the example in godoc to see how to use a regular client session
		cs := nh.GetNoOPSession(exampleShardID)
		for {
			select {
			case v, ok := <-ch:
				if !ok {
					return
				}
				// remove the \n char
				msg := strings.Replace(v, "\n", "", 1)
				{
					// input is a regular message need to be proposed
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					// make a proposal to update the IStateMachine instance
					_, err := nh.SyncPropose(ctx, cs, []byte(msg))
					cancel()
					if err != nil {
						fmt.Fprintf(os.Stderr, "SyncPropose returned error %v\n", err)
					}
				}
			case <-raftStopper.ShouldStop():
				return
			}
		}
	})
	raftStopper.Wait()
}

// ExampleStateMachine is the IStateMachine implementation used in the
// helloworld example.
// See https://github.com/lni/dragonboat/blob/master/statemachine/rsm.go for
// more details of the IStateMachine interface.
type simpleStateMachine struct {
	ShardID   uint64
	ReplicaID uint64
	Count     uint64
}

// NewExampleStateMachine creates and return a new ExampleStateMachine object.
func newSimpleStateMachine(shardID uint64,
	replicaID uint64) sm.IStateMachine {
	return &simpleStateMachine{
		ShardID:   shardID,
		ReplicaID: replicaID,
		Count:     0,
	}
}

// Lookup performs local lookup on the ExampleStateMachine instance. In this example,
// we always return the Count value as a little endian binary encoded byte
// slice.
func (s *simpleStateMachine) Lookup(query interface{}) (interface{}, error) {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, s.Count)
	return result, nil
}

// Update updates the object using the specified committed raft entry.
func (s *simpleStateMachine) Update(e sm.Entry) (sm.Result, error) {
	// in this example, we print out the following hello world message for each
	// incoming update request. we also increase the counter by one to remember
	// how many updates we have applied
	s.Count++
	fmt.Printf("from ExampleStateMachine.Update(), msg: %s, count:%d\n",
		string(e.Cmd), s.Count)
	return sm.Result{Value: uint64(len(e.Cmd))}, nil
}

// SaveSnapshot saves the current IStateMachine state into a snapshot using the
// specified io.Writer object.
func (s *simpleStateMachine) SaveSnapshot(w io.Writer,
	fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	// as shown above, the only state that can be saved is the Count variable
	// there is no external file in this IStateMachine example, we thus leave
	// the fc untouched

	fmt.Printf("from ExampleStateMachine.SaveSnapshot(), count:%d\n", s.Count)

	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, s.Count)
	_, err := w.Write(data)
	return err
}

// RecoverFromSnapshot recovers the state using the provided snapshot.
func (s *simpleStateMachine) RecoverFromSnapshot(r io.Reader,
	files []sm.SnapshotFile,
	done <-chan struct{}) error {
	// restore the Count variable, that is the only state we maintain in this
	// example, the input files is expected to be empty
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	v := binary.LittleEndian.Uint64(data)
	s.Count = v

	fmt.Printf("from ExampleStateMachine.RecoverFromSnapshot(), count:%d\n", s.Count)
	return nil
}

// Close closes the IStateMachine instance. There is nothing for us to cleanup
// or release as this is a pure in memory data store. Note that the Close
// method is not guaranteed to be called as node can crash at any time.
func (s *simpleStateMachine) Close() error { return nil }
