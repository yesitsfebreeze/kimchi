use anyhow::{Result, Context};
use tree_sitter::{wasmtime::Engine, Parser, WasmStore, Node};
use serde_json::{json, Map, Value};
use std::fs;

use crate::jobs::language;

pub async fn handle(lang: String, code: Option<String>, path: Option<String>) -> Result<String> {
	let source = match code {
		Some(code) => code.clone(),
		None => {
			if let Some(path) = path {
				std::fs::read_to_string(&path).map_err(|e| anyhow::anyhow!("Failed to read file '{}': {}", path, e))?
			} else {
				return Err(anyhow::anyhow!("No source code provided"));
			}
		}
	};

	let result = analyze(&lang, &source).await?;
	Ok(serde_json::to_string(&result)?)
}

async fn analyze(lang: &str, code: &str) -> Result<serde_json::Value> {

	let engine = Engine::default();
	let mut store = WasmStore::new(&engine).unwrap();

	language::install(lang.to_string()).await?;
	let wasm_path = language::get_wasm_path(lang)
		.with_context(|| format!("Failed to get wasm path for language: {}", lang))?;

	println!("Loading WASM for language: {}", wasm_path.display());
	let wasm_bytes = fs::read(&wasm_path)
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
