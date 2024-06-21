package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

type ClientMessage struct {
	Filename      string `json:"filename"`
	OperationTime string `json:"operation_time"`
	Path          string `json:"path"`
	FileType      string `json:"file_type"`
	FileAction    string `json:"file_action"`
	FileSize      int64  `json:"file_size"`
	LastModified  string `json:"last_modified"`
}

func main() {
	// 读取key文件
	key, err := os.ReadFile("key.txt")
	if err != nil {
		log.Println("Unable to read key:", err)
		return
	}

	// 去除key中的空格
	apiKey := strings.TrimSpace(string(key))

	header := http.Header{}
	header.Add("Authorization", "Bearer "+apiKey)

	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/cable"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	identifier, err := json.Marshal(map[string]string{"channel": "FileChannel"})
	if err != nil {
		log.Println("Unable to marshal identifier:", err)
		return
	}

	// 订阅频道
	err = c.WriteJSON(map[string]interface{}{
		"command":    "subscribe",
		"identifier": string(identifier),
	})
	if err != nil {
		log.Println("write:", err)
		return
	}

	done := make(chan struct{})

	// 接收消息
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)

			// 解析消息
			var msg map[string]interface{}
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println("Unable to unmarshal message:", err)
				continue
			}

			// 处理消息
			if messageMap, ok := msg["message"].(map[string]interface{}); ok {
				if fileACtion, ok := messageMap["file_action"].(string); ok {
					filename := messageMap["filename"].(string)
					path := messageMap["path"].(string)
					switch fileACtion {
					case "delete":
						// 删除文件
						err := os.Remove(path + "/" + filename)
						if err != nil {
							log.Println("Unable to delete file:", err)
							return
						}
					}
				}
			}

		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 发送消息
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err := c.WriteMessage(websocket.PingMessage, nil)
				if err != nil {
					log.Println("write:", err)
					return
				}
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("Created file:", event.Name)

					// 获取文件名和文件类型
					filename := filepath.Base(event.Name)
					fileType := strings.TrimPrefix(filepath.Ext(event.Name), ".")

					// 向服务器发送消息
					fileInfo, err := os.Stat(event.Name)
					if err != nil {
						log.Println("Unable to get file info:", err)
						return
					}

					data, err := json.Marshal(ClientMessage{
						Filename:      filename,
						OperationTime: fileInfo.ModTime().Format(time.RFC3339),
						Path:          filepath.Dir(event.Name),
						FileType:      fileType,
						FileAction:    "created",
						FileSize:      fileInfo.Size(),
						LastModified:  fileInfo.ModTime().Format(time.RFC3339),
					})
					if err != nil {
						log.Println("Unable to marshal data:", err)
						return
					}

					err = c.WriteJSON(map[string]interface{}{
						"command":    "message",
						"identifier": string(identifier),
						"data":       string(data),
					})
					if err != nil {
						log.Println("write:", err)
						return
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("File removed: ", event.Name)

					// 获取文件名和文件类型
					filename := filepath.Base(event.Name)
					fileType := strings.TrimPrefix(filepath.Ext(event.Name), ".")

					// 封装数据
					data, err := json.Marshal(ClientMessage{
						Filename:      filename,
						OperationTime: time.Now().Format(time.RFC3339),
						Path:          filepath.Dir(event.Name),
						FileType:      fileType,
						FileAction:    "deleted",
					})
					if err != nil {
						log.Println("Unable to marshal data:", err)
						return
					}

					// 向服务器发送消息
					err = c.WriteJSON(map[string]interface{}{
						"command":    "message",
						"identifier": string(identifier),
						"data":       string(data),
					})
					if err != nil {
						log.Println("write:", err)
						return
					}
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("File moved (possibly deleted): ", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./watch_files")
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
