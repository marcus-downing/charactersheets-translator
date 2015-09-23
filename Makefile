run:
	./run-translator.sh

all:
	go build src/translator.go

log:
	tailf /var/log/translator
