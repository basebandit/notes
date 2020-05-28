From golang:1.14-alpine as base

WORKDIR /api

FROM aquasec/trivy:0.4.4 as trivy

RUN trivy --debug --timeout 4m golang:1.14-alpine && \
  echo "No image vulnerabilities" > result

FROM base as dev

COPY go.* ./

RUN go mod download

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go env

RUN go get github.com/go-delve/delve/cmd/dlv && \ 
  go get github.com/githubnemo/CompileDaemon

EXPOSE 8080 2345

FROM dev as test

COPY . .

RUN export CGO_ENABLED=0 && \
  go test -v ./...

FROM test as build-stage

#Build the api with "-ldflags" aka linker flags to reduce binary size
# -s = disable symbol table
# -w = disable DWARF generation
RUN go build -ldlflags "-s -w" -o notes ./cmd/notes

FROM base as prod

COPY --from=trivy result secure
COPY --from=build-stage /api/notes notes

#Create a new group and user, recursively change directory ownership then allowed the binary to be executed
RUN addgroup basebandit && adduser -D -G basebandit basebandit \
  && chown -R basebandit:basebandit /api && \
  chmod +x ./notes

#Change to a non-root user
USER basebandit

#Provide meta data about the port the container must expose
EXPOSE 8080

#Define how Docker should test the container to check that it is still working
HEALTHCHECK CMD ["wget","-q","0.0.0.0:8080"]

#Provide the default command for the production container
CMD ["./notes","start"]