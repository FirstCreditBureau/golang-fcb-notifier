FROM alpine:3.11

RUN echo -e "http://nl.alpinelinux.org/alpine/v3.11/main\nhttp://nl.alpinelinux.org/alpine/v3.11/community" > /etc/apk/repositories

RUN apk add --update --no-cache tzdata ca-certificates && update-ca-certificates
ENV TZ=Asia/Almaty
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

EXPOSE 8080
COPY ./bin/golang-fcb-notifier /golang-fcb-notifier
ENTRYPOINT ["/golang-fcb-notifier"]