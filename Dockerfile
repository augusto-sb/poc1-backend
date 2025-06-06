FROM docker.io/library/golang:1.22.12-alpine3.21 AS compiler
WORKDIR /app/
COPY . .
ARG VERSION=1.0.0
RUN go build -o main -ldflags "-X main.version=${VERSION}" .

FROM scratch AS runner
COPY --from=compiler /app/main /main
USER 1001:1001
ENTRYPOINT ["/main"]
EXPOSE 8080/tcp