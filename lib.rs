#![doc = include_str!(concat!(env!("CARGO_MANIFEST_DIR"), "/README.md"))]

use std::os::raw::{c_char, c_int};

#[no_mangle]
pub extern "C" fn Run(_args: *const c_char) -> c_int {
    eprintln!("If you are seeing this error, reason could be");
    eprintln!("");
    eprintln!("1) network is disabled during cargo build, or env OFFLINE=1 is set");
    eprintln!(
        "2) your os/arch is unsupported (see https://github.com/btwiuse/pub/releases/latest/)"
    );
    eprintln!("3) the ./staticlib/$os/$arch directory doesn't exist");
    eprintln!("");
    eprintln!("Please consider installing via `go install` instead:");
    eprintln!("");
    eprintln!("   go install github.com/btwiuse/pub/cmd/pub@latest");
    eprintln!("");
    0
}
