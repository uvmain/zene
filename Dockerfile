FROM node:22-alpine AS frontend-build

WORKDIR /frontend

COPY ./frontend .

RUN npm install

RUN npm run build

FROM golang:1.24.2 AS backend-build

WORKDIR /app

COPY . .

COPY --from=frontend-build /frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 go build -o zene .

FROM gcr.io/distroless/static-debian12

COPY --from=backend-build /app/zene .

EXPOSE 8080

CMD ["./zene"]