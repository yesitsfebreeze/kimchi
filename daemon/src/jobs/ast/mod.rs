use anyhow::{Result, Context};
use tree_sitter::{wasmtime::Engine, Parser, WasmStore, Node};
use serde_json::{json, Map, Value};
use std::{env, fs, path::{PathBuf}};
use reqwest::Client;

const PROVIDER: &str = "https://github.com/tree-sitter";
const SUB_URL: &str = "releases/latest/download";

pub fn grammar_cache_dir() -> Result<PathBuf> {
	let exe_dir = env::current_exe()
		.context("Failed to get current executable path")?
		.parent()
		.context("Failed to get parent directory")?
		.to_path_buf();

	let path = exe_dir.join("grammars");
	std::fs::create_dir_all(&path).context("Failed to create grammar cache directory")?;

	Ok(path)
}

pub async fn ensure_wasm(lang: &str) -> Result<PathBuf> {
	let name = format!("tree-sitter-{}", lang);
	let url = format!("{}/{}/{}/{}.wasm", PROVIDER, name, SUB_URL, name);

	let cache_dir = grammar_cache_dir().expect("Failed to get grammar cache directory");
	let wasm_path = cache_dir.join(format!("{}.wasm", lang));

	if !wasm_path.exists() {
		println!("[kitsuned] downloading wasm for '{}'", lang);
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


pub async fn analyze_ast(lang: &str, code: &str) -> Result<serde_json::Value> {

	let engine = Engine::default();
	let mut store = WasmStore::new(&engine).unwrap();

	let wasm_path = &ensure_wasm(lang).await?;

	println!("Loading WASM for language: {}", wasm_path.display());
	let wasm_bytes = fs::read(wasm_path)
		.with_context(|| format!("Failed to read wasm grammar at {}", &wasm_path.display()))?;

	let wasm_lang = store
			.load_language(lang, &wasm_bytes)
			.unwrap();

	let mut parser = Parser::new();
	parser.set_wasm_store(store).unwrap();
	parser.set_language(&wasm_lang)?;

	let tree = parser.parse(code, None)
		 .ok_or_else(|| anyhow::anyhow!("Tree-sitter failed to parse"))?;

	let root = tree.root_node();
	Ok(serialize_node(root))
}

fn serialize_node(node: Node) -> Value {
	let mut obj = Map::new();

	obj.insert("type".into(), json!(node.kind()));
	obj.insert("range".into(), json!([node.start_byte(), node.end_byte()]));

	let children: Vec<_> = (0..node.child_count())
		.filter_map(|i| node.child(i))
		.map(|c| serialize_node(c))
		.collect();

	if !children.is_empty() {
		obj.insert("children".into(), json!(children));
	}

	json!(obj)
}
