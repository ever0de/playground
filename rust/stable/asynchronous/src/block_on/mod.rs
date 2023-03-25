#[test]
#[should_panic(
    expected = "cannot execute `LocalPool` executor from within another executor: EnterError"
)]
fn block_on_panic() {
    use futures::executor::block_on;
    block_on(async { block_on(async {}) })
}
