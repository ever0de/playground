use macros::bind;

#[bind]
enum Temp {}

#[bind]
impl Temp {
    #[bind]
    fn temp() -> u32 {
        42
    }
}
