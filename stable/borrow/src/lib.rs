use std::borrow::Borrow;

// 내부가 Copy구현을 하는 타입
struct Wrapper(i32);

impl Borrow<i32> for Wrapper {
    fn borrow(&self) -> &i32 {
        &self.0
    }
}
