use proc_macro::TokenStream;
use quote::quote;
use syn::{parse_macro_input, Item, Type};

pub fn bind(item: TokenStream, attr: TokenStream) -> TokenStream {
    // enum Temp {}
    println!("item: {}", item);
    // impl Temp { #[bind] fn temp() -> u32 { 42 } }
    println!("attr: {}", attr);

    let attr = parse_macro_input!(attr as Item);
    println!("{attr:#?}");

    let output = match &attr {
        Item::Impl(attr) => {
            println!("impl:\n{attr:#?}");

            let Type::Path(self_ty) = &*attr.self_ty else {
					return quote! {
						compile_error!("unsupported `impl Type`");
					}
					.into()
				};

            for item in attr.items.iter() {
                match item {
                    syn::ImplItem::Method(attr) => {
                        println!("method:\n{attr:#?}");

                        let attr = attr.attrs.get(0);
                        let is_bind = attr
                            .map(|attr| {
                                attr.path.segments.len() == 1
                                    && attr.path.segments[0].ident == "bind"
                            })
                            .unwrap_or(false);

                        println!("method is_bind: {is_bind}");

                        return quote! {}.into();
                    }
                    _ => {
                        println!("other:\n{attr:#?}");

                        return quote! {
                            compile_error!("unsupported item");
                        }
                        .into();
                    }
                }
            }

            println!("self_ty: {self_ty:#?}");

            quote! {}
        }
        Item::Struct(attr) => {
            println!("struct:\n{attr:#?}");

            quote! {
                #[wasm_bindgen]
                #attr
            }
        }
        Item::Enum(attr) => {
            println!("enum:\n{attr:#?}");

            let enum_ident = &attr.ident;
            let name = format!("__Private_{}_Bind", enum_ident);
            let struct_ident = syn::Ident::new(&name, enum_ident.span());

            quote! {
                // enum
                #attr
                #[allow(non_camel_case_types)]
                #[wasm_bindgen(js_name = #enum_ident)]
                pub struct #struct_ident(#enum_ident);
            }
        }
        _ => {
            quote! {
                compile_error!("unsupported item");
            }
        }
    };

    println!("output: {}\n", output);

    output.into()
}
