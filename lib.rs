use std::os::raw::{c_char, c_int};

#[link(name = "pub", kind = "static")]
extern "C" {
    fn Run(cstrs: *const c_char) -> c_int;
}

pub fn encode(args: Vec<String>) -> Vec<u8> {
    if args.is_empty() {
        "\0".into()
    } else {
        args.join("\0")
    }
    .as_bytes()
    .into()
}

pub fn run(args: Vec<String>) -> i32 {
    unsafe { Run(encode(args).as_ptr() as *const c_char) as i32 }
}
