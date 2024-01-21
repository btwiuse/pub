#!/usr/bin/env bash

mkdir -p staticlib/linux/{aarch64,x86_64}

cargo build --release --target=aarch64-unknown-linux-gnu
aarch64-linux-gnu-strip ../target/aarch64-unknown-linux-gnu/release/libpub.a
cp -v ../target/aarch64-unknown-linux-gnu/release/libpub.a staticlib/linux/aarch64/

cargo build --release --target=x86_64-unknown-linux-gnu
strip ../target/x86_64-unknown-linux-gnu/release/libpub.a
cp -v ../target/x86_64-unknown-linux-gnu/release/libpub.a staticlib/linux/x86_64/

tar cvJ staticlib > staticlib.txz

du -sh staticlib.txz
