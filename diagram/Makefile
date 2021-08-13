ifeq ($(DEBUG),1)
	GOFLAGS = -debug
endif

run : main.js
	go run . -web $(GOFLAGS) tests/c2h4.png

main.js : src/Main.elm
	sed -i 's/\t/        /g' src/Main.elm
	elm make src/Main.elm --output=main.js

test:
	go run . -grid 16,16 tests/c2h4.cap tests/c2h4.png

help:
	go run . -h

install:
	go install .

bench:
	go test . -bench 'DrawGrid'

benchprof:
	go test . -bench 'DrawGrid' -cpuprofile profiles/grid.out
