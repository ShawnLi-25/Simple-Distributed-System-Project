package node

import (
	"fmt"
	"log"
	"net"
	msg "../Helper"
	"time"
)

// Listener is a type that implements the ListenMsg(), ListenJoinMsg() "method"
type Listener struct {
}

func HandleListenMsg(conn *net.UDPConn) {
	msgBuf := make([]byte, 128)

	n, msgAddr, err := conn.ReadFromUDP(msgBuf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listender: Recieve Msg from UDP client: %s", msgAddr)

	receivedMsg := msg.JSONToMsg([]byte(string(msgBuf[:n])))
	msg.PrintMsg(receivedMsg)

	switch receivedMsg.MessageType {
	case msg.HeartbeatMsg:
		fmt.Println("===Receive Heartbeat===")
		SendIntroduceMsg(conn, receivedMsg.NodeID)
	case msg.JoinMsg:
		fmt.Println("===Receive JoinMsg===")
	case msg.FailMsg:
		fmt.Println("===Receive FailMsg===")
	case msg.LeaveMsg:
		fmt.Println("===Receive LeaveMsg===")
	case msg.IntroduceMsg:
		fmt.Println("===Receive IntroduceMsg===")
	default:
		fmt.Println("Listener:Can't recognize the msg")
	}
}

func (l *Listener) NodeListen() {
	go ListenHeartbeat()
	go ListenMsg()
}

//Listen to Heartbeat and Check timeout
func ListenHeartbeat() {
	fmt.Println("HBListener:Initialize heartbeat listener...")
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+msg.HeartbeatPort)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HBListener:Listen Heartbeat on port %s\n", msg.ConnPort)
	defer ln.Close()
	hbBuf := make([]byte, 128)
	ln.SetReadDeadline(time.Now().Add(msg.Timeout))
	for{
		n, msgAddr, err := ln.ReadFromUDP(hbBuf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Listener:Recieve Heartbeat from UDP client: %s", msgAddr)
		if n > 0 {
			//No delay, refresh deadline
			ln.SetReadDeadline(time.Now().Add(msg.Timeout))
			receivedMsg := msg.JSONToMsg([]byte(string(hbBuf[:n])))
			msg.PrintMsg(receivedMsg)
		}
		
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			//Timeout error
			fmt.Printf("HBListener: Client %s Timeout!\n", msgAddr)
			//TODO Send timeout msg
		}
	}
}

//Listen to Msg (Introducer-only)
func ListenMsg() {
	fmt.Println("MSGListener:Initialize new listener...")
	udpAddr, err := net.ResolveUDPAddr(msg.ConnType, ":"+msg.ConnPort)
	if err != nil {
		log.Fatal(err)
	}

	ln, err := net.ListenUDP(msg.ConnType, udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MSGListener:Listen on port %s\n", msg.ConnPort)
	defer ln.Close()
	for {
		HandleListenMsg(ln)
	}
	
}
