pub mod public_ip;
pub mod yaml;

#[test]
fn if_address() {
    let ifaces = get_if_addrs::get_if_addrs().unwrap();
    println!("{ifaces:#?}");
}
