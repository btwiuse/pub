use r#pub::run;
use std::env;
use std::process;

fn main() {
    let args: Vec<String> = env::args().skip(1).collect();
    let default = String::from(".");
    let arg0 = args.get(0).unwrap_or(&default);
    let result = if dbg!(arg0 == ".") {
        run(&format!("{}\0\0", arg0))
    } else {
        run(&args.join("\0"))
    };
    process::exit(result as i32);
}
