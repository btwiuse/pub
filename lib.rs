use libc::c_char;

#[link(name = "pub", kind = "static")]
extern "C" {
    fn Run(cstrs: *const c_char) -> libc::c_int;
}

pub fn run(arg: &str) -> i32 {
    let cstr = arg.as_bytes();
    let result = unsafe { Run(cstr.as_ptr() as *const c_char) };
    result as i32
}
