package cli

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/pjbgf/go-test/should"
)

func TestNewConsole(t *testing.T) {
	assertThat := func(assumption string, stdOut, stdErr *bytes.Buffer) {
		should := should.New(t)
		panicked := false

		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()

		NewConsole(stdOut, stdErr)

		should.BeTrue(panicked, assumption)
	}

	var stdOut, stdErr bytes.Buffer
	assertThat("should panic for nil stdOut", nil, &stdErr)
	assertThat("should panic for nil stdErr", &stdOut, nil)
}

func TestRun(t *testing.T) {
	assertThat := func(assumption, expectedOutput string) {
		should := should.New(t)
		var stdOut bytes.Buffer
		hasErrored := false
		c := &console{
			stdOut: &stdOut,
			commandFactory: func(args []string) (cliCommand, error) {
				return nil, errors.New("invalid syntax")
			},
			onError: func(writer io.Writer, err error) { hasErrored = true }}

		c.Run([]string{})

		actualOutput := stdOut.String()

		should.BeEqual(expectedOutput, actualOutput, assumption)
		should.BeTrue(hasErrored, assumption)
	}

	assertThat("should error for invalid commands", "")
}

func TestGetCommand(t *testing.T) {
	assertThat := func(assumption string, command string, expectedCmd cliCommand, expectedErr error) {
		should := should.New(t)

		actualCmd, actualErr := getCommand(strings.Split(command, " "))

		should.BeEqual(expectedErr, actualErr, assumption)
		should.HaveSameType(expectedCmd, actualCmd, assumption)
	}

	assertThat("should return additionCommand for 'calc 1 + 1'", "calc 1 + 1", &additionCommand{}, nil)
	assertThat("should error for 'calc 1 +'", "calc 1 +", nil, errors.New("invalid syntax"))
}

func TestAdditionCommand(t *testing.T) {
	assertThat := func(assumption string, v1, v2 int, expectedOutput string) {
		should := should.New(t)
		cmd := additionCommand{v1, v2}
		var stdOut bytes.Buffer

		cmd.run(&stdOut)

		actualOutput := stdOut.String()

		should.BeEqual(expectedOutput, actualOutput, assumption)
	}

	assertThat("should print 'sum total: 5\n' for 3 + 2", 3, 2, "sum total: 5\n")
	assertThat("should print 'sum total: 79\n' for 40 + 39", 40, 39, "sum total: 79\n")
}
