FROM 192.168.202.12:5000/goxiaomei/appserver

COPY {{ .ProName }} ./
COPY config  ./config
COPY views   ./views

CMD [ "./{{ .ProName }}" ]

