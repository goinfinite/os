FROM docker.io/bitnami/minideb:bookworm-amd64

WORKDIR /speedia

RUN apt-get update \
    && apt-get install -y --only-upgrade $(apt-get -s -o Debug::NoLocking=true upgrade | awk '/^Inst.*ecurity/ {print $2}') \
    && apt-get install -y gcc make libpam0g-dev wget tar procps supervisor

COPY /bin/sam /speedia/sam

COPY supervisord.conf /speedia/supervisord.conf

EXPOSE 22/tcp 80/tcp 443/tcp 3306/tcp 10000/tcp

ENTRYPOINT ["/usr/bin/supervisord", "-c", "/speedia/supervisord.conf"]
