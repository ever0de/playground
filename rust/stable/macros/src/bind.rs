pub mod item;

use proc_macro::TokenStream;
use quote::quote_spanned;
use syn::{parse_macro_input, spanned::Spanned, Item};

use crate::bind::item::{item_enum, item_impl, item_struct};

// enum Temp {}
//
// impl Temp { #[bind] fn temp() -> u32 { 42 } }
pub fn bind(_item: TokenStream, attr: TokenStream) -> TokenStream {
    let item = parse_macro_input!(attr as Item);

    match &item {
        Item::Impl(attr) => item_impl(attr),
        Item::Struct(attr) => item_struct(attr),
        Item::Enum(attr) => item_enum(attr),
        // TODO:
        // Item::Fn(attr) => {
        //     quote! {}.into()
        // }
        _ => quote_spanned! {item.span()=>
            compile_error!("`#[bind]` can only be used on `impl`, `struct`, `enum` or `fn`");
        }
        .into(),
    }
}
