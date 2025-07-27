const SOCKET_PATH: &str = "/tmp/kitsuned.sock";
const CLIENT_CHECK_INTERVAL: u64 = 5; // seconds

mod jobs;
use crate::jobs::language;
use crate::jobs::highlight;

use std::fs;
use std::collections::HashSet;
use std::sync::Arc;
use anyhow::Result;
use serde::Serialize;
use serde::Deserialize;
use tokio::net::UnixListener;
use tokio::net::UnixStream;
use tokio::io::AsyncBufReadExt;
use tokio::io::AsyncWriteExt;
use tokio::io::BufReader;
use tokio::sync::Mutex;
use tokio::task;
use tokio::time::interval;
use tokio::time::Duration;
use sysinfo::System;
use sysinfo::ProcessesToUpdate;

type SharedUsers = Arc<Mutex<HashSet<String>>>;


 // if true, will check for client processes and shutdown if none are found
const USE_SHUTDOWN_CHECK: bool = false;

#[tokio::main]
async fn main() -> Result<()> {
	socket().await?;
	Ok(())
}

#[derive(Debug, Serialize, Deserialize)]
pub enum Request {
	Connect {
		user: String,
	},
	Disconnect {
		user: String,
	},
	InstallLanguage {
		lang: String,
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
			Ok(format!("You disconnected as: {}", user))
		},
		Request::InstallLanguage { lang } => {
			return language::install(lang).await;
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

	if USE_SHUTDOWN_CHECK {
		check_for_processes();
	}

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
				match request(&buffer, &users).await {
					Ok(res) => {
						let _ = response(&mut writer, &res).await;
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

async fn response<W: AsyncWriteExt + Unpin>(
	writer: &mut W,
	response: &str,
) -> anyhow::Result<()> {
	let bytes = response.as_bytes();
	let size = bytes.len() as u32;

	writer.write_all(&size.to_be_bytes()).await?;
	writer.write_all(bytes).await?;

	Ok(())
}

pub fn check_for_processes() {
	task::spawn(async move {
		let check_time = Duration::from_secs(CLIENT_CHECK_INTERVAL);
		tokio::time::sleep(check_time).await;

		let mut interval = interval(check_time);
		loop {
			interval.tick().await;

			let mut system = System::new_all();
			system.refresh_processes(ProcessesToUpdate::All, true);

			let any_running = system.processes()
				.values()
				.any(|proc| {
					let name = proc.name().to_str()
						.map(|s| s.to_lowercase())
						.unwrap_or_default();
					name == "kitsune"
				});

			if !any_running {
				println!("no client detected. shutting down.");
				std::fs::remove_file(SOCKET_PATH).ok(); // cleanup
				std::process::exit(0);
			}
		}
	});
}
