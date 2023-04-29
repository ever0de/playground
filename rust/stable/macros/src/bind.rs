pub mod item;

use proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, Item};

use crate::bind::item::{item_enum, item_impl, item_struct};

// enum Temp {}
//
// impl Temp { #[bind] fn temp() -> u32 { 42 } }
pub fn bind(_item: TokenStream, attr: TokenStream) -> TokenStream {
    let attr = parse_macro_input!(attr as Item);

    match &attr {
        Item::Impl(attr) => item_impl(attr),
        Item::Struct(attr) => item_struct(attr),
        Item::Enum(attr) => item_enum(attr),
        _ => quote! {
            compile_error!("unsupported item");
        }
        .into(),
    }
}
