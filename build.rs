use std::env;
use std::fs::File;
use std::io::copy;
use std::io::BufReader;
use std::path::PathBuf;
use std::process::Command;

use reqwest::blocking::get;
use tar::Archive;
use xz2::read::XzDecoder;

fn main() {
    let os = env::var("CARGO_CFG_TARGET_OS").unwrap();
    let arch = env::var("CARGO_CFG_TARGET_ARCH").unwrap();
    let version = env::var("CARGO_PKG_VERSION").unwrap();

    let dir = PathBuf::from(env::var("CARGO_MANIFEST_DIR").unwrap());
    let lib_dir = dir.join("staticlib").join(os.clone()).join(arch.clone());

    // use local go static lib if found
    if lib_dir.exists() {
        println!("cargo:rustc-link-search=native={}", lib_dir.display());
        return;
    }

    // use offline rust static lib as fallback
    let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());
    let target_dir = out_dir.join("offline");
    // run cargo b -p offline -r --target-dir target_dir
    let _ = Command::new("cargo")
        .args(&[
            "build",
            "--release",
            "--package",
            "offline",
            "--target-dir",
            target_dir.to_str().unwrap(),
        ])
        .output()
        .expect("Failed to execute command");
    let mut lib_dir = target_dir.join("release");

    if let Ok(d) = download_and_extract_staticlib(&out_dir, &version, &os, &arch) {
        lib_dir = d;
    }

    println!("cargo:rustc-link-search=native={}", lib_dir.display());
}

fn download_and_extract_staticlib(
    out_dir: &PathBuf,
    version: &str,
    os: &str,
    arch: &str,
) -> std::io::Result<PathBuf> {
    let url = format!(
        "https://github.com/btwiuse/pub/releases/download/v{}/staticlib-{}-{}.txz",
        version, os, arch,
    );

    // Download the file if network is available during build
    let mut response = match get(&url) {
        Ok(response) => response,
        Err(e) => {
            return Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        }
    };
    let mut out = File::create(out_dir.join("staticlib.txz"))?;
    copy(&mut response, &mut out)?;

    let tar_xz = File::open(out_dir.join("staticlib.txz"))?;

    // Decompress and extract the file
    let tar = XzDecoder::new(BufReader::new(tar_xz));
    let mut archive = Archive::new(tar);
    archive.unpack(&out_dir)?;

    let lib_dir = out_dir.join("staticlib").join(os).join(arch);
    if !lib_dir.exists() {
        panic!("unsupported target: {}-{}", os, arch);
    }

    Ok(lib_dir)
}
