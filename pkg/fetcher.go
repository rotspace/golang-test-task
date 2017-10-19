package tagc

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/fetchbot"
)

func (t *Tagc) initFetcher() {
	f := fetchbot.New(fetchbot.HandlerFunc(t.fetchbotHandler))
	t.fb = f
	t.queue = f.Start()
	return
}

func (t *Tagc) fetchbotHandler(ctx *fetchbot.Context, res *http.Response, err error) {
	cmd, ok := ctx.Cmd.(*command)
	if !ok {
		panic("unreacheble")
	}
	defer cmd.Done()
	if err != nil {
		return
	}

	var (
		chunked     = res.Header.Get("Transfer-Encoding") == "chunked"
		contentType = strings.TrimSpace(strings.Split(res.Header.Get("Content-type"), ";")[0])

		status = res.StatusCode
	)

	dup, contentLen, err := CountTags(res.Body, chunked)
	if err != nil {
		cmd.err = err
		return
	}
	if status >= 200 && status < 300 && contentLen > 0 {
		cmd.result = newResponse(cmd.rawurl, status, contentType, contentLen, newResponseElements(dup))
	}
}

type command struct {
	uri    *url.URL
	rawurl string
	done   chan struct{}
	result response
	err    error
}

func newCmd(rawurl string) (*command, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if !uri.IsAbs() {
		err = fmt.Errorf("Not a valid url")
		return nil, err
	}
	return &command{
		uri:    uri,
		rawurl: rawurl,
		done:   make(chan struct{}, 1),
	}, nil
}

func (c *command) URL() *url.URL {
	return c.uri
}

func (c *command) Method() string {
	return "GET"
}

func (c *command) Done() {
	c.done <- struct{}{}
}

func (c *command) Wait() {
	<-c.done
}
