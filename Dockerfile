FROM docker.io/bitnami/minideb:bullseye-amd64

WORKDIR /speedia

RUN install_packages ca-certificates wget curl tar procps cron \
    && touch /var/spool/cron/crontabs/root

RUN wget -nv https://github.com/ochinchina/supervisord/releases/download/v0.7.3/supervisord_0.7.3_Linux_64-bit.tar.gz \
    && tar -xzf supervisord_0.7.3_Linux_64-bit.tar.gz \
    && mv supervisord_0.7.3_Linux_64-bit/supervisord /usr/bin/supervisord \
    && rm -rf supervisord_0.7.3_Linux_64-bit supervisord_0.7.3_Linux_64-bit.tar.gz

COPY supervisord.conf /speedia/supervisord.conf

COPY /bin/sos /speedia/sos

RUN chmod +x /speedia/sos \
    && ln -s /speedia/sos /usr/bin/sos

EXPOSE 22/tcp 80/tcp 443/tcp 3306/tcp 10000/tcp

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/speedia/supervisord.conf"]
