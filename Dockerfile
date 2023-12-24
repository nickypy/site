FROM golang:1.21 as builder

WORKDIR /usr/src/site

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./ ./

RUN go build -v -o /usr/local/bin/site
RUN ["site", "build"]

FROM nginx
COPY ./config/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /usr/src/site/dist /usr/share/nginx/html
