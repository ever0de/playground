use macros::bind;

#[bind]
pub enum Temp {}

#[bind]
impl Temp {
    // #[bind]
    pub fn temp() -> u32 {
        42
    }
}
