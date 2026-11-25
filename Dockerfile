FROM golang:1.22-alpine AS build
WORKDIR /app
COPY server ./server
WORKDIR /app/server
RUN go mod tidy && go build -o /server .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /server /app/server
COPY web /app/web
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
