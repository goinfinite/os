FROM docker.io/bitnami/minideb:bullseye-amd64

WORKDIR /speedia

RUN install_packages ca-certificates wget curl tar procps debian-archive-keyring lsb-release gnupg2 haveged \
    && curl -skL "https://nginx.org/keys/nginx_signing.key" | gpg --dearmor > "/usr/share/keyrings/nginx-archive-keyring.gpg" \ 
    && echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] http://nginx.org/packages/debian $(lsb_release -cs) nginx" > "/etc/apt/sources.list.d/nginx.list" \
    && install_packages nginx cron \
    && touch /var/spool/cron/crontabs/root \
    && mkdir -p /app/logs/nginx /app/conf/pki /app/html \
    && chown -R nobody:nogroup /app

RUN wget -qO supervisord.tar.gz https://github.com/ochinchina/supervisord/releases/download/v0.7.3/supervisord_0.7.3_Linux_64-bit.tar.gz \
    && tar -xzf supervisord.tar.gz \
    && mv supervisord_*/supervisord /usr/bin/supervisord \
    && rm -rf supervisord*

RUN curl -skL "https://rtx.pub/install.sh" | sh \
    && ln -s /root/.local/share/rtx/bin/rtx /usr/bin/rtx \
    && chmod +x /usr/bin/rtx \
    && echo 'eval "$(/usr/bin/rtx activate bash)"' >> /etc/profile

COPY /container/nginx/root/* /etc/nginx/

COPY --chown=nobody:nogroup /container/nginx/user/ /app/conf/nginx/

COPY /container/supervisord.conf /speedia/supervisord.conf

COPY /bin/sos /speedia/sos

RUN chmod +x /speedia/sos \
    && ln -s /speedia/sos /usr/bin/sos

EXPOSE 22/tcp 80/tcp 443/tcp 3306/tcp 5432/tcp 6379/tcp 1618/tcp

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/speedia/supervisord.conf"]
