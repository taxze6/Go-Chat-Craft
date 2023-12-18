package models

import (
	"GoChatCraft/global"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Message struct {
	Model
	FormId   int64  `json:"userId"`
	TargetId int64  `json:"targetId"`
	Type     int    `json:"type"`
	Media    int    `json:"media"`
	Content  string `json:"content"`
	Pic      string `json:"pic"`
	Url      string `json:"url"`
	Desc     string
	Amount   int
}

func (m *Message) MsgTableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	Addr      string
	DataQueue chan []byte //消息
	GroupSets set.Interface
}

var clientMap map[int64]*Node = make(map[int64]*Node) //Mapping relationship table (the key of the map is the userId, and the value is the Node, a global map shared by all coroutines)
var rwlocker sync.RWMutex                             //Read-write lock is needed to ensure thread safety when binding a Node.

func Chat(w http.ResponseWriter, r *http.Request, Id string) {
	//query := r.URL.Query()
	//Id := query.Get("userId")
	userId, err := strconv.ParseInt(Id, 10, 64)
	if err != nil {
		zap.S().Info("Type conversion failed", err)
		return
	}
	//Upgrade to socket
	var isvalida = true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//Get socket connection, construct message node
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// Bind userId and Node together
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()

	//Service sends message
	go sendProc(node)

	//Service receives message
	go recProc(node)

	//test
	sendMsg(userId, []byte("Welcome to the chat."))
}

// sendProc retrieves information from the node and writes it into the WebSocket.
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Info("Failed to write the message", err)
				return
			}
			fmt.Println("The data has been successfully sent through the socket.")
		}
	}
}

// recProc retrieves the message body from the WebSocket,
// then parses it, performs message type identification, and finally sends the message to the destination user's node.
func recProc(node *Node) {
	for {
		//Retrieve information
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("Failed to read the message.", err)
			return
		}
		//Put the message body into the global channel.
		brodMsg(data)
	}
}

var upSendChan chan []byte = make(chan []byte, 1024)

func brodMsg(data []byte) {
	upSendChan <- data
}

func init() {
	go UdpSendProc()
	go UdpRecProc()
}

// The UdpSendProc completes UDP data sending by connecting to the UDP server and writing the message body from the global channel to the UDP server.
func UdpSendProc() {
	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
		Zone: "",
	})
	if err != nil {
		zap.S().Info("Failed to dial UDP port.")
		return
	}
	defer udpConn.Close()
	for {
		select {
		case data := <-upSendChan:
			_, err := udpConn.Write(data)
			if err != nil {
				zap.S().Info("Failed to write UDP message.", err)
				return
			}
			fmt.Println("The data has been successfully sent to the UDP server:", string(data))
		}
	}
}

// UdpRecProc completes the reception of UDP data, starts the UDP service, and retrieves the messages written by UDP clients.
func UdpRecProc() {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
	})
	if err != nil {
		zap.S().Info("Failed to listen on UDP port.", err)
		return
	}
	defer udpConn.Close()
	for {
		var buf [1024]byte
		n, err := udpConn.Read(buf[0:])
		if err != nil {
			zap.S().Info("Failed to read UDP data.", err)
			return
		}
		fmt.Println("UDP server receives UDP data:", buf[0:n])
		dispatch(buf[0:n])
	}
}

// Dispatch: Parsing the message and determining the chat type.
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		zap.S().Info("Failed to parse the message.", err)
		return
	}

	fmt.Println("Parse the data:", "msg.FormId", msg.FormId, "targetId:", msg.TargetId, "type:", msg.Type)
	switch msg.Type {
	case 1:
		sendMsg(msg.TargetId, data)
	case 2:
		sendGroupMsg(uint(msg.FormId), uint(msg.TargetId), data)
	}
}

// Sending a message to the user in a private chat.
func sendMsg(id int64, msg []byte) {
	rwlocker.Lock()
	node, ok := clientMap[id]
	rwlocker.Unlock()
	if !ok {
		zap.S().Info("There is no corresponding node for the userID.")
		return
	}
	zap.S().Info("targetID:", id, "node:", node)
	if ok {
		node.DataQueue <- msg
	}
}

func sendGroupMsg(formId, target uint, data []byte) (int, error) {
	//Get all users in the group, and then send a message to each user except yourself.
	userIDs, err := FindUsersId(target)
	if err != nil {
		return -1, err
	}
	for _, userId := range *userIDs {
		//Do not forward the message to the member who is currently sending the message.
		if formId != userId {
			//Multiple calls to the individual chat function have turned the group chat into multiple individual chats.
			sendMsgAndSave(int64(userId), data)
		}
	}
	return 0, nil
}

func sendMsgAndSave(userId int64, msg []byte) {
	rwlocker.RLock()              //Ensure thread safety by locking.
	node, ok := clientMap[userId] //Whether the other party is online
	rwlocker.RUnlock()
	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.FormId))
	if ok {
		//If the current user is online, forward the message to the current user's WebSocket connection and then store it.
		node.DataQueue <- msg
	}
	//Concatenate userIdStr and targetIdStr to create a unique key.
	var key string
	if userId > jsonMsg.FormId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	//
	res, err := global.RedisDB.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	//Write the chat records into Redis cache.
	score := float64(cap(res)) + 1
	ress, e := global.RedisDB.ZAdd(ctx, key, &redis.Z{Score: score, Member: msg}).Result()
	if e != nil {
		zap.S().Info(e)
		return
	}
	fmt.Println(ress)
}

// isRev is a boolean parameter used to indicate whether to retrieve chat records from the cache in reverse order (from largest to smallest).
// If isRev is true, the ZRange function is used to retrieve records in ascending order.
// If isRev is false, the ZRevRange function is used to retrieve records in reverse order.
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))

	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	var rels []string
	var err error
	if isRev {
		rels, err = global.RedisDB.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = global.RedisDB.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		//The message could not be found.
		zap.S().Info(err)
		fmt.Println(err)
	}
	return rels
}
