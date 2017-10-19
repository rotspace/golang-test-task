package tagc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/PuerkitoBio/fetchbot"
	"golang.org/x/net/html"
)

type Tagc struct {
	fb    *fetchbot.Fetcher
	queue *fetchbot.Queue
	port  int
}

func New(opts ...func(*Tagc)) *Tagc {
	t := &Tagc{}
	for _, fn := range opts {
		fn(t)
	}
	t.initFetcher()
	return t
}

func (t *Tagc) Run() {
	if t.port == 0 {
		t.port = 8080
	}

	portStr := fmt.Sprintf(":%d", t.port)
	if err := http.ListenAndServe(portStr, t); err != nil {
		log.Fatalln(err)
	}
}

func (t *Tagc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var (
		urls      []string
		dec       = json.NewDecoder(req.Body)
		urlsCount int
		cmds      []*command
		result    []response
	)

	err := dec.Decode(&urls)
	if err != nil {
		fmt.Fprintf(w, "{error:%s}", err)
		return
	}

	{
		urlsCount = len(urls)
		cmds = make([]*command, urlsCount)
	}

	for i, url := range urls {
		cmd, err := newCmd(url)
		if err != nil {
			fmt.Fprintf(w, `{error:"%s"}`, err)
			return
		}
		cmds[i] = cmd
		t.queue.Send(cmd)
	}

	for _, cmd := range cmds {
		cmd.Wait()
		if cmd.err != nil {
			fmt.Fprintf(w, `{error:"%s"}`, err)
			return
		}
		result = append(result, cmd.result)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("  ", "  ")
	enc.Encode(result)
}

func CountTags(rdr io.Reader, chunked bool) (dup map[string]int, contentLength int, err error) {

	var (
		tokenizer *html.Tokenizer
		sCounter  *sizeCounter
	)

	if chunked {
		sCounter = newSizeCounter(httputil.NewChunkedReader(rdr))
	} else {
		sCounter = newSizeCounter(rdr)
	}
	dup = make(map[string]int)
	tokenizer = html.NewTokenizer(sCounter)

	for {
		tt := tokenizer.Next()
		if tt == html.StartTagToken ||
			tt == html.SelfClosingTagToken {
			tn, _ := tokenizer.TagName()
			dup[string(tn)]++
		} else if tt == html.ErrorToken {
			if tokenizer.Err() != io.EOF {
				err = tokenizer.Err()
			}
			contentLength = sCounter.readed
			return
		}
	}
}
