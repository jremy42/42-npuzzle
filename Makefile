all:
	cd srcs && go build -o ..

clean:
	rm -f npuzzle

re: clean all
