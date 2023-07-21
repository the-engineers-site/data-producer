package main

import (
	"bufio"
	"github.com/the-engineers-site/data-producer/pkg/handlers"
	"github.com/the-engineers-site/data-producer/pkg/logger"
	"github.com/the-engineers-site/data-producer/pkg/producer"
	"log"
	"os"
)

func main() {
	logger.GetLogger().Info("starting service")
	mapping := handlers.LoadCollection()
	handlers.ReadLogs(mapping)
	producers := producer.FormProducers(mapping)
	producer.StartAll(producers)
}

func readFile() []string {
	var lines []string
	file, err := os.Open("/Users/yjagdale/Documents/code/rnd/data-producer/cmd/data.log")
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lines
}
