import json , os ,tempfile, subprocess, shutil

ROOT = os.getcwd()
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))

print(os.getcwd())

def run(cmd, cwd):
	print(f"→ {cmd}")
	subprocess.run(cmd, shell=True, check=True, cwd=cwd)

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
		grammar_file = os.join(clone_path, wasm_file)
		target = os.path.join(target_dir, f"{name}.wasm")
		shutil.copy(grammar_file, os.path.join(target_dir, f"{name}.wasm"))
		print(f"✅ {name} → {target}")

with open(os.path.join(SCRIPT_DIR, "grammars.json")) as f:
	config = json.load(f)

for name, cfg in config.items():
	print(f"🔍 Building grammar: {name}")
	# try:
	# 	build_grammar(name, cfg)
	# except Exception as e:
	# 	print(f"❌ Failed to build {name}: {e}")
