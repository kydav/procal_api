# Build stage
FROM golang:1.24.5 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY entity/ ./entity
COPY routes/ ./routes
COPY services/ ./services
COPY wrappers/ ./wrappers
COPY repository/ ./repository
COPY procal.env ./procal.env
RUN go build -o godocker

# Deployment stage
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/godocker ./godocker
COPY --from=build /app/procal.env ./procal.env

EXPOSE 8000
CMD ["./godocker"]
