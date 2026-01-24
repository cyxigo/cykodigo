package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cyxigo/cykodigo/bot/data"
	"github.com/gorilla/websocket"
	"modernc.org/libc/limits"
)

func SetupProxy() (*websocket.Dialer, *http.Client, bool) {
	useProxy, ok := data.GetEnvVariable("USE_PROXY")

	if strings.ToLower(useProxy) != "true" || !ok {
		log.Println("Running without proxy")
		return nil, nil, false
	}

	address, ok := data.GetEnvVariable("PROXY_ADDRESS")
	address = strings.TrimSpace(address)

	if !ok {
		log.Println("Can't use proxy: no address")
		return nil, nil, false
	}

	portStr, ok := data.GetEnvVariable("PROXY_PORT")
	portStr = strings.TrimSpace(portStr)

	if !ok {
		log.Println("Can't use proxy: no port")
		return nil, nil, false
	}

	port, err := strconv.Atoi(portStr)

	if port < 1 || port > limits.USHRT_MAX || err != nil {
		log.Printf("Invalid proxy port '%v': %v", portStr, err)
		return nil, nil, false
	}

	username, ok := data.GetEnvVariable("PROXY_USERNAME")
	username = strings.TrimSpace(username)
	username = url.QueryEscape(username)

	if !ok {
		log.Println("Can't use proxy: no username")
		return nil, nil, false
	}

	password, ok := data.GetEnvVariable("PROXY_PASSWORD")
	password = strings.TrimSpace(password)
	password = url.QueryEscape(password)

	if !ok {
		log.Println("Can't use proxy: no password")
		return nil, nil, false
	}

	proxyURL, err := url.Parse(fmt.Sprintf("http://%v:%v@%v:%v",
		username,
		password,
		address,
		portStr))

	if err != nil {
		log.Printf("failed to parse proxy URL: %v", err)
		return nil, nil, false
	}

	log.Printf("Using proxy: %v\n", proxyURL.Redacted())

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	dialer := &websocket.Dialer{
		Proxy: http.ProxyURL(proxyURL),
	}

	return dialer, httpClient, true
}
