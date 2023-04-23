mod bind;

use proc_macro::TokenStream;

#[proc_macro_attribute]
pub fn bind(item: TokenStream, attr: TokenStream) -> TokenStream {
    bind::bind(item, attr)
}
