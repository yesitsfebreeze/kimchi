import os
import tempfile
import subprocess
import shutil
import json
import commentjson

ROOT = os.getcwd()
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))


def get_git_repo_url():
	try:
		url = subprocess.check_output(
			["git", "config", "--get", "remote.origin.url"],
			cwd=ROOT,
			stderr=subprocess.DEVNULL
		).decode("utf-8").strip()
		# Convert SSH to HTTPS if needed
		if url.startswith("git@"):
			parts = url.split(":")
			if len(parts) == 2:
				url = f"https://github.com/{parts[1].replace('.git', '')}"
		elif url.endswith(".git"):
			url = url[:-4]
		return url
	except Exception:
		print("⚠️  Could not determine git remote URL, falling back to hardcoded path")
		return "https://github.com/yesitsfebreeze/kitsune"
	
REPO_URL = get_git_repo_url()

def get_last_commit_for_path(path: str) -> str:
	try:
		sha = subprocess.check_output(
			["git", "log", "-n", "1", "--format=%H", "--", path],
			cwd=ROOT
		).decode("utf-8").strip()
		return sha
	except Exception:
		print(f"⚠️  Could not get commit for {path}, using 'unknown'")
		return "unknown"

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

		files = [
			f"{name}.wasm"
		]
		if os.path.isdir(query_src):
			print(f"📦 Found queries for {name}, copying...")
			query_dst = os.path.join(target_dir, "queries")
			shutil.copytree(query_src, query_dst, dirs_exist_ok=True)
			for f in sorted(os.listdir(query_dst)):
				files.append(os.path.join("queries", f))
		else:
			print(f"⚠️  No queries found for {name}, skipping.")


		lang_path = os.path.join("grammars", name)
		version = get_last_commit_for_path(lang_path)

		AVAILABLE_LANGUAGES.append({
			"name": name,
			"version": version,
			"files": files
		})
		
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
	data = {
		"url": f"{REPO_URL}/raw/refs/heads/master/grammars/",
		"languages": AVAILABLE_LANGUAGES,
	}
	json.dump(data, f, indent=2)
	print(f"\n📝 Available languages written to {available_path}")
