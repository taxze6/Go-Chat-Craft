package models

import (
	"GoChatCraft/global"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Message struct {
	Model
	MsgId                string     `json:"msgId"`
	FormId               int64      `json:"userId"`
	TargetId             int64      `json:"targetId"`
	Type                 int        `json:"type"`
	ContentType          int        `json:"contentType"`
	Status               int        `json:"status"`
	Content              string     `json:"content"`
	MessageSenderName    string     `json:"messageSenderName"`
	MessageSenderFaceUrl string     `json:"messageSenderFaceUrl"`
	Pic                  string     `json:"pic"`
	Url                  string     `json:"url"`
	Image                ImageModel `json:"image"`
	Sound                SoundModel `json:"sound"`
	QuoteMessage         *Message   `json:"quoteMessage"`
	Desc                 string
	Amount               int
}

type ImageModel struct {
	ImageUrl    string  `json:"imageUrl"`
	ImageWidth  float64 `json:"imageWidth"`
	ImageHeight float64 `json:"imageHeight"`
	ImageSize   float64 `json:"fileSize"`
}

type SoundModel struct {
	SourceUrl string `json:"sourceUrl"`

	SoundPath string `json:"soundPath"`

	DataSize float64 `json:"dataSize"`

	Duration int `json:"duration"`
}

// CustomError is a custom error type.
type CustomError struct {
	message string
}

// Error implements the error interface for CustomError.
func (e *CustomError) Error() string {
	return e.message
}

const (
	Sending   = 1
	Succeeded = 2
	Failed    = 3
	Deleted   = 4
)

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
	//sendMsgAndSave(userId, []byte("{\"msgId\":\"11111\",\"userId\":7,\"targetId\":1,\"type\":101,\"contentType\":101,\"content\":\"hello\",\"CreateAt\":\"2023-12-20 11:13:56.71999 +0800 CST\"}"))
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
		// Retrieve information
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("Failed to read the message.", err)
			return
		}

		// Put the message body into the global channel.
		// brodMsg(data)
		// Handle the received message
		handleReceivedMessage(node, data)
		var jsonData map[string]interface{}
		err = json.Unmarshal(data, &jsonData)

		if err != nil {
			zap.S().Info("Failed to parse JSON data.", err)
		}
		msgID := jsonData["msgId"]
		if msgID == "-1" {
			//This is a heartbeat message not stored in Redis
		} else {
			pushMsg(data)
		}
	}
}

// handleReceivedMessage processes the received message and sends a response.
func handleReceivedMessage(node *Node, receivedData []byte) {
	// Your message processing logic here
	// For demonstration, let's assume echoing the received message as a response.

	// Process the received data (you can replace this with your own logic)
	processedData, err := processReceivedData(receivedData)

	if err != nil {
		return
	}
	// Send a response back to the client
	err = node.Conn.WriteMessage(websocket.TextMessage, processedData)
	if err != nil {
		zap.S().Info("Failed to write the response message", err)
		return
	}

	fmt.Println("Response sent to the client.")
}

func processReceivedData(data []byte) ([]byte, error) {
	// For demonstration, let's simply echo the received data.
	// Update data with timestamp and status
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	//if msgID, ok := jsonData["msgId"]; ok && msgID == "-1" {
	//	// If msgId is present and equals -1, return a custom error
	//	return nil, &CustomError{"Invalid msgId: -1"}
	//}

	if err != nil {
		zap.S().Info("Failed to parse JSON data.", err)
	}
	// Modify the data
	jsonData["createAt"] = time.Now()
	jsonData["status"] = Succeeded
	// Marshal the updated data
	updatedData, err := json.Marshal(jsonData)
	if err != nil {
		zap.S().Info("Failed to marshal JSON data.", err)
	}
	// Parse the updated data into a Message struct (optional)
	msg := Message{}
	err = json.Unmarshal(updatedData, &msg)
	if err != nil {
		zap.S().Info("Failed to parse the message.", err)
	}
	// Print the parsed data (optional)
	fmt.Println("Parse the data:", "msg.FormId", msg.FormId, "targetId:", msg.TargetId, "type:", msg.Type, "time:", msg.CreateAt)
	// Return the updated data
	return updatedData, nil
}

var upSendChan chan []byte = make(chan []byte, 2048)
var mqSendChan chan []byte = make(chan []byte, 2048)

func brodMsg(data []byte) {
	upSendChan <- data
}

func pushMsg(data []byte) {
	mqSendChan <- data
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

// Creating an Exchange in RabbitMQ
func RabbitmqCreateExchange() {
	//rabbitmqUser := "guest"
	//rabbitmqPassword := "guest"
	//rabbitmqIp := "127.0.0.1"
	//rabbitmqPort := "5672"
	//
	//conn, err := amqp.Dial("amqp://" + rabbitmqUser + ":" + rabbitmqPassword + "@" + rabbitmqIp + ":" + rabbitmqPort + "/")

	rabbitmqHost := global.ServiceConfig.RabbitMQConfig.Host
	rabbitmqPort := global.ServiceConfig.RabbitMQConfig.Port
	rabbitmqUser := global.ServiceConfig.RabbitMQConfig.User
	rabbitmqPassword := global.ServiceConfig.RabbitMQConfig.Password

	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPassword, rabbitmqHost, rabbitmqPort)
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Println(err)
		zap.S().Info("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	zap.S().Info("Failed to open a channel", err)
	defer ch.Close()
	err = ch.ExchangeDeclare(
		"chat-craft-exchange", // name
		"fanout",              // typed
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	zap.S().Info("Failed to declare an exchange", err)
}

func RabbitmqRecProc() {
	//rabbitmqUser := "guest"
	//rabbitmqPassword := "guest"
	//rabbitmqIp := "127.0.0.1"
	//rabbitmqPort := "5672"
	//
	//conn, err := amqp.Dial("amqp://" + rabbitmqUser + ":" + rabbitmqPassword + "@" + rabbitmqIp + ":" + rabbitmqPort + "/")
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", global.ServiceConfig.RabbitMQConfig.User, global.ServiceConfig.RabbitMQConfig.Password, global.ServiceConfig.RabbitMQConfig.Host, global.ServiceConfig.RabbitMQConfig.Port)
	conn, err := amqp.Dial(connString)
	defer conn.Close()
	if err != nil {
		zap.S().Info("Failed to connect to the message queue.", err)
	} else {
		log.Println("Successfully connected to RabbitMQ!")
	}
	ch, err := conn.Channel()
	if err != nil {
		zap.S().Info("Failed to create the MQ channel.", err)
	} else {
		log.Println("Successfully created the MQ channel.")
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"chat-craft-exchange", // name
		"fanout",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		zap.S().Info("Failed to declare the exchange.", err)
	} else {
		log.Println("Successfully declared the exchange.")
	}
	q, err := ch.QueueDeclare(
		"chat-craft-queue-2", // name
		false,                // durable
		false,                // delete when unused
		true,                 // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		zap.S().Info("Failed to declare the queue.", err)
	} else {
		log.Println("Successfully declared the queue.")
	}
	err = ch.QueueBind(
		q.Name,                // queue name
		"",                    // routing key
		"chat-craft-exchange", // exchange
		false,
		nil)
	if err != nil {
		zap.S().Info("Failed to bind to the exchange.", err)
	} else {
		log.Println("Successfully bound")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		zap.S().Info("Failed to consume the MQ message.", err)
	} else {
		log.Println("Successfully consumed the MQ message.")
	}
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			//接收到别的机器发送来的消息
			log.Printf(" [x] %s", d.Body)
			dispatch(d.Body)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

// rabbitmq发送协程
func RabbitmqSendProc() {
	//rabbitmqUser := "guest"
	//rabbitmqPassword := "guest"
	//rabbitmqIp := "127.0.0.1"
	//rabbitmqPort := "5672"
	//conn, err := amqp.Dial("amqp://" + rabbitmqUser + ":" + rabbitmqPassword + "@" + rabbitmqIp + ":" + rabbitmqPort + "/")
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", global.ServiceConfig.RabbitMQConfig.User, global.ServiceConfig.RabbitMQConfig.Password, global.ServiceConfig.RabbitMQConfig.Host, global.ServiceConfig.RabbitMQConfig.Port)
	conn, err := amqp.Dial(connString)
	if err != nil {
		zap.S().Info("Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		zap.S().Info("Failed to open a channel", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"chat-craft-exchange", // name
		"fanout",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)

	if err != nil {
		zap.S().Info("Failed to declare an exchange", err)
	}
	//RabbitMQ Coroutine for Sending: Continuously read from mqsendchan and deliver the message to MQ when there is a message.
	for {
		select {
		case body := <-mqSendChan:
			//If a message is delivered, send this message to the exchange.
			err = ch.Publish(
				"chat-craft-exchange", // exchange
				"",                    // routing key
				false,                 // mandatory
				false,                 // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			if err != nil {
				zap.S().Info("Failed to publish a message", err)
			} else {
				log.Printf(" [x] Sent %s", body)
			}
		}
	}
}

// Dispatch: Parsing the message and determining the chat type.
func dispatch(data []byte) {
	// Update the time field in the data.
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		zap.S().Info("Failed to parse JSON data.", err)
		return
	}
	jsonData["createAt"] = time.Now()
	jsonData["status"] = Succeeded
	updatedData, err := json.Marshal(jsonData)
	if err != nil {
		zap.S().Info("Failed to marshal JSON data.", err)
		return
	}
	msg := Message{}
	err = json.Unmarshal(updatedData, &msg)
	if err != nil {
		zap.S().Info("Failed to parse the message.", err)
		return
	}
	fmt.Println("Parse the data:", "msg.FormId", msg.FormId, "targetId:", msg.TargetId, "type:", msg.Type, "time:", msg.CreateAt)
	switch msg.Type {
	case 1:
		sendMsgAndSave(msg.TargetId, updatedData)
	case 2:
		sendGroupMsg(uint(msg.FormId), uint(msg.TargetId), updatedData)
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
	msgID := jsonMsg.MsgId
	if msgID == "typingId" {
		//This is an input status message and is not stored in Redis
	} else {
		//Concatenate userIdStr and targetIdStr to create a unique key.
		var key string
		if userId > jsonMsg.FormId {
			key = "msg_" + userIdStr + "_" + targetIdStr
		} else {
			key = "msg_" + targetIdStr + "_" + userIdStr
		}
		fmt.Println(key)
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
	fmt.Println(key)
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
