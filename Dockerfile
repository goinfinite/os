FROM docker.io/bitnami/minideb:bullseye-amd64

WORKDIR /speedia

RUN install_packages ca-certificates gcc make libpam0g-dev wget tar procps

RUN wget -nv https://github.com/ochinchina/supervisord/releases/download/v0.7.3/supervisord_0.7.3_Linux_64-bit.tar.gz \
    && tar -xzf supervisord_0.7.3_Linux_64-bit.tar.gz \
    && mv supervisord_0.7.3_Linux_64-bit/supervisord /usr/bin/supervisord \
    && rm -rf supervisord_0.7.3_Linux_64-bit supervisord_0.7.3_Linux_64-bit.tar.gz

COPY /bin/sam /speedia/sam

COPY supervisord.conf /speedia/supervisord.conf

EXPOSE 22/tcp 80/tcp 443/tcp 3306/tcp 10000/tcp

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/speedia/supervisord.conf"]
