#![allow(dead_code, unused_variables)]

#[derive(Debug, Clone, thiserror::Error)]
enum Error {
    #[error("test error: {0}")]
    Test(&'static str),
}

type Result<T, E = Error> = std::result::Result<T, E>;

fn main() -> Result<i64> {
    let result: Result<i64> = Err(Error::Test("hi"));

    let result = &result;
    let owned_ok = result.as_ref().map_err(|err| err.clone())?;

    Ok(1)
}
