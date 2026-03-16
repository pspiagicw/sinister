VERSION := unversioned
NAME := sinister
LDFLAG := -X main.VERSION=${VERSION}
build: 
	CGO_ENABLED=1 go build -ldflags "${LDFLAGS}" -o ${NAME}


