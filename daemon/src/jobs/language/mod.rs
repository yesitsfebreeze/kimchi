use anyhow::Result;
use anyhow::Context;
use std::env;
use std::fs;
use std::path::PathBuf;
use reqwest::Client;
use serde::Deserialize;

const MANIFEST_URL: &str = "https://raw.githubusercontent.com/yesitsfebreeze/kitsune/refs/heads/master/grammars/available.json";

#[derive(Debug, Deserialize)]
struct Language {
	name: String,
	files: Vec<String>,
}

#[derive(Debug, Deserialize)]
struct Manifest {
	url: String,
	languages: Vec<Language>,
}

fn get_cache_dir(lang: &str) -> Result<PathBuf> {
	let exe_dir = env::current_exe()
		.context("Failed to get current executable path")?
		.parent()
		.context("Failed to get parent directory")?
		.to_path_buf();

	let path = exe_dir.join("grammars").join(lang);
	std::fs::create_dir_all(&path).context("Failed to create grammar cache directory")?;

	Ok(path)
}

pub async fn install(lang: String) -> Result<String> {
	let client = Client::new();

	let resp = client.get(MANIFEST_URL)
		.send()
		.await
		.context("Failed to download manifest")?
		.text()
		.await
		.context("Failed to read manifest content")?;

	 let manifest: Manifest = serde_json::from_str(&resp)
		.context("Failed to parse manifest JSON")?;

	let language = manifest.languages.iter()
		.find(|l| l.name == lang)
		.with_context(|| format!("Language '{}' is not available", lang))?;

	let cache_dir = get_cache_dir(&lang)?;

	for file in &language.files {
		let url = format!("{}/{}/{}", manifest.url.trim_end_matches('/'), lang, file);
		let target = cache_dir.join(file);
		if let Some(parent) = target.parent() {
				fs::create_dir_all(parent)?;
		}

		let content = client.get(&url)
			.send()
			.await
			.with_context(|| format!("Failed to download {}", file))?
			.bytes()
			.await
			.with_context(|| format!("Failed to read content for {}", file))?;

		fs::write(&target, &content)
			.with_context(|| format!("Failed to write to {}", target.display()))?;
	}

	let msg = format!("Language '{}' installed successfully.", lang);
	println!("{}", msg);
	Ok(msg)
}



// async fn ensure_wasm(lang: &str) -> Result<PathBuf> {
// 	let name = format!("tree-sitter-{}", lang);
// 	let url = format!("{}/{}/{}/{}.wasm", WASM_PROVIDER, name, WASM_SUB_URL, name);


// 	let wasm_path = cache_dir.join(format!("{}.wasm", lang));

// 	if !wasm_path.exists() {
// 		println!("downloading wasm for '{}'", lang);
// 		let bytes = Client::new()
// 			.get(url)
// 			.send()
// 			.await
// 			.context("failed to download grammar wasm")?
// 			.bytes()
// 			.await
// 			.context("failed to read wasm bytes")?;
// 		fs::write(&wasm_path, &bytes)?;
// 	}

// 	Ok(wasm_path)
// }

// async fn ensure_queries(lang: &str) -> Result<&str> {
// 	let name = format!("tree-sitter-{}", lang);
// 	let base_url = format!("{}/{}/{}", QUERY_PROVIDER, name, QUERY_SUB_URL);

// 	let client = reqwest::Client::new();

// 	let cache_dir = get_cache_dir(lang)?.join("queries");
// 	fs::create_dir_all(&cache_dir)?;

// 	for file in KNOWN_QUERY_FILES {
// 		let url = format!("{}{}", base_url, file);
// 		let response = client.get(&url).send().await?;
// 		if response.status().is_success() {
// 			println!("downloading query file '{}'", url);
// 			let content = response.bytes().await?;
// 			fs::write(cache_dir.join(file), content)?;
// 		} else if response.status().as_u16() != 404 {
// 			return Err(anyhow::anyhow!("Unexpected error fetching {}: {}", file, response.status()));
// 		}
// 	}

// 	Ok("OK")
// }


pub fn get_wasm_path(lang: &str) -> Result<PathBuf> {
	let cache_dir = get_cache_dir(lang)?;
	let wasm_path = cache_dir.join(format!("{}.wasm", lang));

	if !wasm_path.exists() {
		return Err(anyhow::anyhow!("WASM file for language '{}' not found", lang));
	}

	Ok(wasm_path)
}
