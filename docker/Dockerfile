FROM scratch

ADD in/ca-certificates.pem /etc/ssl/certs/ca-certificates.pem
COPY in/a3s /a3s

EXPOSE 443

ENTRYPOINT [ "/a3s" ]
