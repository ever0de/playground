use macros::bind;
use wasm_bindgen::prelude::wasm_bindgen;

#[bind]
pub enum EnumTemp {}

// #[bind]
// impl EnumTemp {
//     #[bind]
//     pub fn temp() -> u32 {
//         42
//     }
// }

// #[bind]
// pub struct StructTemp {}
