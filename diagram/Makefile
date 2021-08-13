run : main.js
	go run . -web tests/c2h4.png

main.js : src/Main.elm
	elm make src/Main.elm --output=main.js

test:
	go run . -grid 16,16 tests/c2h4.cap tests/c2h4.png

help:
	go run . -h

install:
	go install .

react:
	elm reactor
