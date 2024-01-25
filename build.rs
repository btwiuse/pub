use std::env;
use std::fs::File;
use std::io::copy;
use std::io::BufReader;
use std::path::PathBuf;
use std::process::Command;

use tar::Archive;
use ureq::get;
use xz2::read::XzDecoder;

fn main() {
    println!("cargo:rerun-if-env-changed=OFFLINE");

    let os = env::var("CARGO_CFG_TARGET_OS").unwrap();
    let arch = env::var("CARGO_CFG_TARGET_ARCH").unwrap();
    let version = env::var("CARGO_PKG_VERSION").unwrap();
    let pwd = PathBuf::from(env::var("CARGO_MANIFEST_DIR").unwrap());
    let out_dir = PathBuf::from(env::var("OUT_DIR").unwrap());
    let staticlib_dir = pwd.join("staticlib").join(os.clone()).join(arch.clone());

    if dbg!(env::var("IS_STUB")) == Ok(String::from("1")) {
        return;
    }

    // use local go static lib if found
    if dbg!(staticlib_dir.exists()) {
        println!("(0) using local ./staticlib");
        println!("cargo:rustc-link-search=native={}", staticlib_dir.display());
        return;
    }

    // build stub lib as fallback
    let stub_dir = dbg!(must_build_stub(&out_dir));

    // use offline stub and avoid downloading precompile if OFFLINE=1
    if dbg!(env::var("OFFLINE")) == Ok(String::from("1")) {
        println!("(1) using offline stub");
        println!("cargo:rustc-link-search=native={}", stub_dir.display());
        return;
    }

    // download precompile if network access is enabled, fallback to stub_dir on error
    if let Ok(download_dir) = download_and_extract_precompile(&out_dir, &version, &os, &arch) {
        println!("(2) using online precompile");
        println!("cargo:rustc-link-search=native={}", download_dir.display());
    } else {
        println!("(3) failed to get precompile for {version}, using offline stub");
        println!("cargo:rustc-link-search=native={}", stub_dir.display());
    }
}

fn must_build_stub(out_dir: &PathBuf) -> PathBuf {
    env::set_var("IS_STUB", "1");
    let stub_dir = out_dir.join("stub");
    let _ = Command::new("cargo")
        .args(&[
            "build",
            "--release",
            "--lib",
            "--target-dir",
            dbg!(stub_dir.to_str().unwrap()),
        ])
        .output()
        .expect("Failed to execute command");
    let stub_release_dir = stub_dir.join("release");
    assert!(stub_release_dir.exists());
    stub_release_dir
}

fn download_and_extract_precompile(
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
    let mut response = get(&url)
        .call()
        .or_else(|e| {
            Err(std::io::Error::new(
                std::io::ErrorKind::Other,
                e.to_string(),
            ))
        })?
        .into_reader();
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
