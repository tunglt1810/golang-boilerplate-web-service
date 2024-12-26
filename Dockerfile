#syntax=docker/dockerfile:experimental
#-----------------------------------------------
# Builder stage.
FROM golang:1.23-bookworm AS builder
ENV GO111MODULE=on

# ARG SSH_PRIVATE_KEY

# RUN mkdir /root/.ssh/
# RUN printf "Host *\t\nStrictHostKeyChecking no\n" > /root/.ssh/config && chmod 400 /root/.ssh/config
# RUN echo "${SSH_PRIVATE_KEY}" > /root/.ssh/id_rsa && chmod 600 /root/.ssh/id_rsa
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN mkdir -p /src
WORKDIR /src

COPY . ./

RUN --mount=type=ssh,mode=741,uid=100,gid=102 cd cmd/app && CGO_ENABLED=0 GOOS=linux go build -o app

#-----------------------------------------------
# Runner stage.
FROM debian:buster AS runner

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates

RUN update-ca-certificates
COPY --from=builder /src/cmd/app/app /app/
COPY ./entrypoint.sh /app/
COPY ./config/config.yml /app/config/config.yml

ENV ENV_CONFIG_ONLY=true
WORKDIR /app
ENTRYPOINT ["bash", "./entrypoint.sh"]
