SRC=src/events.jsx

all: clean build/bundle.js

build/bundle.js: $(SRC)
	mkdir -p build
	browserify $^ -t [ babelify --presets [ react ] ] -o $@

.PHONY: clean dist watch
clean:
	rm -rf build
