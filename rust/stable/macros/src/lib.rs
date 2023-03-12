use proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, Item, Type};

#[proc_macro_attribute]
pub fn bind(item: TokenStream, attr: TokenStream) -> TokenStream {
    // enum Temp {}
    println!("item: {}", item);
    // impl Temp { #[bind] fn temp() -> u32 { 42 } }
    println!("attr: {}", attr);

    let attr = parse_macro_input!(attr as Item);
    println!("{attr:#?}");

    match attr {
        Item::Impl(attr) => {
            println!("impl:\n{attr:#?}");

            let Type::Path(self_ty) = *attr.self_ty else {
                return quote! {
                    compile_error!("unsupported `impl Type`");
                }
                .into()
            };

            println!("self_ty: {self_ty:#?}");
        }
        Item::Struct(attr) => {
            println!("struct:\n{attr:#?}");
        }
        Item::Enum(attr) => {
            println!("enum:\n{attr:#?}");
        }
        _ => {
            return quote! {
                compile_error!("unsupported item");
            }
            .into()
        }
    };

    let token = quote! {
        // fn answer() -> u32 { 42 }
    };

    println!("output: {}\n", token);

    token.into()
}
