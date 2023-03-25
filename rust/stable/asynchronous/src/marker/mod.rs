use std::{rc::Rc, sync::Arc};

struct MyType {
    _phantom: std::marker::PhantomData<Rc<()>>,
}

fn foo(_my_val: Arc<MyType>) {}

fn main() {
    let _my_val = Arc::new(MyType {
        _phantom: std::marker::PhantomData,
    });

    // NOTE: ERROR not implemented Send, Sync
    // tokio::spawn(async move {
    //     foo(my_val.clone());
    // });
}
