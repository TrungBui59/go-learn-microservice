# build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY templates /templates
COPY mailServiceApp /app

CMD [ "/app/mailServiceApp" ]
