## Tags - an json service for counting html tags

[godoc](godoc.org/github.com/rotspace/tagc)

It's obedience to robots.txt rules.

## Instalation

`go get -u github.com/rotspace/tagc/cmd/tagc`

## Usage

Run server in your command:

`tagc`

Then send json array with urls:

```bash
curl -XPOST http://localhost:8080 -d '["https://google.com", "https://ya.ru", "https://lenta.ru"]'
```

## Docker

`docker build -t tagc .`

`docker run -it --rm -p 8080:8080 --name tagc-app tagc`
