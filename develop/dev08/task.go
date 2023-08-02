package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качестве аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена
команда выхода (например \quit).

*/

// cd changes the current working directory.
func cd(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("expected 1 argument, but got %d", len(args)-1)
	}
	return os.Chdir(args[1])
}

// pwd prints the current working directory.
func pwd() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

// execCmd executes a command in a new process.
func execCmd(args []string, stdin io.Reader, stdout io.Writer) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout

	return cmd.Run()
}

func pipeline(commands [][]string) error {
	// Create a slice to hold the readers and writers of the pipes
	readers := make([]*io.PipeReader, len(commands)-1)
	writers := make([]*io.PipeWriter, len(commands)-1)
	for i := range readers {
		readers[i], writers[i] = io.Pipe()
		defer readers[i].Close()
		defer writers[i].Close()
	}

	// Start each command in a goroutine, with its stdin connected to the previous pipe's reader,
	// and its stdout connected to the next pipe's writer.
	for i := range commands[:len(commands)-1] {
		go func(i int) {
			execCmd(commands[i], readers[i], writers[i])
		}(i)
	}

	// Start the final command, with its stdin connected to the last pipe's reader,
	// and its stdout connected to os.Stdout.
	return execCmd(commands[len(commands)-1], readers[len(readers)-1], os.Stdout)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		// Print the shell prompt
		fmt.Print("> ")

		// Read a line of input from the user
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Parse the line into a command and arguments
		line = strings.TrimSpace(line)
		args := strings.Split(line, " ")
		pipelineCmds := strings.Split(line, "|")

		if len(pipelineCmds) > 1 {
			commands := make([][]string, len(pipelineCmds))
			for i, cmdLine := range pipelineCmds {
				commands[i] = strings.Split(strings.TrimSpace(cmdLine), " ")
			}

			// Handle pipeline of commands
			if err := pipeline(commands); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}

		// Handle the command
		switch args[0] {
		case "cd":
			if err := cd(args); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		case "pwd":
			if err := pwd(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		case "exit":
			return
		default:
			// If the command is not a built-in command, execute it in a new process
			if err := execCmd(args, os.Stdin, os.Stdout); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}
