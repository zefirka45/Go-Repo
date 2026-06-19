// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"encoding/json"
	"go-server/jwt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	// Время, отведенное на отправку сообщения другому узлу.
	writeWait = 10 * time.Second

	// Время, отведенное на чтение следующего сообщения Pong от удаленного узла.
	pongWait = 60 * time.Second

	// Отправка пингов пиру с указанным периодом. Должен быть меньше, чем pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Максимально допустимый размер сообщения от партнера.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Клиент выступает посредником между соединением WebSocket и хабом.
type Client struct {
	hub *Hub

	// Соединение через веб-сокет.
	conn *websocket.Conn

	// Буферизованный канал исходящих сообщений.
	send   chan []byte
	userID string
}

// Функция readPump перенаправляет сообщения из соединения WebSocket в хаб.

// Приложение запускает readPump в горутине для каждого соединения. Приложение
// гарантирует, что на соединении находится не более одного читателя, выполняя все
// операции чтения из этой горутины.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		payload = bytes.TrimSpace(bytes.Replace(payload, newline, space, -1))
		var msg message
		if err := json.Unmarshal(payload, &msg); err != nil {
			log.Printf("invalid message format from user %s: %v", c.userID, err)
			continue
		}
		msg.From = c.userID
		c.hub.broadcast <- msg
	}
}

// Функция writePump перенаправляет сообщения из хаба в соединение WebSocket.

// Для каждого соединения запускается горутина, выполняющая writePump.
// Приложение гарантирует, что к соединению обращается не более одного пользователя,
// выполняя все операции записи из этой горутины..
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Хаб закрыл канал.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Добавить сообщения чата из очереди в текущее сообщение веб-сокета.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs обрабатывает запросы WebSocket от удаленного узла.
func ServeWs(hub *Hub, c *gin.Context) {
	// Пример получения ID: /ws?userId=user123
	token := c.Query("token")
	claims, err := jwt.ValidateToken(token)
	if err != nil {
		log.Println("Rejecting connection: userId is required")
		c.JSON(400, gin.H{"error": "userId query param is required"})
		return
	}

	userID, ok := claims["id"].(string)
	if !ok {
		log.Println("Rejecting connection: 'id' claim is missing or not a string")
		c.JSON(400, gin.H{"error": "Bad Request: invalid token payload"})
		return // Обязательно прерываем выполнение хендлера!
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
