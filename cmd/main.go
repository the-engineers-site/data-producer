package main

import (
	"bytes"
	"fmt"
	"github.com/the-engineers-site/data-producer/pkg/store"
	"log"
	"math/rand"
	"net"
	"os"
	"text/template"
	"time"
)

type Formatter struct {
	DateSyslog    string
	FullTimeStamp string
	PrivateIp     string
	PublicIp      string
	Hostname      string
	Port          int
	Url           string
}

func main() {
	log.Println("starting service")
	logLine := os.Getenv("LOG")
	logMessage := "test {{ .DateSyslog}} {{ .PrivateIp }} {{ .Hostname }}|src={{ .PrivateIp }}|srcport=63917|dst={{ .PublicIp }}|dstport=443|username=-|devicetime=[{{ .FullTimeStamp }}]|s-action=TCP_TUNNELED|sc-status=200|cs-method=CONNECT|time-taken=123918|sc-bytes=14767|cs-bytes=13941|cs-uri-scheme=tcp|cs-host={{ .Url }}|cs-uri-path=/|cs-uri-query=-|cs-uri-extension=-|cs-auth-group=-|rs(Content-Type)=-|cs(User-Agent)=golang-producer|cs(Referer)=-|sc-filter-result=OBSERVED|filter-category=Chat_(IM)/SMS|cs-userdn=-|cs-uri={{ .FullTimeStamp }}/|x-virus-id=-|s-ip={{ .PrivateIp }}|s-sitename=http.proxy"

	if logLine == "" {
		logLine = logMessage
	}

	tmpl, err := template.New("myTemplate").Parse(logLine)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	for {
		rand.Seed(time.Now().UnixNano())
		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, getRandomObject())
		store.Send(buffer.String())
	}

	// Dec 1 06:46:09 10.246.7.37 Bluecoat 01/12/2022:11:44:19 GMT
}

func getRandomObject() (formatter Formatter) {
	formatter = Formatter{
		DateSyslog:    time.Now().Format("Jan 2 15:04:05"),
		FullTimeStamp: time.Now().Format("01/02/2006:15:04:05 MST"),
		PrivateIp:     generateRandomIP("10.0.0.0/8").String(),
		PublicIp:      generateRandomIP("54.0.0.0/8").String(),
		Hostname:      generateRandomHostname(),
		Port:          generateRandomPort(),
		Url:           generateRandomUrl(),
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

func generateRandomIP(cidr string) net.IP {
	// Parse the CIDR range
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}

	// Generate a random IP within the CIDR range
	ip := make(net.IP, len(ipNet.IP))
	for i, b := range ipNet.IP {
		ip[i] = b
	}

	// Generate random bytes for the host portion
	for i := len(ip) - 1; i >= len(ip)-len(ipNet.Mask); i-- {
		ip[i] = byte(rand.Intn(256))
	}

	return ip
}
