use r#pub::run;
use std::env;
use std::process;

fn main() {
    let args: Vec<String> = env::args().collect();
    let default = String::from(".");
    let arg = args.get(1).unwrap_or(&default);
    let result = run(arg);
    process::exit(result as i32);
}
