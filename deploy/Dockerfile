FROM  ccr.ccs.tencentyun.com/hyy-yu/alpine-timezone:3.15

ARG serviceName
ENV SERVICENAME=$serviceName
RUN echo $serviceName
RUN echo $SERVICENAME
WORKDIR /app

# Copy the controller binary from the context into our
# container image
COPY service /app/service
COPY internal/service/shop/config/cfg.toml /app/cfg.toml

ENTRYPOINT ["./service"]