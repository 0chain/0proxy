FROM golang:1.11.4-alpine3.8 as zproxy_build

RUN apk add --update --no-cache build-base linux-headers git cmake bash perl grep

# Install Herumi's cryptography
RUN apk add gmp gmp-dev openssl-dev && \
    cd /tmp && \
    wget -O - https://github.com/herumi/mcl/archive/master.tar.gz | tar xz && \
    wget -O - https://github.com/herumi/bls/archive/master.tar.gz | tar xz && \
    mv mcl* mcl && \
    mv bls* bls && \
    make -C mcl -j $(nproc) lib/libmclbn256.so install && \
    cp mcl/lib/libmclbn256.so /usr/local/lib && \
    make -C bls -j $(nproc) install && \
    rm -R /tmp/mcl && \
    rm -R /tmp/bls

ENV SRC_DIR=/0proxy
ENV GO111MODULE=on

# Download the dependencies:
# Will be cached if we don't change mod/sum files
COPY ./code/go/0proxy.io/core/go.mod          ./code/go/0proxy.io/core/go.sum          $SRC_DIR/go/0proxy.io/core/
COPY ./code/go/0proxy.io/zproxycore/go.mod     ./code/go/0proxy.io/zproxycore/go.sum     $SRC_DIR/go/0proxy.io/zproxycore/

WORKDIR $SRC_DIR/go/0proxy.io/zproxycore/zproxy
RUN go mod download

# Add the source code:
ADD ./code/go/0proxy.io $SRC_DIR/go/0proxy.io

RUN go build -v -tags "bn256 development" -ldflags "-X 0proxy.io/core/build.BuildTag=$GIT_COMMIT"

# Copy the build artifact into a minimal runtime image:
FROM golang:1.11.4-alpine3.8
ENV APP_DIR=/0proxy
WORKDIR $APP_DIR
RUN apk add gmp gmp-dev openssl-dev
COPY --from=zproxy_build  /usr/local/lib/libmcl*.so \
                        /usr/local/lib/libbls*.so \
                        /usr/local/lib/
COPY --from=zproxy_build $APP_DIR/go/0proxy.io/zproxycore/zproxy/zproxy $APP_DIR/bin/zproxy

