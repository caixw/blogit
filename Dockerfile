# blogit dockerfile

FROM scratch

COPY ./cmd/blogit/blogit ./blogit

ENTRYPOINT ['./blogit']
