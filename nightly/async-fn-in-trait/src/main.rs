#![feature(async_fn_in_trait)]
// https://github.com/rust-lang/rust/issues/91611
// - https://github.com/rust-lang/rust/issues/103854
// https://blog.rust-lang.org/inside-rust/2022/11/17/async-fn-in-trait-nightly.html
// - https://blog.rust-lang.org/inside-rust/2022/11/17/async-fn-in-trait-nightly.html#fn4
// https://smallcultfollowing.com/babysteps/blog/2023/02/01/async-trait-send-bounds-part-1-intro/

// Recap: How async/await works in Rust
// async fn fetch_data(db: &MyDb) -> String { ... }
// ->
// fn fetch_data<'a>(db: &'a MyDb) -> impl Future<Output = String> + 'a

use std::future::Future;

trait Database {
    async fn fetch_data(&self) -> String;
}
// # async-trait
// https://rust-lang.github.io/wg-async/vision/submitted_stories/status_quo/barbara_benchmarks_async_trait.html
// trait Database {
//     fn fetch_data<'async_trait>(
//         &'async_trait self,
//     ) -> Pin<Box<dyn Future<Output = String> + Send + 'async_trait>>;
// }

// # The historic problem of async fn in trait
trait DatabaseAssociate {
    type FetchData<'a>: Future<Output = String> + 'a
    where
        Self: 'a;
    fn fetch_data<'a>(&'a self) -> Self::FetchData<'a>;
}
// impl Database for MyDb {
//     type FetchData<'a> = /* what type goes here??? */;
//     fn fetch_data<'a>(&'a self) -> FetchData<'a> { async move { ... } }
// }

impl Database for MyDb {
    async fn fetch_data(&self) -> String {
        self.fetch_data().await
    }
}

struct MyDb;

impl MyDb {
    async fn fetch_data(&self) -> String {
        "Hello, world!".to_string()
    }
}

#[tokio::main]
async fn main() {
    let data = fetch_data(MyDb).await;
    println!("{data}");
}

async fn fetch_data(db: impl Database) -> String {
    db.fetch_data().await
}
