package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	msgBuf = 10
)

type Message struct {
	Text string
}

type Response map[string]interface{}

func (r Response) String() (s string) {
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}

type Broker struct {
	subscribers map[chan []byte]bool
}

func (b *Broker) Subscribe() chan []byte {
	ch := make(chan []byte, msgBuf)
	b.subscribers[ch] = true
	return ch
}

func (b *Broker) Unsubscribe(ch chan []byte) {
	delete(b.subscribers, ch)
}

func (b *Broker) Publish(msg []byte) {
	for ch := range b.subscribers {
		ch <- msg
	}
}

func NewBroker() *Broker {
	return &Broker{make(map[chan []byte]bool)}
}

var msgBroker *Broker

func messageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type")[:16] != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusNotAcceptable)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(body))
	msgBroker.Publish([]byte(string(body)))

	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, Response{"success": true, "message": "OK"})
		return
	}

	fmt.Fprintln(w, "OK")
}

func timerEventSource(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := msgBroker.Subscribe()
	defer msgBroker.Unsubscribe(ch)

	for {
		msg := <-ch
		fmt.Fprintf(w, "data: %s\n\n", msg)
		f.Flush()
	}
}

type Server struct {
	broker *Broker
	port   uint64
}

func NewServer(port uint64, debug bool) Server {
	msgBroker = NewBroker()
	return Server{broker: msgBroker, port: port}
}

func (s Server) Run() {
	http.HandleFunc("/update", messageHandler)
	http.HandleFunc("/events", timerEventSource)
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Printf("Listening on port %d", s.port)
	err := http.ListenAndServe(":"+strconv.FormatUint(s.port, 10), nil)
	if err != nil {
		panic(err)
	}
}

func main() {
	server := NewServer(3333, false)
	go server.Run()

	bio := bufio.NewReader(os.Stdin)
	for {
		// we ignore isPrefix since we are reading small lines anyway
		line, _, err := bio.ReadLine()
		// 398605931381 total_distance_meters:76
		tokens := strings.Split(string(line), " ")
		if len(tokens) < 2 {
			continue
		}
		timestamp := tokens[0]
		tokens = strings.Split(tokens[1], ":")
		key := tokens[0]
		value := tokens[1]
		msg := fmt.Sprintf("{\"timestamp\": %s, \"label\": \"%s\", \"value\": %s}", timestamp, key, value)
		if err != nil {
			log.Fatal(err)
		}
		server.broker.Publish([]byte(msg))
		time.Sleep(25 * time.Millisecond)
	}

}
