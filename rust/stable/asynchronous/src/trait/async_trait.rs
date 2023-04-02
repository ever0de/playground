use async_trait::async_trait;

#[async_trait]
pub trait ObjectSafe: Sync {
    async fn f(&self) {}

    async fn new() -> Self
    where
        Self: Sized;

    async fn g(&mut self);
}

pub struct MyType {}

#[async_trait]
impl ObjectSafe for MyType {
    async fn new() -> Self {
        Self {}
    }

    async fn g(&mut self) {}
}

#[test]
fn object_safe() {
    let value = MyType {};
    // make trait object
    let _object = &value as &dyn ObjectSafe;
}
