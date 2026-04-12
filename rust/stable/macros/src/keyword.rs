pub const MACRO_NAME: &str = "bind";

pub fn to_private_name(text: impl AsRef<str>) -> String {
    format!("__private_{}_bind", text.as_ref())
}
