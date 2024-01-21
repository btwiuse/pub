use std::env;
use std::fs::File;
use std::io::copy;
use std::io::BufReader;
use std::path::PathBuf;

use flate2::read::GzDecoder;
use reqwest::blocking::get;
use tar::Archive;

fn main() {
    let arch = env::var("CARGO_CFG_TARGET_ARCH").unwrap();
    let version = env::var("CARGO_PKG_VERSION").unwrap();

    let dir = PathBuf::from(env::var("CARGO_MANIFEST_DIR").unwrap());
    let lib_dir = dir.join("staticlib").join("linux").join(arch.clone());
    if lib_dir.exists() {
        println!("cargo:rustc-link-search=native={}", lib_dir.display());
        return
    }


    let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());
    // Download the file
    let url = format!("https://github.com/btwiuse/pub/releases/download/v{}/staticlib.tgz", version);
    let mut response = get(url).expect("Failed to download file");
    let mut out = File::create(out_dir.join("staticlib.tgz")).expect("Failed to create file");
    copy(&mut response, &mut out).expect("Failed to copy content");

    // Decompress and extract the file
    let tar_gz = File::open(out_dir.join("staticlib.tgz")).unwrap();
    let tar = GzDecoder::new(BufReader::new(tar_gz));
    let mut archive = Archive::new(tar);
    archive.unpack(&out_dir).expect("Failed to unpack archive");

    let lib_dir = out_dir.join("staticlib").join("linux").join(arch);

    println!("cargo:rustc-link-search=native={}", lib_dir.display());
}
