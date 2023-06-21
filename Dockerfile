FROM registry.access.redhat.com/ubi8/ubi-init:8.8

WORKDIR /speedia

RUN rpm --import https://www.centos.org/keys/RPM-GPG-KEY-CentOS-Official \
    && dnf config-manager --disableplugin subscription-manager --add-repo http://mirror.centos.org/centos/8-stream/BaseOS/x86_64/os \
    && dnf config-manager --disableplugin subscription-manager --add-repo http://mirror.centos.org/centos/8-stream/AppStream/x86_64/os \
    && dnf config-manager --disableplugin subscription-manager --add-repo http://mirror.centos.org/centos/8-stream/PowerTools/x86_64/os \
    && dnf install -qy https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm

RUN dnf -y update-minimal --security --sec-severity=Important --sec-severity=Critical \
    && dnf install --enablerepo=* -qy git pam-devel curl wget tar zip

COPY /bin/sam /speedia/sam
COPY ./sam.service /etc/systemd/system/sam.service

RUN chmod +x /speedia/sam \
    && systemctl enable sam.service

EXPOSE 80/tcp
EXPOSE 443/tcp
EXPOSE 10000/tcp

ENTRYPOINT ["/sbin/init"]
