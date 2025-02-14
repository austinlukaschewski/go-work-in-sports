FROM arm64v8/golang:1.23 AS build

WORKDIR /app

COPY . .

RUN go mod download && go mod verify
RUN ["templ", "generate"]
RUN chmod +x ./tailwindcss
RUN ./tailwindcss -i ./static/styles/tailwind.css -o ./static/public/css/tailwind.css


RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main .

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates && update-ca-certificates

COPY ./src /app/src
COPY ./internal /app/internal
COPY ./static /app/static
COPY --from=build /app/main .

ENV PORT=4000
EXPOSE 4000
USER 1000

CMD ["./main"]