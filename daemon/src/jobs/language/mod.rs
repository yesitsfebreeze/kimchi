use anyhow::Result;
use anyhow::Context;
use std::env;
use std::fs;
use std::path::PathBuf;
use reqwest::Client;

const WASM_PROVIDER: &str = "https://github.com/tree-sitter";
const WASM_SUB_URL: &str = "releases/latest/download";

const QUERY_PROVIDER: &str = "https://github.com/nvim-treesitter/nvim-treesitter";
const QUERY_SUB_URL: &str = "refs/heads/master/queries/";

const KNOWN_QUERY_FILES: &[&str] = &[
	"highlights.scm",
	"injections.scm",
	"locals.scm",
	"textobjects.scm",
	"folds.scm",
	"indents.scm",
	"tags.scm",
];

pub async fn install(lang: String) -> Result<String> {
	ensure_wasm(&lang).await?;
	ensure_queries(&lang).await?;

	Ok("language installed successfully".to_string())
}

fn grammar_cache_dir(lang: &str) -> Result<PathBuf> {
	let exe_dir = env::current_exe()
		.context("Failed to get current executable path")?
		.parent()
		.context("Failed to get parent directory")?
		.to_path_buf();

	let path = exe_dir.join("grammars").join(lang);
	std::fs::create_dir_all(&path).context("Failed to create grammar cache directory")?;

	Ok(path)
}

async fn ensure_wasm(lang: &str) -> Result<PathBuf> {
	let name = format!("tree-sitter-{}", lang);
	let url = format!("{}/{}/{}/{}.wasm", WASM_PROVIDER, name, WASM_SUB_URL, name);

	let cache_dir = grammar_cache_dir(lang).expect("Failed to get grammar cache directory");
	let wasm_path = cache_dir.join(format!("{}.wasm", lang));

	if !wasm_path.exists() {
		println!("downloading wasm for '{}'", lang);
		let bytes = Client::new()
			.get(url)
			.send()
			.await
			.context("failed to download grammar wasm")?
			.bytes()
			.await
			.context("failed to read wasm bytes")?;
		fs::write(&wasm_path, &bytes)?;
	}

	Ok(wasm_path)
}

async fn ensure_queries(lang: &str) -> Result<&str> {
	let name = format!("tree-sitter-{}", lang);
	let base_url = format!("{}/{}/{}", QUERY_PROVIDER, name, QUERY_SUB_URL);

	let client = reqwest::Client::new();

	let cache_dir = grammar_cache_dir(lang)?.join("queries");
	fs::create_dir_all(&cache_dir)?;

	for file in KNOWN_QUERY_FILES {
		let url = format!("{}{}", base_url, file);
		let response = client.get(&url).send().await?;
		if response.status().is_success() {
			println!("downloading query file '{}'", url);
			let content = response.bytes().await?;
			fs::write(cache_dir.join(file), content)?;
		} else if response.status().as_u16() != 404 {
			return Err(anyhow::anyhow!("Unexpected error fetching {}: {}", file, response.status()));
		}
	}

	Ok("OK")
}


pub fn get_wasm_path(lang: &str) -> Result<PathBuf> {
	let cache_dir = grammar_cache_dir(lang)?;
	let wasm_path = cache_dir.join(format!("{}.wasm", lang));

	if !wasm_path.exists() {
		return Err(anyhow::anyhow!("WASM file for language '{}' not found", lang));
	}

	Ok(wasm_path)
}
