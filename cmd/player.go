package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Player struct {
	mark string

	game *Game
	//WebSocket Connection
	conn *websocket.Conn

	board   chan *Board
	message chan string
}

type PlayReq []int

// Reading encomming board data from client
func (p *Player) readPump() {
	defer func() {
		p.game.unregister <- p
		p.conn.Close()
	}()

	for {
		var data PlayReq
		err := p.conn.ReadJSON(&data)
		if err != nil {
			log.Println(err)
			break
		}
		if p.game.board[data[0]][data[1]] == " " {
			p.game.board[data[0]][data[1]] = p.mark
			p.game.boardChannel <- &p.game.board
		}
	}
}

func (p *Player) writePump() {
	defer func() {
		p.game.unregister <- p
		p.conn.Close()
	}()

	for {
		select {
		case board := <-p.board:
			err := p.conn.WriteJSON(board)
			if err != nil {
				log.Println(err)
				return
			}
		case msg := <-p.message:
			p.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

func serveWs(lobby *Lobby, mark string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	player := &Player{
		mark:    mark,
		conn:    conn,
		board:   make(chan *Board),
		message: make(chan string),
	}
	lobby.registerPlayer <- player
}
