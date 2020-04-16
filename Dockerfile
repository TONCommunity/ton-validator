FROM ubuntu:18.04
RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update && \
    apt-get install -y openssl wget git golang-go&& \
    rm -rf /var/lib/apt/lists/*

ENV LD_LIBRARY_PATH "/root/go/src/github.com/mercuryoio/tonlib-go/v2/lib/linux"
ENV PATH $PATH:/ton/bin:/root/go/bin
ENV FIFTPATH /ton/fift/lib
RUN go get -u github.com/mercuryoio/ton-validator-bot
RUN mkdir -p /ton/bin /ton/smartcont /ton/fift/lib /ton/work
COPY --from=it4addict/ton-build /ton/build/lite-client/lite-client /ton/bin
COPY --from=it4addict/ton-build /ton/build/validator-engine-console/validator-engine-console /ton/bin
COPY --from=it4addict/ton-build /ton/build/crypto/fift /ton/bin
COPY --from=it4addict/ton-build /ton/crypto/fift/lib /ton/fift/lib
COPY --from=it4addict/ton-build /ton/crypto/smartcont /ton/smartcont
WORKDIR /ton/work
