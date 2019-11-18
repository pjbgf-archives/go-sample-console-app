package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/pjbgf/go-sample-console-app/pkg/calc"
)

type console struct {
	commandFactory func(args []string) (cliCommand, error)
	stdOut         io.Writer
	stdErr         io.Writer
	onError        func(writer io.Writer, err error)
}

type cliCommand interface {
	run(output io.Writer)
}

func NewConsole(stdOut io.Writer, stdErr io.Writer) *console {
	if stdOut == (*bytes.Buffer)(nil) {
		panic("stdOut was null")
	}
	if stdErr == (*bytes.Buffer)(nil) {
		panic("stdErr was null")
	}

	return &console{
		getCommand,
		stdOut,
		stdErr,
		exitOnError,
	}
}

func exitOnError(writer io.Writer, err error) {
	printf(writer, "error: %s\n", err)
	os.Exit(1)
}

func (c *console) Run(args []string) {
	cmd, err := c.commandFactory(args)
	if err != nil {
		c.onError(c.stdErr, err)
		return
	}

	cmd.run(c.stdOut)
}

func getCommand(args []string) (cliCommand, error) {
	val1, val2, op, err := parse(args)
	if err != nil {
		return nil, err
	}

	if op == "+" {
		return &additionCommand{val1, val2}, nil
	}

	return nil, nil
}

func parse(args []string) (value1, value2 int, op string, err error) {
	if len(args) < 4 {
		err = errors.New("invalid syntax")
		return
	}

	if value1, err = strconv.Atoi(args[1]); err != nil {
		return
	}

	op = args[2]

	value2, err = strconv.Atoi(args[3])

	return
}

type additionCommand struct {
	value1, value2 int
}

func (a *additionCommand) run(output io.Writer) {
	v := calc.Sum(a.value1, a.value2)

	printf(output, "sum total: %d\n", v)
}

func printf(writer io.Writer, format string, args ...interface{}) {
	writer.Write([]byte(fmt.Sprintf(format, args...)))
}
