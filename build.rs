use std::env;
use std::path::PathBuf;

fn main() {
    let dir = env::var("CARGO_MANIFEST_DIR").unwrap();
    let arch = env::var("CARGO_CFG_TARGET_ARCH").unwrap();

    let lib_dir = PathBuf::from(dir)
        .join("staticlib")
        .join("linux")
        .join(arch);

    println!("cargo:rustc-link-search=native={}", lib_dir.display());
}
