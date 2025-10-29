FROM alpine:latest

WORKDIR /app

COPY vjitt-server /app/
COPY .env /app/

RUN chmod +x /app/vjitt-server

CMD ["./vjitt-server"]