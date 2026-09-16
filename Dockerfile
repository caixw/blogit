# blogit dockerfile

FROM scratch

LABEL org.opencontainers.image.authors="caixw <https://caixw.io>"

COPY ./blogit /

ENTRYPOINT ["/blogit"]
