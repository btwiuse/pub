#[no_mangle]
pub extern "C" fn Run(_args: *mut std::os::raw::c_char) -> std::os::raw::c_int {
    println!("Please enable network during cargo build and make sure your os/arch is supported on https://github.com/btwiuse/pub/releases/latest/");
    0
}
