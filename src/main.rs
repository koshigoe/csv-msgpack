extern crate csv;
extern crate rmp_serde as rmps;
extern crate serde;

use rmps::Serializer;
use serde::Serialize;
use std::error::Error;
use std::io::{self, Write};
use std::process;

type Row = Vec<String>;

fn run() -> Result<(), Box<Error>> {
    let mut rdr = csv::ReaderBuilder::new()
        .has_headers(false)
        .flexible(true)
        .from_reader(io::stdin());

    for result in rdr.deserialize() {
        let row: Row = result?;
        let mut buf = Vec::new();
        row.serialize(&mut Serializer::new(&mut buf)).unwrap();
        io::stdout().write(&buf)?;
    }
    Ok(())
}

fn main() {
    if let Err(err) = run() {
        eprintln!("{}", err);
        process::exit(1);
    }
}
