FROM 192.168.202.12:5000/xiaomei/appserver

WORKDIR /home/ubuntu/{{ .ProName }}

COPY {{ .ProName }} ./
COPY config  ./config
COPY views   ./views
RUN mkdir ../logs/{{ .ProName }} && ln -s ../logs/{{ .ProName }} log

CMD [ "./{{ .ProName }}" ]

