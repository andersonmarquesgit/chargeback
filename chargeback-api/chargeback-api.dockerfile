FROM alpine:latest

RUN mkdir /app

COPY chargeback-api /app

CMD [ "/app/chargeback-api" ]