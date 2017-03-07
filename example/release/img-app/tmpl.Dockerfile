FROM 192.168.202.12:5000/xiaomei/appserver

RUN mkdir /home/ubuntu/{{ .ProName }} /home/ubuntu/logs/{{ .ProName }}
WORKDIR /home/ubuntu/{{ .ProName }}

COPY {{ .ProName }} ./
COPY config  ./config
COPY views   ./views
RUN ln -s ../logs/{{ .ProName }} log

CMD [ "./{{ .ProName }}" ]

