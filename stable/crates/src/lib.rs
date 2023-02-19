#[test]
fn if_address() {
    let ifaces = get_if_addrs::get_if_addrs().unwrap();
    println!("{ifaces:#?}");
}

#[tokio::test]
async fn get_public_ip() {
    let ip = public_ip::addr().await.unwrap();
    println!("public ip: {ip}");
}
