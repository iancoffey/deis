FROM prom/alertmanager:0.0.4

COPY rootfs /

CMD ["boot"]

USER        golang
WORKDIR     /alertmanager
EXPOSE      9093
