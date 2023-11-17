FROM docker.io/bitnami/minideb:bullseye-amd64

WORKDIR /speedia

RUN install_packages ca-certificates wget curl tar procps debian-archive-keyring lsb-release gnupg2 \
    && curl -skL "https://nginx.org/keys/nginx_signing.key" | gpg --dearmor > "/usr/share/keyrings/nginx-archive-keyring.gpg" \ 
    && echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] http://nginx.org/packages/debian $(lsb_release -cs) nginx" > "/etc/apt/sources.list.d/nginx.list" \
    && install_packages nginx cron \
    && touch /var/spool/cron/crontabs/root \
    && mkdir -p /app/logs/nginx /app/conf/nginx /app/conf/pki /app/html \
    && chown -R nobody:nogroup /app

COPY /container/nginx.conf /etc/nginx/nginx.conf

COPY --chown=nobody:nogroup /container/primary.conf /app/conf/nginx/primary.conf

RUN wget -nv https://github.com/ochinchina/supervisord/releases/download/v0.7.3/supervisord_0.7.3_Linux_64-bit.tar.gz \
    && tar -xzf supervisord_0.7.3_Linux_64-bit.tar.gz \
    && mv supervisord_0.7.3_Linux_64-bit/supervisord /usr/bin/supervisord \
    && rm -rf supervisord_0.7.3_Linux_64-bit supervisord_0.7.3_Linux_64-bit.tar.gz

COPY /container/supervisord.conf /speedia/supervisord.conf

COPY /bin/sos /speedia/sos

RUN chmod +x /speedia/sos \
    && ln -s /speedia/sos /usr/bin/sos

EXPOSE 22/tcp 80/tcp 443/tcp 3306/tcp 10000/tcp

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/speedia/supervisord.conf"]
