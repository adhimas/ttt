package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"example.com/ttt/internal/tictactoe"
)

func PrintGame(board []tictactoe.Cell) {
	for i, cell := range board {
		if cell == tictactoe.Empty {
			fmt.Printf(" %d", i)
		} else {
			fmt.Printf(" %s", cell.String())
		}

		if i%3 == 2 {
			fmt.Println()
		}
	}
}

func PrintWinner(game tictactoe.GameUpdate) {
	switch *game.Winner {
	case tictactoe.Empty:
		log.Printf("tied!")
	case game.Piece:
		log.Printf("you won!")
	default:
		log.Printf("game over")
	}
}

func Prompt(player tictactoe.Cell) int {
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("enter number for %s: ", player)

		if more := scanner.Scan(); !more {
			continue
		}
		// TODO: handle interrupt
		if err := scanner.Err(); err != nil {
			log.Println("scanner:", err)
			continue
		}

		line := scanner.Text()
		cellNum, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			log.Println("prompt input:", err)
			continue
		}
		if cellNum < 0 || cellNum > 8 { // 0-8
			log.Println("prompt input: invalid move")
			continue
		}

		fmt.Println()
		return cellNum
	}
}
