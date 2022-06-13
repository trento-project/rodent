FROM registry.suse.com/bci/golang:latest AS build

WORKDIR /BUILD

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .
RUN go build -o /rodent


##
## Deploy
##
FROM registry.suse.com/bci/golang:latest

WORKDIR /

COPY --from=build /rodent /rodent
