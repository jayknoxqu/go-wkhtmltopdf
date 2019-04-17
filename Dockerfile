# builder image
FROM surnet/alpine-wkhtmltopdf:3.8-0.12.5-full as builder

# Image
FROM golang:1.11-alpine3.8

# Install needed packages
RUN  echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.8/main" > /etc/apk/repositories \
     && echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.8/community" >> /etc/apk/repositories \
     && apk update && apk add --no-cache \
      libstdc++ \
      libx11 \
      libxrender \
      libxext \
      libssl1.0 \
      ca-certificates \
      fontconfig \
      freetype \
      ttf-dejavu \
      ttf-droid \
      ttf-freefont \
      ttf-liberation \
      ttf-ubuntu-font-family \
    && apk add --no-cache --virtual .build-deps \
      msttcorefonts-installer \
    \
    # Install microsoft fonts
    && update-ms-fonts \
    && fc-cache -f \
    \
    # Clean up when done
    && rm -rf /var/cache/apk/* \
    && rm -rf /tmp/* \
    && apk del .build-deps

COPY --from=builder /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=builder /bin/wkhtmltoimage /bin/wkhtmltoimage

WORKDIR /go/src/app

COPY src/ .

COPY fonts/ /usr/share/fonts

RUN go get -d -v ./... && go install -v ./...

EXPOSE 80

CMD [ "app" ]
