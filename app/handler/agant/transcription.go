package agant

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	fiberws "github.com/gofiber/websocket/v2"
	gorillaws "github.com/gorilla/websocket"
)

type TranscriptionConfig struct {
	Language string `json:"language"`
	LogID    string `json:"log_id"`
}

type TranscriptionResponse struct {
	LogID  string `json:"log_id"`
	Text   string `json:"text"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// HandleTranscriptionWebSocket handles the WebSocket connection for transcription
func HandleTranscriptionWebSocket(c *fiberws.Conn) {
	defer c.Close()

	// Get initial configuration
	var config TranscriptionConfig
	if err := c.ReadJSON(&config); err != nil {
		log.Printf("Error reading config: %v", err)
		return
	}

	// Connect to Nuera service using Gorilla's websocket (client-side)
	nueraConn, _, err := gorillaws.DefaultDialer.Dial(os.Getenv("NUERA_SERVICE_URL"), nil)
	if err != nil {
		log.Printf("Error connecting to Nuera service: %v", err)
		return
	}
	defer nueraConn.Close()

	// Send initial configuration to Nuera service
	if err := nueraConn.WriteJSON(config); err != nil {
		log.Printf("Error sending config to Nuera service: %v", err)
		return
	}

	// Handle messages in both directions
	go func() {
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("Error reading from client: %v", err)
				return
			}

			if err := nueraConn.WriteMessage(messageType, message); err != nil {
				log.Printf("Error writing to Nuera service: %v", err)
				return
			}
		}
	}()

	for {
		_, message, err := nueraConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from Nuera service: %v", err)
			return
		}

		var response TranscriptionResponse
		if err := json.Unmarshal(message, &response); err != nil {
			log.Printf("Error unmarshaling response: %v", err)
			continue
		}

		// Update the transcription log
		updateData := struct {
			ResponsePayload string `json:"response_payload"`
			Status          string `json:"status"`
		}{
			ResponsePayload: fmt.Sprintf(`{"text": "%s"}`, response.Text),
			Status:          response.Status,
		}

		updateJSON, _ := json.Marshal(updateData)
		if err := c.WriteMessage(fiberws.TextMessage, updateJSON); err != nil {
			log.Printf("Error writing to client: %v", err)
			return
		}
	}
}
