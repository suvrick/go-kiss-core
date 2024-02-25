package proxy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	resource = "https://lumtest.com/myip.json"
	host     = "brd.superproxy.io:22225"
	user_id  = "brd-customer-hl_07f044e7-zone-static"
	password = "hcx7fnqnph27"
)

type Proxy struct {
	Ip string
}

func GetNetProxy(session string) func(*http.Request) (*url.URL, error) {

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(&url.URL{
				Scheme: "http",
				Host:   host,
				User:   url.UserPassword(user_id, password),
			}),
		},
	}

	// 55.106
	//
	resp, err := client.Get(resource)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	p := Proxy{}

	err = json.Unmarshal(body, &p)
	if err != nil {
		return nil
	}

	// brd.superproxy.io:22225:brd-customer-hl_07f044e7-zone-static-ip-178.171.116.192:hcx7fnqnph27
	// brd.superproxy.io:22225:brd-customer-hl_07f044e7-zone-static-ip-178.171.116.193:hcx7fnqnph27
	// brd.superproxy.io:22225:brd-customer-hl_07f044e7-zone-static-ip-178.171.116.195:hcx7fnqnph27

	newUser := fmt.Sprintf("%s-ip-%s", user_id, p.Ip)

	switch session {
	case "one":
		newUser = "brd-customer-hl_07f044e7-zone-static-ip-178.171.116.192"
	case "two":
		newUser = "brd-customer-hl_07f044e7-zone-static-ip-178.171.116.193"
	case "tree":
		newUser = "brd-customer-hl_07f044e7-zone-static-ip-178.171.116.195"
	}

	fmt.Printf("generate proxy: %s\n", newUser)

	u := url.URL{
		Scheme: "http",
		Host:   host,
		User:   url.UserPassword(newUser, password),
	}

	return http.ProxyURL(&u)
}

// "{\"ip\":\"212.80.221.241\",\"country\":\"IE\",\"asn\":{\"asnum\":9009,\"org_name\":\"M247 Europe SRL\"},\"geo\":{\"city\":\"Dublin\",\"region\":\"L\",\"region_name\":\"Leinster\",\"postal_code\":\"D12\",\"latitude\":53.323,\"longitude\":-6.3159,\"tz\":\"Europe/Dublin\",\"lum_city\":\"dublin\",\"lum_region\":\"l\"}}
