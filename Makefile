PACKAGE = nugo.exe

RM = del

all: $(PACKAGE)

$(PACKAGE): config.go main.go manifest.go package.go repo.go
	go build -x -o $(PACKAGE)

test: $(PACKAGE) *_test.go
	go test -v

clean:
	$(RM) $(PACKAGE)
	go clean -x

run: $(PACKAGE)
	./$(PACKAGE)

get-deps:
	go get -u github.com/gorilla/mux
	go get -u github.com/codegangsta/negroni
