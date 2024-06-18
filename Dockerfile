FROM golang:1.21-alpine as builder

ARG branch=develop
ENV current_branch=$branch

# Utilities
RUN apk add git gcc g++ tzdata zip ca-certificates
# Add dep for package management
# RUN go get -u -f -v github.com/golang/dep/...

#set workdir
RUN mkdir -p /go/src/turbine-app
WORKDIR /go/src/turbine-api

# COPY go.mod and go.sum files to the workspace
COPY go.mod . 
COPY go.sum .

## Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
RUN go mod tidy

COPY . .
RUN go get

# COPY ENV BASED ON BRANCH
COPY .env.$branch .env

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/turbine-api main.go 

# final stage
FROM alpine:latest
ARG branch=develop
ENV current_branch=$branch

# COPY ENV BASED ON BRANCH
COPY .env.$branch .env

#Env  
ENV TIMEZONE Asia/Jakarta

#set timezone
RUN apk --no-cache add tzdata && echo "Asia/Jakarta" > /etc/timezone
RUN apk add --update tzdata && \
cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime && \
echo "${TIMEZONE}" > /etc/timezone && apk del tzdata

#expose
EXPOSE ${PORT}
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/turbine-api .

RUN printf "#!/bin/sh\n\nwhile true; do\n\techo \"[INFO] Starting Service at \$(date)\"\n\t(./turbine-api >> ./history.log || echo \"[ERROR] Restarting Service at \$(date)\")\ndone" > run.sh
RUN printf "#!/bin/sh\n./run.sh & tail -F ./history.log" > up.sh
RUN chmod +x up.sh run.sh
CMD ["./up.sh"]