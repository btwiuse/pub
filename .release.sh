#!/bin/bash

ROOT_DIR=$(dirname $(realpath -m ${0}))

loop(){
  cd ${ROOT_DIR}
  ls -1d staticlib/*/* | tr / ' ' | while read dir os arch; do
    bash -cv "tar -cvJ $dir/$os/$arch -f staticlib-$os-$arch.txz"
  done
}

main(){
  loop
}

main
