FROM registry.supos.ai/library/centos:latest

LABEL Author="rangh"

RUN mkdir -p /var/app

WORKDIR /var/app

ADD magicProxy /var/app/magicProxy
ADD ks.yaml /var/app/ks.yaml
ADD bootstrap.sh /var/app/bootstrap.sh

RUN chmod 777 /var/app/magicProxy
RUN chmod 777 /var/app/bootstrap.sh
ENTRYPOINT ["/var/app/bootstrap.sh"]

EXPOSE 3306
