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
	version: String,
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

async fn get_manifest() -> Result<Manifest> {
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

	Ok(manifest)
}

pub async fn install_multi(langs: Vec<String>) -> Result<String> {
	let mut results = Vec::new();
	let manifest = get_manifest().await?;
	for lang in langs {
		match _install(&lang, &manifest).await {
			Ok(msg) => results.push(msg),
			Err(e) => results.push(format!("Failed to install {}: {}", lang, e)),
		}
	}
	Ok(results.join("\n"))
}

pub async fn install(lang: String) -> Result<String> {
	if lang.is_empty() {
		return Err(anyhow::anyhow!("Language name cannot be empty"));
	}

	let manifest = get_manifest().await?;
	_install(&lang, &manifest).await
		.with_context(|| format!("Failed to install language '{}'", lang))
}

async fn _install(lang: &str, manifest: &Manifest) -> Result<String> {
	if lang.is_empty() {
		return Err(anyhow::anyhow!("Language name cannot be empty"));
	}

	let client = Client::new();

	let language = manifest.languages.iter()
		.find(|l| l.name == lang)
		.with_context(|| format!("Language '{}' is not available", lang))?;

	let cache_dir = get_cache_dir(&lang)?;

	let version_path = cache_dir.join(".version");

	if version_path.exists() {
		let current_version = fs::read_to_string(&version_path)
			.unwrap_or_default()
			.trim()
			.to_string();

		if current_version == language.version {
			let msg = format!("Language '{}' is already up to date (version {}).", lang, current_version);
			println!("{}", msg);
			return Ok(msg);
		}
	}

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

	fs::write(&version_path, language.version.trim())
		.with_context(|| format!("Failed to write version to {}", version_path.display()))?;

	let msg = format!("Language '{}' installed successfully.", lang);
	println!("{}", msg);
	Ok(msg)
}

pub fn get_wasm_path(lang: &str) -> Result<PathBuf> {
	let cache_dir = get_cache_dir(lang)?;
	let wasm_path = cache_dir.join(format!("{}.wasm", lang));

	if !wasm_path.exists() {
		return Err(anyhow::anyhow!("WASM file for language '{}' not found", lang));
	}

	Ok(wasm_path)
}
