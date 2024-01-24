use r#pub::run;
use std::env;
use std::process;

fn main() {
    process::exit(run(env::args().skip(1).collect()));
}
