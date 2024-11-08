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

func (h *Game) checkWinner() (result string) {
	result = ""
	winPoses := [][][]int{
		{{0,0},{0,1},{0,2}},
		{{1,0},{1,1},{1,2}},
		{{2,0},{2,1},{2,2}},

		{{0,0},{1,0},{2,0}},
		{{0,1},{1,1},{2,1}},
		{{0,2},{1,2},{2,2}},

		{{0,0},{1,1},{2,2}},
		{{2,0},{1,1},{0,2}},
	}

	for i := 0; i < len(winPoses); i++ {
		pos1 := winPoses[i][0]
		pos2 := winPoses[i][1]
		pos3 := winPoses[i][2]

		mark1 := h.board[pos1[0]][pos1[1]]
		mark2 := h.board[pos2[0]][pos2[1]]
		mark3 := h.board[pos3[0]][pos3[1]]

		if (mark1 == mark2 && mark1 == mark3 && mark2 == mark3) {
			result = mark1
		}
	}

	return
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
			winner := h.checkWinner()
			for player := range h.players {
				player.board <- board
				if(player.mark != h.turn) {
					player.message <- "Your Turn"
				}
				if(winner == player.mark) {
					player.message <- "You Win"
				}
			}
			if h.turn == "X" { h.turn = "O" } else { h.turn = "X" }
		}
	}
}
