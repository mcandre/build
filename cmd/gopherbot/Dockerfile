# Copyright 2017 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
FROM golang:1.10 AS build
LABEL maintainer "golang-dev@googlegroups.com"

# BEGIN deps (run `make update-deps` to update)

# Repo cloud.google.com/go at ba25346 (2017-10-13)
ENV REV=ba2534604a6c2b8854d680e50f9132cade92d02f
RUN go get -d cloud.google.com/go/compute/metadata &&\
    (cd /go/src/cloud.google.com/go && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo github.com/golang/protobuf at 130e6b0 (2017-09-20)
ENV REV=130e6b02ab059e7b717a096f397c5b60111cae74
RUN go get -d github.com/golang/protobuf/proto `#and 5 other pkgs` &&\
    (cd /go/src/github.com/golang/protobuf && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo github.com/google/go-github at 471e67c (2017-11-16)
ENV REV=471e67c93cdd36a5f7c9bba5d6d11665ba12147d
RUN go get -d github.com/google/go-github/github &&\
    (cd /go/src/github.com/google/go-github && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo github.com/google/go-querystring at 53e6ce1 (2017-01-11)
ENV REV=53e6ce116135b80d037921a7fdd5138cf32d7a8a
RUN go get -d github.com/google/go-querystring/query &&\
    (cd /go/src/github.com/google/go-querystring && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo github.com/gregjones/httpcache at 316c5e0 (2017-09-20)
ENV REV=316c5e0ff04e0ce21cfae9fbbb2ab39c7fc74e71
RUN go get -d github.com/gregjones/httpcache &&\
    (cd /go/src/github.com/gregjones/httpcache && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo go4.org at 034d17a (2017-05-25)
ENV REV=034d17a462f7b2dcd1a4a73553ec5357ff6e6c6e
RUN go get -d go4.org/types &&\
    (cd /go/src/go4.org && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo golang.org/x/net at 01c1902 (2017-11-02)
ENV REV=01c190206fbdffa42f334f4b2bf2220f50e64920
RUN go get -d golang.org/x/net/context `#and 2 other pkgs` &&\
    (cd /go/src/golang.org/x/net && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo golang.org/x/oauth2 at bb50c06 (2017-09-28)
ENV REV=bb50c06baba3d0c76f9d125c0719093e315b5b44
RUN go get -d golang.org/x/oauth2 `#and 2 other pkgs` &&\
    (cd /go/src/golang.org/x/oauth2 && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo golang.org/x/sync at f52d181 (2017-05-17)
ENV REV=f52d1811a62927559de87708c8913c1650ce4f26
RUN go get -d golang.org/x/sync/errgroup &&\
    (cd /go/src/golang.org/x/sync && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Repo golang.org/x/time at 26559e0 (2018-03-14)
ENV REV=26559e0f760e39c24d730d3224364aef164ee23f
RUN go get -d golang.org/x/time/rate &&\
    (cd /go/src/golang.org/x/time && (git cat-file -t $REV 2>/dev/null || git fetch -q origin $REV) && git reset --hard $REV)

# Optimization to speed up iterative development, not necessary for correctness:
RUN go install cloud.google.com/go/compute/metadata \
	github.com/golang/protobuf/proto \
	github.com/golang/protobuf/ptypes \
	github.com/golang/protobuf/ptypes/any \
	github.com/golang/protobuf/ptypes/duration \
	github.com/golang/protobuf/ptypes/timestamp \
	github.com/google/go-github/github \
	github.com/google/go-querystring/query \
	github.com/gregjones/httpcache \
	go4.org/types \
	golang.org/x/net/context \
	golang.org/x/net/context/ctxhttp \
	golang.org/x/oauth2 \
	golang.org/x/oauth2/internal \
	golang.org/x/sync/errgroup \
	golang.org/x/time/rate
# END deps

COPY . /go/src/golang.org/x/build/

RUN go install -ldflags "-linkmode=external -extldflags '-static -pthread'" golang.org/x/build/cmd/gopherbot


FROM scratch
LABEL maintainer "golang-dev@googlegroups.com"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/gopherbot /
ENTRYPOINT ["/gopherbot"]
