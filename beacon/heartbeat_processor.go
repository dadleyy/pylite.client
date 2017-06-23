package beacon

import "os"
import "log"
import "sync"
import "time"
import "github.com/gorilla/websocket"

// NewHeartbeatProcessor creates a new processor for heartbeats
func NewHeartbeatProcessor(connection *websocket.Conn) *HeartbeatProcessor {
	logger := log.New(os.Stdout, "heartbeat", log.Ldate|log.Ltime|log.Lshortfile)
	return &HeartbeatProcessor{logger, connection}
}

// HeartbeatProcessor is responsible for keeping the websocket connection alive
type HeartbeatProcessor struct {
	*log.Logger
	connection *websocket.Conn
}

// Start launches the hearbeat sequence
func (processor *HeartbeatProcessor) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Minute)

	for _ = range ticker.C {
		writer, e := processor.connection.NextWriter(websocket.TextMessage)

		if e != nil {
			processor.Printf("unable to open up writer: %s", e.Error())
			ticker.Stop()
			break
		}

		if _, e := writer.Write([]byte("ping")); e != nil {
			processor.Printf("unable to write: %s", e.Error())
			ticker.Stop()
			break
		}

		processor.Printf("successfully pinged api host")
	}

}