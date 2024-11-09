package main

import (
	"fmt"
	"math/rand/v2"

	// "github.com/gorilla/websocket"
)

type Lobby struct {
	players []*Player
	games map[*Game]bool

	registerPlayer chan *Player
	registerGame chan *Game
	unregisterGame chan *Game
}

func newLobby() *Lobby {
	return &Lobby{
		players:  []*Player{},
		games:    make(map[*Game]bool),
		registerPlayer: make(chan *Player),
		registerGame: make(chan *Game, 1),
		unregisterGame: make(chan *Game),
	}
}

func (l *Lobby) run() {
	for {
		select {
		case player := <-l.registerPlayer:
			l.players = append(l.players, player)
		case game := <-l.registerGame:
			l.games[game] = true
			fmt.Println("Register game")
		case game := <-l.unregisterGame:
			if _, ok := l.games[game]; ok {
				delete(l.games, game)
				close(game.boardChannel)
				close(game.register)
				close(game.unregister)
				fmt.Println("Unregister game")
			}
		}

		if len(l.players) >= 2 {
			fmt.Println("start game")
			startGame(l.players[len(l.players)-1], l.players[len(l.players)-2],l)
			//Delete the last two players
			l.players = l.players[:len(l.players)-2]
		}
	}
}

func startGame(p1 *Player, p2 *Player, lobby *Lobby) {
	game := newGame()
	go game.run()

	p1.game = game
	p2.game = game
	p1.game.register <- p1
	p2.game.register <- p2

	go p1.writePump()
	go p2.writePump()
	go p1.readPump()
	go p2.readPump()

	p1.message <- "Got A Game"
	p2.message <- "Got A Game"

	if rand.IntN(2) == 0 {
		p1.mark = "X"
		p1.message <- "Your Turn"

		p2.mark = "O"
	} else {
		p1.mark = "O"

		p2.message <- "Your Turn"
		p2.mark = "X"
	}
	game.lobby = lobby
	game.lobby.registerGame <- game
}
