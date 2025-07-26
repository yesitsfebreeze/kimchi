use serde::{Serialize, Deserialize};
use anyhow::Result;

use crate::jobs::ast;
// use crate::jobs::lsp;


// TODO: maybe need rename/refactor/goto/ etc requests for lsp
#[derive(Debug, Serialize, Deserialize)]
struct Request {
	lang: String,
	job: String,
	code: Option<String>,
	path: Option<String>,
	region_start: Option<usize>,
	region_end: Option<usize>,
}

async fn handle_ast(req: &Request) -> Result<String> {
	let source = match &req.code {
		Some(code) => code.clone(),
		None => {
			if let Some(path) = &req.path {
				std::fs::read_to_string(path).map_err(|e| anyhow::anyhow!("Failed to read file '{}': {}", path, e))?
			} else {
				return Err(anyhow::anyhow!("No source code provided"));
			}
		}
	};

	let ast_result = ast::analyze_ast(&req.lang, &source).await?;
	Ok(serde_json::to_string(&ast_result)?)
}

// async fn handle_lsp(req: &Request) -> Result<String> {
// 	// let lsp_result = lsp::analyze_code(&req.lang, &req.source).await?;
// 	// Ok(serde_json::to_string(req)?)
// 	Ok(serde_json::to_string(&"LSP analysis not implemented yet")?)
// }


pub async fn handle_request(line: &str) -> Result<String> {
	let req: Request = serde_json::from_str(line)?;
	match req.job.as_str() {
		"ast" => {
			return handle_ast(&req).await;
		}
		// "lsp" => {
		// 	return handle_lsp(&req).await;
		// },
		_ => anyhow::bail!("Unknown job '{}'", req.job),
	}
}
