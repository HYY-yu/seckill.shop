FROM alpine:3.13 AS final

ARG serviceName
# The main mirrors are giving us timeout issues on builds periodically.
# So we can try these.
ENV TZ=Asia/Shanghai
# RUN sed -e 's/dl-cdn[.]alpinelinux.org/nl.alpinelinux.org/g' -i~ /etc/apk/repositories
RUN echo "http://mirrors.aliyun.com/alpine/v3.4/main/" > /etc/apk/repositories \
    && apk --no-cache add tzdata zeromq \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo '$TZ' > /etc/timezone

WORKDIR /app

# Copy the controller binary from the context into our
# container image
COPY sk_$serviceName /app/sk_$serviceName
COPY ./config/cfg.toml /app/cfg.toml

ENTRYPOINT ["/app/sk_sk"]