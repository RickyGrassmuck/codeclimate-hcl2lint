ARG BASE=1.13.1-alpine3.10
FROM golang:${BASE} as build

WORKDIR /usr/src/app

COPY engine.json ./engine.json.template
 RUN apk add --no-cache jq
 RUN export go_version=$(go version | cut -d ' ' -f 3) && \
     cat engine.json.template | jq '.version = .version + "/" + env.go_version' > ./engine.json

COPY codeclimate-hcl2lint.go go.mod go.sum Utils ./
RUN apk add --no-cache git
RUN go build -o codeclimate-hcl2lint .

FROM golang:${BASE}
LABEL maintainer="Ricky Grassmuck <rigrassm@gmail.com>"

RUN adduser -u 9000 -D app

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/engine.json /
COPY --from=build /usr/src/app/codeclimate-hcl2lint ./

USER app
VOLUME /code

CMD ["/usr/src/app/codeclimate-hcl2lint"]