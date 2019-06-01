package allocine

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	host string
	port int
	url string
	client http.Client
}

func NewClient(host string, port int) *Client {
	url := fmt.Sprintf("http://%s:%d", host, port)
	return &Client{host: host, port: port, url: url, client: http.Client{}}
}


func (c Client) makeRequest(location string) (result []byte, err error) {
	endpoint := strings.Join([]string{c.url, "showtimes"}, "/")
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Print(err)
		return
	}

	q := req.URL.Query()
	q.Add("location", location) // Hardly define Mouans-Sartoux
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		return
	}

	return body, err
}

func (c Client) GetLastShowTime(location string) (result *gabs.Container, err error) {

	raw, err := c.makeRequest(location)
	if err != nil {
		return
	}

	return gabs.ParseJSON(raw)
}