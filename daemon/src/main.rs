const SOCKET_PATH: &str = "/tmp/kitsuned.sock";

mod jobs;
use crate::jobs::highlight;

use std::{fs};
use anyhow::Result;
use tokio::net::{UnixListener, UnixStream};
use tokio::io::{AsyncBufReadExt, AsyncWriteExt, BufReader};
use serde::{Serialize, Deserialize};
use std::collections::HashSet;
use std::sync::Arc;
use tokio::sync::{Mutex, MutexGuard};

type SharedUsers = Arc<Mutex<HashSet<String>>>;

#[tokio::main]
async fn main() -> Result<()> {
	socket().await?;
	Ok(())
}

#[derive(Debug, Serialize, Deserialize)]
pub enum Request {
	Shutdown,
	Connect {
		user: String,
	},
	Disconnect {
		user: String,
	},
	Highlight {
		lang: String,
		code: Option<String>,
		path: Option<String>,
	},

	// TODO: maybe need rename/refactor/goto/ etc requests for lsp
	// Lsp {
	//	 lang: String,
	//	 path: Option<String>,
	//	 // ...
	// }
}
async fn request(line: &str, users: &SharedUsers) -> Result<String> {

	let req: Request = match serde_json::from_str(line) {
		Ok(r) => r,
		Err(e) => {
			return Err(anyhow::anyhow!("ERR: Could not parse request: {}", e));
		}
	};

	match req {
		Request::Shutdown => {
			println!("Shutting down...");
			std::process::exit(0);
		},
		Request::Connect { user } => {
			let mut set = users.lock().await;
			if !set.insert(user.clone()) {
				return Err(anyhow::anyhow!("User '{}' already connected", user));
			}
			Ok(format!("You connected as: {}", user))
		},

		Request::Disconnect { user } => {
			let mut set = users.lock().await;
			if !set.remove(&user) {
				return Err(anyhow::anyhow!("User '{}' was not connected", user));
			}
			shutdown(set);
			Ok(format!("You disconnected as: {}", user))
		},
		Request::Highlight { lang, code, path } => {
			return highlight::handle(lang, code, path).await;
		}
	}
}

async fn socket() -> Result<()> {
	let _ = fs::remove_file(SOCKET_PATH);
	let listener = UnixListener::bind(SOCKET_PATH)?;
	println!("Listening on {}, Ctrl+C to exit", SOCKET_PATH);

	let users: SharedUsers = Arc::new(Mutex::new(HashSet::new()));

	loop {
		let (stream, _) = listener.accept().await?;
		tokio::spawn(handle(stream, users.clone()));
	}
}

async fn handle(stream: UnixStream, users: SharedUsers) {
	let (reader, mut writer) = stream.into_split();
	let mut buf_reader = BufReader::new(reader);
	let mut buffer = String::new();

	loop {
		match buf_reader.read_line(&mut buffer).await {
			Ok(0) => break, // EOF
			Ok(_n) => {
				// println!("input: {}", buffer.trim());
				match request(&buffer, &users).await {
					Ok(response) => {
						let _ = writer.write_all(response.as_bytes()).await;
						let _ = writer.write_all(b"\n").await;
					}
					Err(err) => {
						let msg = format!("{{\"error\":\"{}\"}}\n", err);
						println!("{}", err);
						let _ = writer.write_all(msg.as_bytes()).await;
					}
				}

				buffer.clear();
			}
			Err(err) => {
				println!("ERR: socket read error: {}", err);
				break;
			}
		}
	}
}

fn shutdown(set: MutexGuard<'_, HashSet<String>>) {
	if set.is_empty() {
    println!("No users remaining. Shutting down.");
    std::fs::remove_file(SOCKET_PATH).ok(); // cleanup
    std::process::exit(0);
	}
}
