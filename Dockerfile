FROM scratch
COPY supervisord /usr/local/bin/supervisord
COPY pidproxy /usr/local/bin/pidproxy
ENTRYPOINT ["/usr/local/bin/supervisord"]
