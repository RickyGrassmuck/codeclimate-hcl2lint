ARG BASE=1.16.2-alpine3.13
FROM golang:${BASE} as build

WORKDIR /usr/src/app

COPY engine.json ./engine.json.template
 RUN apk add --no-cache jq
 RUN export go_version=$(go version | cut -d ' ' -f 3) && \
     cat engine.json.template | jq '.version = .version + "/" + env.go_version' > ./engine.json

COPY codeclimate-hcl2lint.go engine.go go.mod go.sum  ./
RUN apk add --no-cache git
RUN go version && go build -o codeclimate-hcl2lint codeclimate-hcl2lint.go engine.go 

FROM golang:${BASE}
LABEL maintainer="Ricky Grassmuck <rigrassm@gmail.com>"

RUN adduser -u 9000 -D app

WORKDIR /usr/src/app

COPY --from=build /usr/src/app/engine.json /
COPY --from=build /usr/src/app/codeclimate-hcl2lint ./

USER app
VOLUME /code

CMD ["/usr/src/app/codeclimate-hcl2lint"]