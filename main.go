package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	name        string
	age         int
	lastUpdated int
	comment     string
}

var records []Record

func main() {
	records = make([]Record, 0)

	readDB("database.csv")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		command, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		command = command[0 : len(command)-1]
		arguments := strings.Split(command, " ")
		if len(arguments) == 0 {
			fmt.Println("Please enter a command")
		} else if command == "print" && len(arguments) == 1 {
			printDB()
		} else if command == "exit" && len(arguments) == 1 {
			break
		} else if command == "update" && len(arguments) == 1 {
			writeDB("database.csv")
		} else if len(arguments) == 3 && arguments[0] == "add" {
			fmt.Println("enter a comment: ")
			comment, err := reader.ReadString('\n')
			if err != nil {
				continue
			}
			comment = comment[0 : len(comment)-1]
			age, err := strconv.Atoi(arguments[2])
			if err != nil {
				fmt.Println("please enter valid age")
				continue
			}
			lastUpdated := time.Now().UTC()
			record := Record{name: arguments[1], age: age, lastUpdated: int(lastUpdated.Unix()), comment: comment}
			records = append(records, record)
		} else if len(arguments) == 2 && arguments[0] == "remove" {
			for i, record := range records {
				if record.name == arguments[1] {
					records = append(records[:i], records[i+1:]...)
					continue
				}
			}
		} else {
			fmt.Println("please enter a command")
		}
	}
	printDB()
	writeDB("database.csv")
}

func readDB(fileName string) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err.Error())
		}
		line = line[0 : len(line)-1]

		arguments := strings.Split(line, ",")

		if len(arguments) != 4 {
			return
		}
		var record Record
		record.name = arguments[0]
		record.age, _ = strconv.Atoi(arguments[1])
		record.lastUpdated, _ = strconv.Atoi(arguments[2])
		record.comment = arguments[3]
		records = append(records, record)
	}
}

func writeDB(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer file.Close()

	for _, record := range records {
		_, err := file.WriteString(record.name + "," + strconv.Itoa(record.age) + "," + strconv.Itoa(record.lastUpdated) + "," + record.comment + "\n")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	return
}

func printDB() {
	for i := 0; i < len(records); i++ {
		fmt.Println("-------------------------------------------------------------------------------")
		tm := time.Unix(int64(records[i].lastUpdated), 0)
		fmt.Printf("|%15v|%4v|%25v|%30v|\n", records[i].name, records[i].age, tm.Format("2006-01-02 15:04:05"), records[i].comment)
	}
	fmt.Println("-------------------------------------------------------------------------------")
}
