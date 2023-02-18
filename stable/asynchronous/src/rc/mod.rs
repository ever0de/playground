use std::rc::Rc;

struct NonClone;

async fn use_rc_in_function(i: NonClone) {
    let _rc = Rc::new(i);
}

async fn get_rc_ref_but_sync_error(i: &Rc<NonClone>) {
    let _rc = Rc::clone(i);
}

async fn get_rc_but_send_error(i: Rc<NonClone>) {
    let _rc = Rc::clone(&i);
}

async fn send_sync_test() {
    let _ok = tokio::spawn(use_rc_in_function(NonClone));

    // captured value is not `Send` because `&` references cannot be sent unless their referent is `Sync`
    // let err = tokio::spawn(get_rc_ref_but_sync_error(&Rc::new(NonClone)));

    // future cannot be sent between threads safely
    // within `[async block@src/lib.rs:27:28: 27:74]`, the trait `Send` is not implemented for `Rc<NonClone>`
    // let err = tokio::spawn(async move { get_rc_but_send_error(Rc::new(NonClone)).await });
}
