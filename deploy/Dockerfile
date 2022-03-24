FROM ccr.ccs.tencentyun.com/hyy-yu/alpine-timezone:3.15

ARG serviceName
ENV SERVICENAME=$serviceName

WORKDIR /app

# Copy the controller binary from the context into our
# container image
COPY $serviceName /app/$serviceName
COPY ./internal/service/$SERVICENAME/config/cfg.toml /app/cfg.toml

ENTRYPOINT ["/app/${SERVICENAME}"]