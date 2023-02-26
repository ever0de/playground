#[tokio::test]
async fn get_public_ip() {
    let ip = public_ip::addr().await.unwrap();
    println!("public ip: {ip}");
}
