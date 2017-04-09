FROM 192.168.202.12:5000/xiaomei/baseimage

USER ubuntu

EXPOSE 3000
ENV GOPORT=3000
HEALTHCHECK --interval=3s --timeout=3s --retries=20 \
  CMD sh -c 'wget -q --spider http://localhost:$GOPORT/_alive || exit 1'

