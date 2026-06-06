FROM ghcr.io/pnpm/pnpm:latest AS frontend-build

WORKDIR /app

COPY ./frontend ./frontend
COPY ./pnpm-lock.yaml ./pnpm-workspace.yaml ./package.json ./

RUN pnpm runtime set node 24 -g
RUN pnpm install --frozen-lockfile
RUN pnpm run build:frontend

FROM golang:1.26.4 AS backend-build

WORKDIR /app

COPY . .

COPY --from=frontend-build /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o zene .

FROM gcr.io/distroless/static-debian12

COPY --from=backend-build /app/zene .

EXPOSE 8080

CMD ["./zene"]