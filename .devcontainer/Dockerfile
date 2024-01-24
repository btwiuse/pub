FROM btwiuse/k0s AS k0s
FROM btwiuse/ufo AS ufo

FROM btwiuse/arch:golang

RUN pacman -Syu --noconfirm --overwrite='*' aarch64-linux-gnu-gcc \
    docker clang rustup deno pkgfile yarn npm wget htop neofetch cmake gopls staticcheck delve go-tools &&\
    yes | pacman -Scc

RUN pkgfile -u

USER btwiuse

# RUN go install -v golang.org/x/tools/gopls@latest
# RUN go install -v honnef.co/go/tools/cmd/staticcheck@latest
# RUN go install -v github.com/go-delve/delve/cmd/dlv@latest
# RUN go install -v golang.org/dl/gotip@latest && gotip download
# RUN go install -v github.com/btwiuse/baize@latest

RUN echo 'set -o vi' >> $HOME/.bashrc &&\
    echo 'export EDITOR=vim' >> $HOME/.bashrc &&\
    echo 'export PATH=${HOME}/go/bin:${PATH}' >> $HOME/.bashrc

COPY --from=k0s /usr/bin/k0s /usr/bin/k0s
COPY --from=ufo /usr/bin/ufo /usr/bin/ufo
