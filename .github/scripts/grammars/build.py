import json # , os, subprocess, tempfile, shutil

# def run(cmd, cwd):
# 	print(f"→ {cmd}")
# 	subprocess.run(cmd, shell=True, check=True, cwd=cwd)

# def build_grammar(name, cfg):
# 	repo_url = cfg["repo"]
# 	grammar_dir = cfg.get("grammar_dir", ".")
# 	build_cmds = cfg.get("build_cmds", ["tree-sitter generate", "tree-sitter build-wasm"])

# 	with tempfile.TemporaryDirectory() as tmp:
# 		print(f"🧬 Cloning {name}...")
# 		clone_path = os.path.join(tmp, name)
# 		run(f"git clone {repo_url} {name}", tmp)

# 		grammar_path = os.path.join(clone_path, grammar_dir)
# 		for cmd in build_cmds:
# 			run(cmd, grammar_path)

# 		wasm_file = next(f for f in os.listdir(grammar_path) if f.endswith(".wasm"))
# 		out_path = os.path.join("out", f"{name}.wasm")
# 		os.makedirs("out", exist_ok=True)
# 		shutil.copy(os.path.join(grammar_path, wasm_file), out_path)
# 		print(f"✅ {name} → {out_path}")

with open("grammars.json") as f:
	config = json.load(f)

for name, cfg in config.items():
	print(f"🔍 Building grammar: {name}")
	# try:
	# 	build_grammar(name, cfg)
	# except Exception as e:
	# 	print(f"❌ Failed to build {name}: {e}")
