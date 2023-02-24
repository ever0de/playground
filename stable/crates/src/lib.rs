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

#[test]
fn yaml() {
    use serde::{Deserialize, Deserializer, Serialize};

    fn status_to_string<'de, D>(deserializer: D) -> Result<String, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s = BookStatus::deserialize(deserializer)?;

        Ok(match s {
            BookStatus::Todo => "todo",
            BookStatus::Doing => "doing",
            BookStatus::Done => "done",
        }
        .to_owned())
    }

    #[derive(Debug, Serialize, Deserialize)]
    struct Book {
        title: String,
        #[serde(deserialize_with = "status_to_string")]
        status: String,
    }

    #[derive(Debug, Serialize, Deserialize)]
    enum BookStatus {
        #[serde(rename = "todo")]
        Todo,
        #[serde(rename = "doing")]
        Doing,
        #[serde(rename = "done")]
        Done,
    }

    let book = Book {
        title: "The Rust Programming Language".to_owned(),
        status: "todo".to_owned(),
    };

    let yaml = serde_yaml::to_string(&book).unwrap();
    println!("{yaml}");

    let book: Book = serde_yaml::from_str(&yaml).unwrap();
    println!("{book:#?}");
}
