FROM alpine:latest

RUN mkdir /app

COPY chargeback-processor /app

CMD [ "/app/chargeback-processor" ]