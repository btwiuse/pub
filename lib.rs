use libc::c_char;
use std::ffi::CString;

#[link(name = "pub", kind = "static")]
extern "C" {
    fn Run(cstrs: *const c_char) -> libc::c_int;
}

pub fn run(arg: &str) -> i32 {
    let cstr = CString::new(arg).expect("CString::new failed");
    let result = unsafe { Run(cstr.as_ptr()) };
    result as i32
}
