FROM btwiuse/arch:rustup as builder

WORKDIR /app

ADD . /app/

RUN cargo build --release

RUN cp ./target/release/pub /usr/local/bin/

FROM btwiuse/arch

COPY --from=builder /usr/local/bin/pub /usr/bin/pub

ENTRYPOINT ["/usr/bin/pub"]
