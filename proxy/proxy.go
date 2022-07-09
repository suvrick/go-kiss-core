package proxy

import (
	"fmt"
	"strings"
)

type Proxy struct {
	Scheme   string
	Host     string
	User     string
	Password string
}

// TODO: impliment IProxyManager
// proxy := &Proxy{
// 	Scheme:   "http",
// 	Host:     "zproxy.lum-superproxy.io:22225",
// 	User:     "lum-customer-c_07f044e7-zone-static",
// 	Password: "lon38jik65hm",
// 	Timeout:  time.Second * 60,
// }

func Parse(url string, seporator string) *Proxy {

	temp := strings.Split(url, seporator)

	if len(temp) != 4 {
		return nil
	}

	result := [4]string{}

	for i, v := range temp {
		substr := v
		result[i] = substr
	}

	return &Proxy{
		Scheme:   "http",
		Host:     fmt.Sprintf("%s:%s", result[0], result[1]),
		User:     result[2],
		Password: result[3],
	}
}
