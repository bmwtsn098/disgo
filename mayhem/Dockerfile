FROM golang:1.18 as builder

COPY . /disgo
WORKDIR /disgo/mayhem

RUN go install github.com/dvyukov/go-fuzz/go-fuzz@latest github.com/dvyukov/go-fuzz/go-fuzz-build@latest
RUN go get github.com/dvyukov/go-fuzz/go-fuzz-dep
RUN apt update && apt install -y clang

RUN go-fuzz-build -libfuzzer -o fuzzDiscord.a && \
    clang -fsanitize=fuzzer fuzzDiscord.a -o fuzzDiscord.libfuzzer

FROM debian:bookworm-slim
COPY --from=builder /disgo/mayhem/fuzzDiscord.libfuzzer /