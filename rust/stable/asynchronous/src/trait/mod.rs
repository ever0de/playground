pub mod async_trait;

use std::{future::Future, pin::Pin};

pub trait Constructor {
    fn new() -> Pin<Box<dyn Future<Output = Self>>>;
}

struct Foo<'a> {
    _marker: std::marker::PhantomData<&'a ()>,
}

impl<'a> Constructor for Foo<'a> {
    fn new() -> Pin<Box<dyn Future<Output = Self>>> {
        Box::pin(async {
            Foo {
                _marker: std::marker::PhantomData,
            }
        })
    }
}

struct Temp<'a> {
    inner: &'a i64,
}

impl<'a> Temp<'a> {
    // Self -> don't work
    // `impl Trait` return type cannot contain a projection or `Self` that references lifetimes from a parent scope
    // see issue #103532 <https://github.com/rust-lang/rust/issues/103532> for more information
    fn new(inner: &'a i64) -> impl Future<Output = Temp<'a>> {
        async move { Self { inner } }
    }
}
