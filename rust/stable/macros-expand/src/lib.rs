use macros::bind;
use wasm_bindgen::prelude::wasm_bindgen;

#[bind]
pub enum Ast {
    FromStruct(StructTemp),
    Number(u32),
}

#[bind]
impl Ast {
    #[bind]
    pub fn to_num(&self) -> u32 {
        let Self::Number(num) = self  else {
            panic("not a number");
        };

        num
    }

    #[bind]
    pub fn from_num(num: u32) -> Self {
        Self::Number(num)
    }
}

#[bind]
#[derive(Debug)]
pub struct StructTemp {}

#[bind]
impl StructTemp {
    #[bind]
    pub fn print_self(&self) {
        println!("{self:?}");
    }
}
