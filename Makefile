ifeq ($(OS),Windows_NT)
	PACKAGE = nugo.exe
	RM = del
else
	PACKAGE = nugo
	RM = rm -f
endif

all: $(PACKAGE)

$(PACKAGE): config.go log.go main.go manifest.go middleware.go package.go repo.go version.go
	go build -x -o $(PACKAGE)

test: $(PACKAGE)
	go test -v

clean:
	$(RM) $(PACKAGE)
	go clean -x

run: $(PACKAGE)
	./$(PACKAGE)

get-deps:
	go get -u github.com/codegangsta/negroni
	go get -u gopkg.in/yaml.v2

