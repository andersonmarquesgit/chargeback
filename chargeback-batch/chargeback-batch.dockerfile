FROM alpine:latest

RUN mkdir /app

COPY chargeback-batch /app

CMD [ "/app/chargeback-batch" ]