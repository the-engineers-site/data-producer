package handlers

import (
	"bufio"
	"github.com/the-engineers-site/data-producer/pkg/constants"
	"github.com/the-engineers-site/data-producer/pkg/logger"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	Eps     string
	Path    string
	Records []string
}

func LoadCollection() map[string][]*Configuration {
	logger.GetLogger().Info("loading collections")
	mapping := make(map[string][]*Configuration)
	mappingDoc, err := os.Open(constants.MappingFilePath + constants.MappingFileName)
	if err != nil {
		log.Panic("error while reading mapping", err.Error())
	}
	defer mappingDoc.Close()
	scanner := bufio.NewScanner(mappingDoc)
	for scanner.Scan() {
		line := scanner.Text()
		d := strings.Split(line, " ")
		logger.GetLogger().Info("loading server info", zap.String("server", d[0]))
		if _, ok := mapping[d[0]]; ok {
			mapping[d[0]] = append(mapping[d[0]], &Configuration{
				Eps:  d[1],
				Path: d[2],
			})
		} else {
			mapping[d[0]] = []*Configuration{
				{
					Eps:  d[1],
					Path: d[2],
				},
			}
		}
	}
	logger.GetLogger().Info("loaded server", zap.Int("unique servers", len(mapping)))
	return mapping
}

func ReadLogs(mapping map[string][]*Configuration) {
	for _, config := range mapping {
		readLogs(config)
	}
}

func readLogs(config []*Configuration) {
	for _, c := range config {
		records := readFile(c.Path)
		c.Records = records
	}
}

func readFile(path string) []string {
	var data []string
	mappingDoc, err := os.Open(path)
	if err != nil {
		log.Panic("error while reading mapping", err.Error())
	}
	defer mappingDoc.Close()
	scanner := bufio.NewScanner(mappingDoc)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return data
}
