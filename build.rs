use std::env;
use std::fs::File;
use std::io::copy;
use std::io::BufReader;
use std::path::PathBuf;

use reqwest::blocking::get;
use tar::Archive;
use xz2::read::XzDecoder;

fn main() {
    let arch = env::var("CARGO_CFG_TARGET_ARCH").unwrap();
    let version = env::var("CARGO_PKG_VERSION").unwrap();

    let dir = PathBuf::from(env::var("CARGO_MANIFEST_DIR").unwrap());
    let lib_dir = dir.join("staticlib").join("linux").join(arch.clone());
    if lib_dir.exists() {
        println!("cargo:rustc-link-search=native={}", lib_dir.display());
        return;
    }

    let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());

    // nop stub for offline build
    let mut tar_xz = File::open(dir.join("offline").join("staticlib.txz")).unwrap();

    let url = format!(
        "https://github.com/btwiuse/pub/releases/download/v{}/staticlib.txz",
        version
    );

    // Download the file if network is available during build
    if let Ok(mut response) = get(url) {
        let mut out = File::create(out_dir.join("staticlib.txz")).expect("Failed to create file");
        copy(&mut response, &mut out).expect("Failed to copy content");
        tar_xz = File::open(out_dir.join("staticlib.txz")).unwrap();
    }

    // Decompress and extract the file
    let tar = XzDecoder::new(BufReader::new(tar_xz));
    let mut archive = Archive::new(tar);
    archive.unpack(&out_dir).expect("Failed to unpack archive");

    let lib_dir = out_dir.join("staticlib").join("linux").join(arch);

    println!("cargo:rustc-link-search=native={}", lib_dir.display());
}
