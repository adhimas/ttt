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

func PrintGame(game *tictactoe.GameUpdate) {
	for i, cell := range game.Board {
		if cell == tictactoe.Empty {
			fmt.Printf(" %d", i)
		} else {
			fmt.Printf(" %s", cell)
		}

		if i%3 == 2 {
			fmt.Println()
		}
	}

	if game.Winner != nil {
		printWinner(game.Piece, *game.Winner)
	}
}

func printWinner(player, winner tictactoe.Cell) {
	switch winner {
	case tictactoe.Empty:
		fmt.Println("tied!")
	case player:
		fmt.Println("you won!")
	default:
		fmt.Println("game over")
	}
}
