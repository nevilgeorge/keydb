package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const DB_FILENAME = "db/local.db"

func main() {
	fmt.Println("Hi from KeyDB.")

	recordStore, err := NewFileRecordStore(DB_FILENAME)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create RecordStore: ", err)
		return
	}
	defer recordStore.Close()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			// EOF
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		cmd, rest, _ := strings.Cut(input, " ")
		rest = strings.TrimSpace(rest)

		switch cmd {
		case "put":
			err := handlePut(recordStore, rest)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Record handling error: ", err)
			}
		case "get":
			err := handleGet(recordStore, rest)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Record handling error: ", err)
			}
		case "exit", "quit":
			return
		default:
			continue
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Input error: ", err)
		}
	}
}

func handlePut(recordStore RecordStore, rest string) error {
	if !strings.Contains(rest, ":") {
		return fmt.Errorf("Expected a key/value pair split by :, you entered %s", rest)
	}

	stringParts := strings.Split(rest, ":")
	for i, p := range stringParts {
		stringParts[i] = strings.TrimSpace(p)
	}

	if len(stringParts) != 2 {
		return fmt.Errorf("Expected one key and one value, go this instead: %s", stringParts)
	}
	key := stringParts[0]
	value := stringParts[1]
	err := recordStore.Put(key, value)
	if err != nil {
		return err
	}
	fmt.Printf("Record stored: %s: %s", key, value)
	fmt.Println()

	return nil
}

func handleGet(recordStore RecordStore, key string) error {
	if key == "" {
		return fmt.Errorf("Expected a key")
	}

	value, err := recordStore.Get(key)
	if err != nil {
		return err
	}
	fmt.Println(value)
	return nil
}
