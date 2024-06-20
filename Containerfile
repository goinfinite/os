FROM docker.io/bitnami/minideb:bullseye-amd64

WORKDIR /speedia

RUN apt-get update && apt-get upgrade -y \
    && install_packages ca-certificates wget curl tar procps debian-archive-keyring lsb-release gnupg2 haveged rsync zip unzip bind9-dnsutils build-essential git certbot \
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

RUN curl -skL "https://mise.run" | sh \
    && mv /root/.local/bin/mise /usr/bin/mise \
    && chmod +x /usr/bin/mise \
    && echo 'eval "$(/usr/bin/mise activate bash)"' >> /etc/profile

COPY /container/nginx/root/* /etc/nginx/

COPY --chown=nobody:nogroup /container/nginx/user/ /app/conf/nginx/

COPY /container/supervisord.conf /speedia/supervisord.conf

COPY /bin/os /speedia/os

RUN chmod +x /speedia/os \
    && ln -s /speedia/os /usr/bin/os

EXPOSE 22/tcp 80/tcp 443/tcp 3306/tcp 5432/tcp 6379/tcp 1618/tcp

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/speedia/supervisord.conf"]
