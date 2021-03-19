# blogit dockerfile

FROM scratch

COPY ./blogit /

ENTRYPOINT ['/blogit']
