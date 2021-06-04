package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in    *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		store: store,
		in:    bufio.NewScanner(in),
	}
}

func extractWinner(line string) string {
	return strings.Replace(line, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func (cli *CLI) PlayPoker() {
	line := cli.readLine()
	cli.store.RecordWin(extractWinner(line))
}
