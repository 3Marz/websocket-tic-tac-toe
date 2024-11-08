package main

import "fmt"

type Board [][]string

func newBoard() Board {
	return [][]string {
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
	}
}

type Game struct {
	players map[*Player]bool
	turn    string

	register chan *Player
	unregister chan *Player

	board        Board
	boardChannal chan *Board
}

func newGame() *Game {
	return &Game{
		players:      make(map[*Player]bool),
		turn:         "X",
		register:     make(chan *Player),
		unregister:     make(chan *Player),
		board:        newBoard(),
		boardChannal: make(chan *Board),
	}
}

func (h *Game) run() {
	for {
		select {
		case player := <-h.register:
			h.players[player] = true
			fmt.Println("new player joined")
		case player := <-h.unregister:
			for p := range h.players {
				if p != player {
					p.message <- "Opponent left the game"
				}
				if _, ok := h.players[player]; ok {
					player.conn.Close()
					delete(h.players, player)
					close(player.board)
					close(player.message)
					fmt.Println("player left")
				}
			}
		case board := <-h.boardChannal:
			for player := range h.players {
				player.board <- board
				if(player.mark != h.turn) {
					fmt.Println(player.mark)
					player.message <- "Your Turn"
				}
			}
			if h.turn == "X" { h.turn = "O" } else { h.turn = "X" }
		}
	}
}
