#![doc = include_str!(concat!(env!("CARGO_MANIFEST_DIR"), "/README.md"))]

use std::os::raw::{c_char, c_int};

#[no_mangle]
pub extern "C" fn Run(_args: *const c_char) -> c_int {
    eprintln!("{}", include_str!("error_msg.txt"));
    1
}
