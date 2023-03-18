#[test]
fn let_else() {
    enum Value {
        Integer(i32),
    }

    let options = Some(Value::Integer(1));

    let Some(Value::Integer(_)) = options else {
        panic!("expected an integer");
    };
}
