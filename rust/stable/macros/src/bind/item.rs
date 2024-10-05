use proc_macro::TokenStream;
use quote::quote;
use syn::{ImplItem, ItemEnum, ItemImpl, ItemStruct, Type};

use crate::keyword::MACRO_NAME;

///
pub fn item_impl(item: &ItemImpl) -> TokenStream {
    let mut item = item.clone();

    // self_ty is enum name
    // e.g) impl `Ast` { ... }
    let Type::Path(_self_ty) = &*item.self_ty else {
		return quote! {compile_error!("unsupported `impl Type`");}.into()
	};

    // TODO: for loop
    let item = &item.items[0];
    // e.g)
    // #[bind]
    // fn to_num(&self) -> 32 {...}
    let ImplItem::Method(mut impl_attr) = item.clone() else {
		return quote! {compile_error!("unsupported item, only supported `ImplItem::Method`");}.into();
	};

    // e.g)
    // #[bind]
    let fn_attr = impl_attr.attrs.get_mut(0);
    let Some(is_bind)= fn_attr
        .and_then(|attr| {
            attr.path
                .segments
                .iter_mut()
                .find(|segement| segement.ident == MACRO_NAME)
        }) else {
	        return quote! {#item}.into();
		};

    // remove #[bind] attribute
    impl_attr.attrs = Vec::new();
    let result = quote! {
        #[wasm_bindgen]
        #item
    };

    dbg!(result.to_string());

    quote! {#item}.into()
}

pub fn item_struct(attr: &ItemStruct) -> TokenStream {
    quote! {
        #[wasm_bindgen]
        #attr
    }
    .into()
}

// wrap struct with js_name
pub fn item_enum(attr: &ItemEnum) -> TokenStream {
    let enum_ident = &attr.ident;
    let name = crate::keyword::to_private_name(enum_ident.to_string());
    let struct_ident = syn::Ident::new(&name, enum_ident.span());

    quote! {
        #attr

        #[allow(non_camel_case_types)]
        #[wasm_bindgen(js_name = #enum_ident)]
        pub struct #struct_ident(#enum_ident);
    }
    .into()
}
