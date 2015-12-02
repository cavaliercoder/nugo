ifeq ($(OS),Windows_NT)
	PACKAGE = nugo.exe
	PACKAGE_ENV = set CGO_ENABLED=0 &&
	RM = del
else
	PACKAGE = nugo
	PACKAGE_ENV = "CGO_ENABLED=0"
	RM = rm -f
endif

PACKAGE_LDFLAGS = -ldflags="-s"

all: $(PACKAGE)

$(PACKAGE): config.go log.go main.go manifest.go middleware.go package.go repo.go version.go
	$(PACKAGE_ENV) go build -x $(PACKAGE_LDFLAGS) -o $(PACKAGE)

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

