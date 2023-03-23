pub fn foo() -> impl IntoIterator<Item = impl IntoIterator> {
    [[1]]
}
