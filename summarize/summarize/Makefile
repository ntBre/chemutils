install:
	go install .

summarize: *.go
	go build .

deploy: summarize
	scp -C summarize 'woods:Programs/summarize/.'
	scp -C summarize.1 'woods:Programs/pbqff/docs/man1/.'

eland: summarize
	scp -C summarize 'eland:programs/summarize/.'

man_test: summarize.1
	groff -Tpdf -man $< > /tmp/summarize.pdf
