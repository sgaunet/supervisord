FROM scratch
COPY supervisord /usr/local/bin/supervisord
ENTRYPOINT ["/usr/local/bin/supervisord"]
