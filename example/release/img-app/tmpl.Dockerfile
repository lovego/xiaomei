FROM goxiaomei/appserver

COPY {{ .ProName }} ./
COPY config  ./config
COPY views   ./views

CMD [ "./{{ .ProName }}" ]

