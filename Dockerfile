FROM node:22-alpine AS ff

WORKDIR /

COPY ./package.json .

RUN npm install

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
COPY --from=ff ./node_modules/bin/ffprobe-baron/ffprobe ./bin/ffprobe
COPY --from=ff ./node_modules/bin/ffmpeg-baron/ffmpeg ./bin/ffmpeg

EXPOSE 8080

CMD ["./zene"]