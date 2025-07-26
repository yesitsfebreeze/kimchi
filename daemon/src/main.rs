mod socket;
mod handler;
mod jobs;
use anyhow::Result;

#[tokio::main]
async fn main() -> Result<()> {
	socket::start().await?;
	Ok(())
}
