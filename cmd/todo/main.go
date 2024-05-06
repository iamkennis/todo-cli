package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	todo "github.com/iamkennis/todocli"
)

var todoFileName = "./todo.json"

type Stringer interface {
	String() string
}

func main() {
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "list all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Delete task in the list")
	date := flag.Int("date", 0, "Show time in a task")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	notdone := flag.Bool("notdone", false, "Not done tasks")

	flag.Parse()

	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Print(l)
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		}
	case *del > 0:
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		}

	case *date > 0:
		if err := l.ShowDateTime(*date, *verbose); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *notdone:
		if err := l.NotDoneItem(*notdone); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		}

	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		}
		l.Add(t)
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)

		}

	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)

	}

	//		if os.Getenv("TODO_FILENAME") != "" {
	//		todoFileName = os.Getenv("TODO_FILENAME")
	//	}
}

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	var tasks []string

	s := bufio.NewScanner(r)
	for s.Scan() {
		text := s.Text()
		if len(text) == 0 {
			continue
		}
		tasks = append(tasks, text)
	}
	if err := s.Err(); err != nil {
		return "", err
	}

	return strings.Join(tasks, "\n"), nil
}
