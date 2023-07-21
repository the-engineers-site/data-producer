package store

import (
	"bytes"
	"fmt"
	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/the-engineers-site/data-producer/pkg/logger"
	"go.uber.org/zap"
	"html/template"
	"log"
	"math/rand"
	"net"
	"time"
)

type Formatter struct {
	DateSyslog    string
	LongTimeStamp string
	FullTimeStamp string
	PrivateIp     string
	PublicIp      string
	Hostname      string
	Port          int
	Url           string
	Username      string
	RandomNumber  int
}

func Send(counter *int, connection *net.Conn, message string) {
	tmpl, err := template.New("myTemplate").Parse(message)
	if err != nil {
		logger.GetLogger().Error("error while processing template", zap.Error(err))
		return
	}
	rand.Seed(time.Now().UnixNano())
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, getRandomObject())
	sendLineAsync(counter, connection, buffer.String())
}

func sendLineAsync(count *int, connection *net.Conn, message string) {
	_, err := fmt.Fprintln(*connection, message)
	if err != nil {
		log.Println("Error while publishing ", err)
	}
	time.Sleep(1 * time.Second)
	*count++
}

func getRandomObject() (formatter Formatter) {
	formatter = Formatter{
		DateSyslog:    time.Now().Format("Jan 2 15:04:05"),
		FullTimeStamp: time.Now().Format("01/02/2006:15:04:05 MST"),
		LongTimeStamp: time.Now().Format("2006-02-01 15:04:05"),
		PrivateIp:     fake.IPv4Address(),
		PublicIp:      fake.IPv4Address(),
		Hostname:      generateRandomHostname(),
		Port:          generateRandomPort(),
		Url:           generateRandomUrl(),
		RandomNumber:  fake.Number(0, 200000),
		Username:      fake.Username(),
	}
	return formatter
}

func generateRandomPort() int {
	// Generate a random number within the valid port number range (0 to 65535)
	port := rand.Intn(65536)

	return port
}

func generateRandomUrl() string {
	// Generate a random string of length 8 for the hostname
	randomString := generateRandomString(4) + "." + generateRandomString(4)

	// Combine the random string with the base domain
	hostname := randomString + ".databahn.ai"

	return hostname
}

func generateRandomHostname() string {
	// Generate a random string of length 8 for the hostname
	randomString := generateRandomString(8)

	// Combine the random string with the base domain
	hostname := randomString + ".databahn.ai"

	return hostname
}

func generateRandomString(length int) string {
	// Characters allowed in the random string
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"

	// Generate a random string of specified length
	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomString[i] = chars[rand.Intn(len(chars))]
	}

	return string(randomString)
}
