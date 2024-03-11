package proxy

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	host     = "brd.superproxy.io:22225"
	user_id  = "brd-customer-hl_07f044e7-zone-static-session-rand"
	password = "hcx7fnqnph27"
)

type Proxy struct {
	Ip string
}

func GetNetProxy2(session string) func(*http.Request) (*url.URL, error) {

	file, err := os.ReadFile("./proxy/ips-static.txt")
	if err != nil {
		log.Println(err)
	}

	text := string(file)

	lines := strings.Split(text, "\n")

	i := rand.Intn(len(lines))

	line := lines[i]

	log.Printf("net-proxy: %s\n", line)

	arr := strings.Split(line, ":")

	u := url.URL{
		Scheme: "http",
		Host:   host,
		User:   url.UserPassword(arr[2], arr[3]),
	}

	return http.ProxyURL(&u)
}

func GetUserAgent(session string) string {

	file, err := os.ReadFile("./proxy/user-agents.txt")
	if err != nil {
		log.Println(err)
	}

	text := string(file)

	lines := strings.Split(text, "\n")

	i := rand.Intn(len(lines))

	line := lines[i]

	log.Printf("user-agent: %s\n", line)

	return line
}

func GetNetProxy(session string) func(*http.Request) (*url.URL, error) {

	user_id2 := fmt.Sprintf("brd-customer-hl_07f044e7-zone-static-session-%s", session)
	log.Printf("sesstion: %s, login: %s\n", session, user_id2)
	u := url.URL{
		Scheme: "http",
		Host:   host,
		User:   url.UserPassword(user_id2, password),
	}

	return http.ProxyURL(&u)
}
