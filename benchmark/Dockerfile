FROM alpine

RUN apk add jq util-linux

ADD https://github.com/tsenart/vegeta/releases/download/v12.8.3/vegeta-12.8.3-linux-amd64.tar.gz /vegeta.tar.gz
RUN tar xzvf vegeta.tar.gz
RUN chmod +x vegeta
RUN mv vegeta /usr/local/bin/vegeta

RUN ls
#RUN ln -s vegeta /usr/local/bin/vegeta

ADD tryme.sh ./
RUN chmod +x tryme.sh

RUN ls

CMD ["sh", "tryme.sh"]
