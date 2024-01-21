use libc::c_char;
use std::ffi::CString;

#[link(name = "pub", kind = "static")]
extern "C" {
    fn Run(cstrs: *const c_char) -> libc::c_int;
}

fn main() {
    let strs = CString::new("https://k0s.io").unwrap();
    let result = unsafe { Run(strs.as_ptr()) };
    println!("Result: {}", result);
}
