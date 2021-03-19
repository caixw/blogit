# blogit dockerfile

FROM scratch

COPY /home/runner/work/blogit/blogit/dist/blogit_linux_amd64/blogit ./blogit

ENTRYPOINT ['./blogit']
