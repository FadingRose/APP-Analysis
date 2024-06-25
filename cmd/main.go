package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"fadingrose/app-analyzer/internal/handler"
	"fadingrose/app-analyzer/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {

	logger.Info.Println("Starting server...")

	r := gin.Default()

	// Define the upload route and handler
	r.POST("/upload", handler.UploadHandler)
	r.POST("/upload/url", handler.UploadURLHandler)
	r.POST("/upload/qrcode", handler.UploadQRCodeHandler)
	r.GET("/ws", wsHandler)

	// Serve static frontend
	r.StaticFile("/", "../frontend/index.html")

	cleanup()

	// create /tmp/ if it doesn't exist
	if _, err := os.Stat("/tmp"); os.IsNotExist(err) {
		os.Mkdir("/cache", 0755)
	}

	go func() {
		logger.Info.Println("Starting server on http://localhost:8080")
		if err := r.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	openBrowser("http://localhost:8080")

	select {}
}

func cleanup() {
	// Remove the uploads directory
	if err := os.RemoveAll("./cache"); err != nil {
		panic(err)
	}
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux" and other UNIX-like systems
		cmd = "xdg-open"
		args = []string{url}
	}

	exec.Command(cmd, args...).Start()
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var wsConn *websocket.Conn

func wsHandler(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	wsConn = conn
	go wsListen()
}

func wsListen() {
	defer func() {
		if wsConn != nil {
			wsConn.Close()
			os.Exit(0)
		}
	}()
	for {
		_, _, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
	}
}

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func SendReport(report string) {
	if wsConn != nil {
		msg := Message{
			Type:    "report",
			Message: report,
		}
		msgBytes, _ := json.Marshal(msg)
		wsConn.WriteMessage(websocket.TextMessage, msgBytes)
	}
}
