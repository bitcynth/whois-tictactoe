package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

var board map[int]map[int]string
var lastPlay string
var wmsg string

// Generate a string to display the current status of the board
func boardString() string {
	res := ""
	res = res + "+-------+\r\n"
	res = res + fmt.Sprintf("| %s %s %s |\r\n", board[0][0], board[0][1], board[0][2])
	res = res + fmt.Sprintf("| %s %s %s |\r\n", board[1][0], board[1][1], board[1][2])
	res = res + fmt.Sprintf("| %s %s %s |\r\n", board[2][0], board[2][1], board[2][2])
	res = res + "+-------+\r\n"
	return res
}

// Check if someone has won
func checkWin() string {
	if (board[0][0] != " ") && (board[0][0] == board[0][1]) && (board[0][1] == board[0][2]) {
		return board[0][0]
	}
	if (board[1][0] != " ") && (board[1][0] == board[1][1]) && (board[1][1] == board[1][2]) {
		return board[1][0]
	}
	if (board[2][0] != " ") && (board[2][0] == board[2][1]) && (board[2][1] == board[2][2]) {
		return board[2][0]
	}
	if (board[0][0] != " ") && (board[0][0] == board[1][0]) && (board[1][0] == board[2][0]) {
		return board[0][0]
	}
	if (board[0][1] != " ") && (board[0][1] == board[1][1]) && (board[1][1] == board[2][1]) {
		return board[0][1]
	}
	if (board[0][2] != " ") && (board[0][2] == board[1][2]) && (board[1][2] == board[2][2]) {
		return board[0][2]
	}
	if (board[0][0] != " ") && (board[0][0] == board[1][1]) && (board[1][1] == board[2][2]) {
		return board[0][0]
	}
	if (board[0][2] != " ") && (board[0][2] == board[1][1]) && (board[1][1] == board[2][0]) {
		return board[0][2]
	}
	return ""
}

func handleConn(c net.Conn) {
	inp, err := bufio.NewReader(c).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	if inp == "blink\r\n" {
		res := "\033[5m\033[93mBlinky!\033[0m"
		c.Write([]byte(res))
	}

	cmd := strings.ReplaceAll(inp, "\r\n", "")
	fmt.Println(cmd)
	parts := strings.Split(cmd, " ")

	if lastPlay == parts[0] {
		c.Write([]byte("You just played!\r\n"))
		c.Write([]byte(boardString()))
		c.Close()
		return
	}

	x, _ := strconv.Atoi(parts[1])
	y, _ := strconv.Atoi(parts[2])
	x = x - 1
	y = y - 1
	if x < 0 || x > 2 || y < 0 || y > 2 {
		c.Write([]byte("Invalid coordinates!\r\n"))
		c.Write([]byte(boardString()))
		c.Close()
		return
	}
	if wmsg != "" {
		c.Write([]byte(wmsg))
		c.Write([]byte(boardString()))
		c.Close()
		return
	}
	if board[x][y] == " " {
		board[x][y] = string(parts[0][0])
		fmt.Println(boardString())
		lastPlay = parts[0]
		w := checkWin()
		if w != "" {
			fmt.Println(w, "won!")
			wmsg = fmt.Sprintf("\033[32m\033[5m%s won!\033[0m\r\n", w)
			c.Write([]byte(wmsg))
		}
	} else {
		c.Write([]byte("\033[5mThat tile is taken, take another.\r\n\033[0m"))
	}

	c.Write([]byte(boardString()))
	c.Close()
}

func main() {
	wmsg = ""
	board = make(map[int]map[int]string)
	for x := range []int{0, 1, 2} {
		board[x] = make(map[int]string)
		for y := range []int{0, 1, 2} {
			board[x][y] = " "
		}
	}
	l, err := net.Listen("tcp4", ":8043")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConn(c)
	}
}
