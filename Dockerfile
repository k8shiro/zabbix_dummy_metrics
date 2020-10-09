FROM golang:latest

RUN go get -u github.com/AlekSi/zabbix \
              github.com/AlekSi/zabbix-sender
RUN mkdir /zabbix_dummy_metrics
WORKDIR /zabbix_dummy_metrics
