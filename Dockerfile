FROM node:24-alpine AS frontend-build

WORKDIR /frontend

COPY frontend/ ./
RUN npm install

COPY ./frontend .
RUN npm run build

FROM golang:1.25.5 AS backend-build

WORKDIR /app

COPY go.mod go.sum ./
COPY main.go router.go ./
COPY core/ ./core/
RUN go mod download

COPY --from=frontend-build /frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o zene .

FROM cgr.dev/chainguard/static:latest

COPY --from=backend-build /app/zene /zene

USER root
RUN mkdir -p /data && chown -R 65532:65532 /data

USER nonroot

EXPOSE 8080

CMD ["/zene"]
