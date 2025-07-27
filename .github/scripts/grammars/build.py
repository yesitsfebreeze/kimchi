import os
import tempfile
import subprocess
import shutil
import json
import commentjson

ROOT = os.getcwd()
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))

def run(cmd, cwd):
	print(f"→ {cmd}")
	subprocess.run(cmd, shell=True, check=True, cwd=cwd)

def clone_nvim_treesitter_queries():
	tmp = tempfile.mkdtemp()
	nvim_repo = os.path.join(tmp, "nvim-treesitter")
	print("🌳 Cloning nvim-treesitter...")
	run("git clone --depth 1 https://github.com/nvim-treesitter/nvim-treesitter", tmp)
	return os.path.join(nvim_repo, "queries")

NVIM_QUERY_DIR = clone_nvim_treesitter_queries()

AVAILABLE_LANGUAGES = []

def build_grammar(name, cfg):
	target_dir = os.path.join(ROOT, "grammars", name)
	repo_url = cfg["repo"]
	build = cfg.get("build", ["tree-sitter generate ./grammar.js", "tree-sitter build --wasm"])

	with tempfile.TemporaryDirectory() as tmp:
		print(f"🧬 Cloning {name}...")
		clone_path = os.path.join(tmp, name)
		run(f"git clone {repo_url} {name}", tmp)
		for cmd in build:
			run(cmd, clone_path)
		os.makedirs(target_dir, exist_ok=True)
		wasm_file = "tree-sitter-" + name + ".wasm"
		grammar_file = os.path.join(clone_path, wasm_file)
		target = os.path.join(target_dir, f"{name}.wasm")
		shutil.copy(grammar_file, os.path.join(target_dir, f"{name}.wasm"))
		
		query_src = os.path.join(NVIM_QUERY_DIR, name)
		if os.path.isdir(query_src):
			print(f"📦 Found queries for {name}, copying...")
			query_dst = os.path.join(target_dir, "queries")
			shutil.copytree(query_src, query_dst, dirs_exist_ok=True)
		else:
			print(f"⚠️  No queries found for {name}, skipping.")

		AVAILABLE_LANGUAGES.append(name)
		print(f"✅ {name} → {target}")

with open(os.path.join(SCRIPT_DIR, "grammars.jsonc")) as f:
	config = commentjson.load(f)

for name, cfg in config.items():
	print(f"🔍 Building grammar: {name}")
	try:
		build_grammar(name, cfg)
	except Exception as e:
		print(f"❌ Failed to build {name}: {e}")

available_path = os.path.join(ROOT, "grammars", "available.json")
with open(available_path, "w") as f:
	json.dump(AVAILABLE_LANGUAGES, f, indent=2)
	print(f"\n📝 Available languages written to {available_path}")
