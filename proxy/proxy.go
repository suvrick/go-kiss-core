package proxy

import (
	"errors"
	"net/url"
	"sync"
	"time"
)

var (
	ErrEmptyProxy = errors.New("empty proxy")
)

var Instance *ProxyManager
var one sync.Once

type Proxy struct {
	URL     *url.URL
	IsBad   bool
	TimeUse time.Time
}

type ProxyManager struct {
	P       chan *Proxy
	proxies []Proxy
}

// Create new instanse
func NewProxyManager() *ProxyManager {
	if Instance == nil {
		one.Do(func() {
			Instance = &ProxyManager{
				P:       make(chan *Proxy),
				proxies: make([]Proxy, 0),
			}
		})
	}

	return Instance
}

func (pm *ProxyManager) Load([]Proxy) {

}

func loop() {
	for i := 0; i > len(Instance.proxies); i++ {
		v := Instance.proxies[i]
		if v.IsBad {
			continue
		}
		Instance.P <- &v
	}
}

func (pm *ProxyManager) GetProxy() *Proxy {
	return <-pm.P
}

func (pm *ProxyManager) UpdateProxy(p *Proxy) {

}

func GetDefaultProxy() *Proxy {
	return &Proxy{
		URL: &url.URL{
			Scheme: "http",
			Host:   "zproxy.lum-superproxy.io:22225",
			User:   url.UserPassword("lum-customer-c_07f044e7-zone-static", "hcx7fnqnph27"),
		},
	}
}
