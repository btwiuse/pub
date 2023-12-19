use std::env;
use std::path::PathBuf;

fn main() {
    let dir = env::var("CARGO_MANIFEST_DIR").unwrap();
    let lib_dir = PathBuf::from(dir).join(".");

    println!("cargo:rustc-link-search=native={}", lib_dir.display());
}
