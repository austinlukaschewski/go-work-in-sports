FROM golang:1.23-alpine AS build

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY . .

# RUN source ./tools/postinstall.sh
# RUN mage buildtailwindcss
# RUN mage buildtempl
# RUN mage buildscss

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./main .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main /app
COPY --from=build /app/static /app/static
COPY --from=build /app/internal /app/internal

EXPOSE 4000
USER 1000

ENV PORT=4000

CMD ["./main"]