FROM node:24-alpine AS frontend-build

WORKDIR /frontend

COPY ./frontend .
COPY ./package-lock.json ./

RUN npm ci
RUN npm run build

FROM golang:1.26.0 AS backend-build

WORKDIR /app

COPY . .

COPY --from=frontend-build /frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o zene .

FROM gcr.io/distroless/static-debian12

COPY --from=backend-build /app/zene .

EXPOSE 8080

CMD ["./zene"]