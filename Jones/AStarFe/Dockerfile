FROM harbor.adtech.sogou/adtd/nginx:v0

RUN mkdir -p /search/app

WORKDIR /usr/share/nginx/html

COPY ./entrypoint.sh /
RUN chmod +x /entrypoint.sh

COPY dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/