# Building Backend
FROM golang:alpine as passport-server

RUN apk add nodejs npm

WORKDIR /source
COPY . .
WORKDIR /source/pkg/view
RUN npm install
RUN npm run build
WORKDIR /source
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs -o /dist ./pkg/cmd/main.go

# Runtime
FROM golang:alpine

COPY --from=passport-server /dist /passport/server

EXPOSE 8444

CMD ["/passport/server"]