FROM centos:6
RUN yum -y install curl sudo
RUN sed -i 's/Defaults    requiretty/#Defaults    requiretty/' /etc/sudoers
RUN curl -fsSL https://repo.stns.jp/scripts/yum-repo.sh | sh
RUN yum -y install stns libnss-stns rsyslog
RUN sed -i \
  '-e s/passwd:     files/passwd:     files stns/g' \
  '-e s/shadow:     files/shadow:     files stns/g' \
  '-e s/group:      files/group:      files stns/g' \
  /etc/nsswitch.conf
ADD files/stns.conf /etc/stns/stns.conf
ADD files/libnss_stns.conf /etc/stns/libnss_stns.conf
CMD service rsyslog start && service stns start && /bin/bash
