FROM 192.168.2.254:54800/alpine:latest

RUN mkdir -p /app-acc/configs

ENV WORKDIR /app-acc

RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.4/main\nhttp://mirrors.aliyun.com/alpine/v3.4/community" >  /etc/apk/repositories \
    && apk update && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Shanghai/Asia" > /etc/timezone \
    && apk del tzdata

WORKDIR $WORKDIR

COPY ./bin/bbx-ragflow-mcp $WORKDIR/bbx-ragflow-mcp

RUN chmod +x $WORKDIR/bbx-ragflow-mcp

EXPOSE 25003
# start
CMD ["./bbx-ragflow-mcp"]
