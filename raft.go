package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"net/http"
	"net/rpc"
	"strings"
	"time"
)

var (
	globalMap = make(map[uint64]int)
)

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func WriteEntry(rf *Raft, clientName string, msg string) int {
	var index int
	qMsgStr := clientName + "$" + msg

	if rf.state == Leader {
		rf.writeChan <- qMsgStr
		index = len(rf.log) + 1
	} else {
		index = -1
	}

	return index

}

func ReadEntry(rf *Raft, idx int) string {
	var str string
	maxIdx := rf.getLastIndex()
	if rf.state != Leader {
		return "error: redirect client call to Leader"
	} else {
		if (idx == 0) {
			str = "error: index has to be Greater than 0"
		} else if idx <= maxIdx {
			str = rf.log[idx-1].LogCMD
		} else {
			str = "error: no log record exists"
		}
		return str
	}
}

func GetEntries(rf *Raft, clientName string) string {
	var str string
	maxIdx := rf.getLastIndex()
	if rf.state != Leader {
		return "error: redirect client call to Leader"
	} else {
		var strArr []string
		for idx := 0; idx < maxIdx; idx++ {
			lEntry := rf.log[idx]
			if lEntry.LogClient == clientName {
				strArr = append(strArr, lEntry.LogCMD)
			}
		}

		str = stringArraySerialise(strArr)
		return str
	}
}

func GetAllEntries(rf *Raft) string {
	var str string
	maxIdx := rf.getLastIndex()
	if rf.state != Leader {
		return "error: redirect client call to Leader"
	} else {
		var strArr []string
		for idx := 0; idx < maxIdx; idx++ {
			lEntry := rf.log[idx]
			strArr = append(strArr, lEntry.LogCMD)
		}

		str = stringArraySerialise(strArr)
		return str
	}
}

func GetAllEntriesArray(rf *Raft) []string {
	maxIdx := rf.getLastIndex()
	var strArr []string
	if rf.state != Leader {
		return strArr
	} else {

		for idx := 0; idx < maxIdx; idx++ {
			lEntry := rf.log[idx]
			strArr = append(strArr, lEntry.LogCMD)
		}
		return strArr
	}
}

func stringArraySerialise(arr []string) string {
	var res string
	res = " [ "
	for _, v := range arr {
		res += "{ " + v + " } "
	}
	res += " ] "
	return res
}

type node struct {
	connect bool
	address string
}

func callMeDaddy() {
	fmt.Println("Hi Daddy")
}

// New node
func newNode(address string) *node {
	node := &node{}
	node.address = address
	return node
}

// State def
type State int

// status of node
const (
	Follower State = iota + 1
	Candidate
	Leader
)

func whatState(st State) string {
	switch st {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	}
	return ""
}

// LogEntry struct
type LogEntry struct {
	LogTerm   int
	LogIndex  int
	LogCMD    string
	LogClient string
}

// Raft Node
type Raft struct {
	me int

	nodes map[int]*node

	state       State
	currentTerm int
	votedFor    int
	voteCount   int

	// Collection of log entries
	log []LogEntry

	// The largest index submitted
	commitIndex int
	// Maximum index applied to the state machine
	lastApplied int

	// Save the index of the next entry that needs to be sent to each node
	nextIndex []int
	// Save the highest index of the log that has been copied to each node
	matchIndex []int

	// channels
	heartbeatC chan bool
	toLeaderC  chan bool

	// writeChannel
	writeChan chan string
}

// RequestVote rpc method
func (rf *Raft) RequestVote(args VoteArgs, reply *VoteReply) error {

	if args.Term < rf.currentTerm {
		reply.Term = rf.currentTerm
		reply.VoteGranted = false
		return nil
	}

	if rf.votedFor == -1 {
		rf.currentTerm = args.Term
		rf.votedFor = args.CandidateID
		reply.Term = rf.currentTerm
		reply.VoteGranted = true
	}

	return nil
}

// Heartbeat rpc method
func (rf *Raft) Heartbeat(args HeartbeatArgs, reply *HeartbeatReply) error {

	// in case leader Node is smaller than current node term
	if args.Term < rf.currentTerm {
		reply.Success = false
		reply.Term = rf.currentTerm
		return nil
	}

	// If only heartbeat
	rf.heartbeatC <- true
	if len(args.Entries) == 0 {
		reply.Success = true
		reply.Term = rf.currentTerm
		return nil
	}

	// If there is entries
	// leader Maintained LogIndex Greater than current Follower of LogIndex
	// Represents the current Follower Lost contact, so Follower To inform Leader It currently
	// The maximum index of the, so that the next heartbeat Leader will return
	if args.PrevLogIndex > rf.getLastIndex() {
		reply.Success = false
		reply.Term = rf.currentTerm
		reply.NextIndex = rf.getLastIndex() + 1
		return nil
	}
	fmt.Println("To be appended Entries ", args.Entries)
	for i := 0; i < len(args.Entries); i++ {
		if _, found := globalMap[hash(args.Entries[i].LogCMD)]; found {
			continue
		}
		rf.log = append(rf.log, args.Entries[i])
		globalMap[hash(args.Entries[i].LogCMD)] = 1
		rf.commitIndex = rf.getLastIndex()
		reply.Success = true
		reply.Term = rf.currentTerm
		reply.NextIndex = rf.getLastIndex() + 1
	}

	return nil
}

func (rf *Raft) rpc(port string) {
	rpc.Register(rf)
	rpc.HandleHTTP()
	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal("listen error: ", err)
		}
	}()
}

func (rf *Raft) start() {
	rf.state = Follower
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.heartbeatC = make(chan bool)
	rf.toLeaderC = make(chan bool)

	rf.writeChan = make(chan string)

	go func() {

		rand.Seed(time.Now().UnixNano())

		for {
			switch rf.state {
			case Follower:
				select {
				case <-rf.heartbeatC:
					log.Printf("follower-%d received heartbeat\n", rf.me)
				case <-time.After(time.Duration(rand.Intn(500-300)+300) * time.Millisecond):
					log.Printf("follower-%d timeout\n", rf.me)
					rf.state = Candidate
				}
			case Candidate:
				fmt.Printf("Node: %d, I'm candidate\n", rf.me)
				rf.currentTerm++
				rf.votedFor = rf.me
				rf.voteCount = 1
				go rf.broadcastRequestVote()

				select {
				case <-time.After(time.Duration(rand.Intn(500)+500) * time.Millisecond):
					rf.state = Follower
				case <-rf.toLeaderC:
					fmt.Printf("Node: %d, I'm leader\n", rf.me)
					rf.state = Leader

					// Initialize nextIndex and matchIndex of peers
					rf.nextIndex = make([]int, len(rf.nodes))
					rf.matchIndex = make([]int, len(rf.nodes))
					for i := range rf.nodes {
						rf.nextIndex[i] = 1
						rf.matchIndex[i] = 0
					}

					go func() {
						for {
							i := len(rf.log) + 1
							msg, ok := <-rf.writeChan
							fmt.Println("msg is " + msg)
							if ok {
								dt := time.Now()
								timeStr := dt.Format("01-02-2006 15:04:05")
								clientMsgArr := strings.Split(msg, "$")
								clientName, clientMsg := clientMsgArr[0], clientMsgArr[1]
								clientString := fmt.Sprintf("time : %s | clientName : %s | index : %d | msg : %s", timeStr, clientName, i, clientMsg)
								rf.log = append(rf.log, LogEntry{rf.currentTerm, i, clientString, clientName})
								globalMap[hash(clientString)] = 1
							} else {
								fmt.Println("no msg to log")
							}
							time.Sleep(500 * time.Millisecond)
						}
					}()
				}
			case Leader:
				rf.broadcastHeartbeat()
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

type VoteArgs struct {
	Term        int
	CandidateID int
}

type VoteReply struct {
	//Current term number so that candidates can update their term number
	Term int
	//True when the candidate won this ballot
	VoteGranted bool
}

func (rf *Raft) broadcastRequestVote() {
	var args = VoteArgs{
		Term:        rf.currentTerm,
		CandidateID: rf.me,
	}

	for i := range rf.nodes {
		go func(i int) {
			var reply VoteReply
			rf.sendRequestVote(i, args, &reply)
		}(i)
	}
}

func (rf *Raft) sendRequestVote(serverID int, args VoteArgs, reply *VoteReply) {
	client, err := rpc.DialHTTP("tcp", rf.nodes[serverID].address)

	if err != nil {
		return
	}

	defer client.Close()
	client.Call("Raft.RequestVote", args, reply)

	// The current candidate node is invalid
	if reply.Term > rf.currentTerm {
		rf.currentTerm = reply.Term
		rf.state = Follower
		rf.votedFor = -1
		return
	}

	if reply.VoteGranted {
		rf.voteCount++
	}

	majority := len(rf.nodes)/2 + 1

	if rf.voteCount >= majority {
		rf.toLeaderC <- true
	}
}

type HeartbeatArgs struct {
	Term     int
	LeaderID int

	// Index before the new log
	PrevLogIndex int
	// PrevLogIndex Term number
	PrevLogTerm int
	// The log entry to be stored (indicating that it is empty at the time of the heartbeat)
	Entries []LogEntry
	// Leader has committed index value
	LeaderCommit int
}

type HeartbeatReply struct {
	Success bool
	Term    int

	// If the Follower Index is less than the Leader Index, it will tell the Leader the index position to start sending next time
	NextIndex int
}

func (rf *Raft) broadcastHeartbeat() {
	for i := range rf.nodes {
		addr := rf.nodes[i].address
		log.Printf("addr %v , lastLogCMD %s", addr, rf.getLastCMD())

		var args HeartbeatArgs
		args.Term = rf.currentTerm
		args.LeaderID = rf.me
		args.LeaderCommit = rf.commitIndex

		// Calculate preLogIndex, preLogTerm
		// Extract the entry after preLogIndex-baseIndex and send it to follower
		prevLogIndex := rf.nextIndex[i] - 1
		if rf.getLastIndex() > prevLogIndex {
			args.PrevLogIndex = prevLogIndex
			args.PrevLogTerm = rf.log[prevLogIndex].LogTerm
			args.Entries = rf.log[prevLogIndex:]
			log.Printf("send entries: %v\n", args.Entries)
		}

		go func(i int, args HeartbeatArgs) {
			var reply HeartbeatReply
			rf.sendHeartbeat(i, args, &reply)
		}(i, args)
	}
}

func (rf *Raft) sendHeartbeat(serverID int, args HeartbeatArgs, reply *HeartbeatReply) {
	client, err := rpc.DialHTTP("tcp", rf.nodes[serverID].address)

	if err != nil {
		return
	}

	defer client.Close()

	client.Call("Raft.Heartbeat", args, reply)

	// If the leader node lags behind the follower node
	if reply.Success {
		if reply.NextIndex > 0 {
			rf.nextIndex[serverID] = reply.NextIndex
			rf.matchIndex[serverID] = rf.nextIndex[serverID] - 1
		}
	} else {
		// If the leader's term is less than the follower's term, you need to change the leader to a follower and re-elect
		if reply.Term > rf.currentTerm {
			rf.currentTerm = reply.Term
			rf.state = Follower
			rf.votedFor = -1
			return
		}
	}
}

func (rf *Raft) getLastCMD() string {
	rlen := len(rf.log)
	if rlen == 0 {
		return ""
	}
	return rf.log[rlen-1].LogCMD
}

func (rf *Raft) getLastIndex() int {
	rlen := len(rf.log)
	if rlen == 0 {
		return 0
	}
	return rf.log[rlen-1].LogIndex
}

func (rf *Raft) getLastTerm() int {
	rlen := len(rf.log)
	if rlen == 0 {
		return 0
	}
	return rf.log[rlen-1].LogTerm
}
