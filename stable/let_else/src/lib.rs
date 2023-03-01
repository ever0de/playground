#[test]
fn let_else() {
    struct Value(i32);

    let options = Some(Value(1));

    let Value(value) = options else {
        return;
    };
}
