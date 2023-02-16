// https://discord.com/channels/487203989830631435/803967640565317635/1075793448554221679
#[tokio::test]
async fn mutex() {
    use std::sync::Arc;
    use tokio::sync::Mutex;

    let data1 = Arc::new(Mutex::new(0));
    let data2 = Arc::clone(&data1);

    let t = tokio::spawn(async move {
        println!("start lock 1");
        let mut lock = data2.lock().await;
        println!("for loop ");
        for _ in 0..1000 {
            *lock += 1;
        }
        println!("finish ");
    });

    println!("lock 222");
    let mut lock = data1.lock().await;
    *lock += 1;
    println!("{}", lock);
    // If you wait for `t: JoinHandle<()>` to end, you will be in a deadlock.
    // need: drop lock
    drop(lock);

    t.await.unwrap();
}
