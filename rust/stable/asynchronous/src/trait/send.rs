use async_trait::async_trait;

pub enum Payload {
    A,
}

pub enum Error {}

#[async_trait(?Send)]
pub trait OptionalSend {
    async fn send(&self) -> Result<Payload, Error>;
}

#[async_trait(?Send)]
pub trait OptionalSendMut {
    async fn send_mut(&mut self) -> Result<Payload, Error>;
}

pub mod required {
    use super::*;

    pub struct Wrapper<T> {
        sender: T,
    }

    impl<T> Wrapper<T>
    where
        T: OptionalSend + OptionalSendMut,
    {
        pub fn new(sender: T) -> Self {
            Self { sender }
        }
    }

    pub struct Sender {}

    // Failed compile with error: method `send` has an incompatible type for trait
    // Since the trait's required return type changes to Pin<Box<dyn Future<...> + "Send">>, it will not compile.

    // #[async_trait]
    // impl OptionalSend for Sender {
    //     async fn send(&self) -> Result<Payload, Error> {
    //         Ok(Payload::A)
    //     }
    // }

    // #[async_trait]
    // impl OptionalSendMut for Sender {
    //     async fn send_mut(&mut self) -> Result<Payload, Error> {
    //         Ok(Payload::A)
    //     }
    // }
}
