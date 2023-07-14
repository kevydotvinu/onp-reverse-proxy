FROM registry.fedoraproject.org/f32/golang
USER root
WORKDIR /
RUN openssl req -x509 \
                -newkey rsa:4096 \
                -nodes \
                -keyout reverse-proxy.key \
                -out reverse-proxy.crt \
                -days 365 \
                -subj "/C=IN/ST=Maharashtra/L=Mumbai/O=ONP/CN=onp.ocp.example.local"
COPY main.go go.mod .
ENTRYPOINT ["go", "run", "main.go", "-key", "/reverse-proxy.key", "-cert", "/reverse-proxy.crt"]
