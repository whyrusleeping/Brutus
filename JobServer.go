package main

import (
	"net"
	"fmt"
	"log"
	"math/big"
	"bufio"
	"bytes"
	"math/rand"
	"strconv"
	"sync"
	"time"
	"encoding/binary"
	"io"
)

type Record struct {
	key string
	val int
	owner string
}

//The cluster manager
type JobServer struct {
	best Record
	completed int64
	compLock sync.Mutex
	jobs []byte
	jobLock sync.Mutex
	blocksize *big.Int
	startTime time.Time
}

//Start up the TCP server, listening on the specified port
func (js *JobServer) Start(port int) error {
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d",port))
	if err != nil {
		return err
	}
	list, err := net.ListenTCP("tcp",laddr)
	if err != nil {
		return err
	}

	//Record start time
	js.startTime = time.Now()
	for {
		ncon, err := list.Accept()
		if err == nil {
			//Asynchrounously handle the connection
			go js.HandleConnection(ncon)
		}
	}
	return nil
}

func (js *JobServer) HandleConnection(c net.Conn) {
	read := bufio.NewReader(c)
	mess, err := read.ReadBytes('|')
	if err != nil {
		log.Println("Client sent malformed message.")
		return
	}
	parts := bytes.Split(mess, []byte(";"))
	typ := string(parts[0])
	switch typ {
		case "ask":
			js.SendJob(c)
		case "update":
			if len(parts) < 3 {
				log.Println("Improper update message")
				return
			}
			js.ReceiveUpdate(parts[1:])
		case "new":
			//js.NewNode(c)
		case "checkin":
			//Error check here
			n,_ := strconv.Atoi(string(parts[1]))
			js.jobs[n] = 2
			js.compLock.Lock()
			js.completed += js.blocksize.Int64()
			js.compLock.Unlock()
	}
}

//Parse updates from a client
func (js *JobServer) ReceiveUpdate(info [][]byte) {
	fmt.Printf("%s - %s - %s\n", info[1], info[2], info[3])
}

//This code is specific to the task
func (js *JobServer) JobForIndex(i int64) string {
	start := big.NewInt(int64(i))
	start.Mul(start, js.blocksize)
	end := big.NewInt(js.blocksize.Int64())
	end.Add(end, start)
	return fmt.Sprintf("%d;%s;%s|", i, start.String(), end.String())
}

//Insert logic for scheduling here
func (js *JobServer) SendJob(c net.Conn) {
	i := int64(0)

	js.jobLock.Lock()
	//TODO: What happens when there are no blocks left?
	for ; js.jobs[i] != 0; i = rand.Int63n(int64(len(js.jobs))) {}
	js.jobs[i] = 1
	js.jobLock.Unlock()

	//Generate Job for selected index
	job := js.JobForIndex(i)
	//Send it to the Node
	c.Write([]byte(job))
}

//TOOLS
func WriteByteArray(w io.Writer, s []byte) {
	binary.Write(w, binary.LittleEndian, uint16(len(s)))
	w.Write(s)
}


func main() {
	js := new(JobServer)
	js.Start(9000)
}
