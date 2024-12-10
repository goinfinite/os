FROM docker.io/bitnami/minideb:bullseye-amd64

WORKDIR /infinite

RUN apt-get update && apt-get upgrade -y \
    && install_packages bind9-dnsutils build-essential ca-certificates certbot cron \
    curl debian-archive-keyring git gnupg2 haveged lsb-release procps rsync supervisor \
    tar unzip vim wget zip unattended-upgrades

RUN curl -skL "https://nginx.org/keys/nginx_signing.key" | gpg --dearmor >"/usr/share/keyrings/nginx-archive-keyring.gpg" \
    && echo "deb [signed-by=/usr/share/keyrings/nginx-archive-keyring.gpg] http://nginx.org/packages/debian $(lsb_release -cs) nginx" >"/etc/apt/sources.list.d/nginx.list" \
    && install_packages nginx \
    && mkdir -p /app/logs/cron /app/logs/nginx /app/conf/pki /app/html \
    && chown -R nobody:nogroup /app

RUN curl -skL "https://mise.run" | sh \
    && mv /root/.local/bin/mise /usr/bin/mise \
    && chmod +x /usr/bin/mise \
    && echo 'eval "$(/usr/bin/mise activate bash)"' >>/etc/profile

COPY /container/nginx/root/* /etc/nginx/

COPY --chown=nobody:nogroup /container/nginx/user/ /app/conf/nginx/

COPY /container/supervisord.conf /infinite/supervisord.conf

COPY /bin/os /infinite/os

RUN chmod +x /infinite/os \
    && ln -s /infinite/os /usr/bin/os

EXPOSE 22/tcp 80/tcp 443/tcp 3306/tcp 5432/tcp 6379/tcp 1618/tcp

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/infinite/supervisord.conf"]
