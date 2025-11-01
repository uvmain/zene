FROM node:22-alpine AS frontend-build

WORKDIR /frontend

RUN apk add --no-cache git openssh

COPY ./frontend .

RUN npm install
RUN npm run build

RUN git describe --tags --always > /frontend/dist/version.txt

FROM golang:1.25.1 AS backend-build

WORKDIR /app

COPY . .

COPY --from=frontend-build /frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o zene .

FROM gcr.io/distroless/static-debian12

COPY --from=backend-build /app/zene .

EXPOSE 8080

CMD ["./zene"]