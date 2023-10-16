package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	todo "github.com/meru_2802/GoDoList"
)

const (
	todoFile = "data/todos.mayur"
	historyFile = "data/history.mayur"
)

func main() {
	add := flag.Bool("add", false, "add a new todo")
	complete := flag.Int("complete", 0, "mark a todo completed")
	del := flag.Int("del", 0, "delete a todo")
	list := flag.Bool("list", false, "show the list of todos")
	showCompleted := flag.Bool("showCompleted", false, "show the list of completed todos")

	flag.Parse()

	todos := &todo.Todos{}

	if err := todos.Load(todoFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		todos.Add(task)

		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Print()
	case *complete > 0:
		err := todos.Complete(*complete, historyFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		// Store todos
		if err := todos.Store(todoFile); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		todos.Print()

	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = todos.Store(todoFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *list:
		todos.Print()
	case *showCompleted:
		todos.PrintCompleted()

	default:
		fmt.Fprintln(os.Stdout, "invalid command")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return scanner.Text(), nil
}
