install:
	go install .

deploy:
	go build .
	scp -C summarize 'woods:Programs/summarize/.'
	scp -C summarize.1 'woods:Programs/pbqff/docs/man1/.'

man_test: summarize.1
	groff -Tpdf -man $< > /tmp/summarize.pdf
