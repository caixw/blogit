# blogit dockerfile

FROM scratch

COPY ./dist/blogit_linux_amd64/blogit ./blogit

ENTRYPOINT ['./blogit']
