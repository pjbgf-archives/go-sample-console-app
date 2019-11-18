package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/pjbgf/go-test/should"
)

func TestMain(t *testing.T) {
	assertThat := func(assumption, command, expectedOutput string) {
		should := should.New(t)
		stdout, _ := ioutil.TempFile("", "calc-fake-stdout.*")
		defer os.Remove(stdout.Name())

		os.Stdout = stdout
		os.Args = strings.Split(command, " ")

		main()

		output, err := ioutil.ReadFile(stdout.Name())
		actualOutput := string(output)

		should.NotError(err, assumption)
		should.BeEqual(expectedOutput, actualOutput, assumption)
	}

	assertThat("should sum 1+1 and return 2", "calc 1 + 1", "sum total: 2\n")
}

func TestMain_ErrorCodes(t *testing.T) {
	assertThat := func(assumption string, command string, expected string) {
		should := should.New(t)
		args := strings.Split(command, " ")

		cmd := exec.Command(args[0], args[1:]...)
		err := cmd.Run()

		e, ok := err.(*exec.ExitError)

		should.BeTrue(ok, assumption)
		should.BeEqual(expected, e.Error(), assumption)
	}

	assertThat("should exit with exit code 1 for invalid syntax", "go test -test.run=TestMain_ErrorCodes main", "exit status 1")
}
