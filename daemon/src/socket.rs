use tokio::net::{UnixListener, UnixStream};
use tokio::io::{AsyncBufReadExt, AsyncWriteExt, BufReader};
use crate::handler::handle_request;
use std::{fs};
use anyhow::Result;



pub async fn start() -> Result<()> {
	let path = "/tmp/kitd.sock";

	// Clean up old socket
	let _ = fs::remove_file(path);

	let listener = UnixListener::bind(path)?;
	println!("[kitd] Listening on {}, Ctrl+C to exit", path);

	loop {
		let (stream, _) = listener.accept().await?;
		tokio::spawn(handle_connection(stream));
	}
}

pub async fn handle_connection(stream: UnixStream) {
	let (reader, mut writer) = stream.into_split();
	let mut buf_reader = BufReader::new(reader);
	let mut buffer = String::new();

	loop {
		match buf_reader.read_line(&mut buffer).await {
			Ok(0) => break, // EOF
			Ok(_n) => {
				println!("[kitd] <- {}", buffer.trim());

				match handle_request(&buffer).await {
					Ok(response) => {
						println!("[kitd] -> {}", truncate(&response, 200));
						let _ = writer.write_all(response.as_bytes()).await;
						let _ = writer.write_all(b"\n").await;
					}
					Err(err) => {
						let msg = format!("{{\"error\":\"{}\"}}\n", err);
						println!("[kitd] !! {}", err);
						let _ = writer.write_all(msg.as_bytes()).await;
					}
				}

				buffer.clear();
			}
			Err(err) => {
				println!("[kitd] !! socket read error: {}", err);
				break;
			}
		}
	}
}

fn truncate(s: &str, max: usize) -> &str {
	if s.len() > max { &s[..max] } else { s }
}
