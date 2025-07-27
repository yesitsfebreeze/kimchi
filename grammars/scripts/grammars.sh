#!/usr/bin/env bash
echo "Building grammars..."

mkdir -p ./tmp

cd ./tmp
git clone https://github.com/tree-sitter/tree-sitter-agda agda
cd ./agda
tree-sitter generate grammar.js
tree-sitter build --wasm
cp ./tree-sitter-agda.wasm ../../grammars/agda.wasm
cd ../..


rm -rf ./tmp

# set -euo pipefail

# # Constants
# NVIM_TS_REPO="https://github.com/nvim-treesitter/nvim-treesitter"
# TS_QUERY_DIR="nvim-treesitter/queries"
# OUT_DIR="./grammars"
# WASM_URL_TEMPLATE="https://github.com/tree-sitter/tree-sitter-%s/releases/latest/download/tree-sitter-%s.wasm"

# # Clean and prepare
# rm -rf nvim-treesitter
# git clone --depth 1 "$NVIM_TS_REPO" nvim-treesitter

# mkdir -p "$OUT_DIR"

# LANGUAGES=()

# REPOS=(
# 	"agda|https://github.com/tree-sitter/tree-sitter-agda"
# 	"bash|https://github.com/tree-sitter/tree-sitter-bash"
# 	"c|https://github.com/tree-sitter/tree-sitter-c"
# 	"cpp|https://github.com/tree-sitter/tree-sitter-cpp"
# 	"css|https://github.com/tree-sitter/tree-sitter-css"
# 	"go|https://github.com/tree-sitter/tree-sitter-go"
# 	"haskell|https://github.com/tree-sitter/tree-sitter-haskell"
# 	"html|https://github.com/tree-sitter/tree-sitter-html"
# 	"java|https://github.com/tree-sitter/tree-sitter-java"
# 	"javascript|https://github.com/tree-sitter/tree-sitter-javascript"
# 	"jsdoc|https://github.com/tree-sitter/tree-sitter-jsdoc"
# 	"julia|https://github.com/tree-sitter/tree-sitter-julia"
# 	"ocaml|https://github.com/tree-sitter/tree-sitter-ocaml"
# 	"php|https://github.com/tree-sitter/tree-sitter-php"
# 	"python|https://github.com/tree-sitter/tree-sitter-python"
# 	"ql|https://github.com/tree-sitter/tree-sitter-ql"
# 	"regex|https://github.com/tree-sitter/tree-sitter-regex"
# 	"ruby|https://github.com/tree-sitter/tree-sitter-ruby"
# 	"rust|https://github.com/tree-sitter/tree-sitter-rust"
# 	"scala|https://github.com/tree-sitter/tree-sitter-scala"
# 	"typescript|https://github.com/tree-sitter/tree-sitter-typescript"
# 	"verilog|https://github.com/tree-sitter/tree-sitter-verilog"
# 	"yaml|https://github.com/ikatyang/tree-sitter-yaml"
# )


# mkrdir -p /tmp/grammars
# cd /tmp/grammars
# git clone https://github.com/tree-sitter/tree-sitter-agda agda
# cd agda
# tree-sitter generate ./grammar.js
# tree-sitter build --wasm

# for dir in nvim-treesitter/queries/*; do
# 	if [ -d "$dir" ]; then
# 		lang=$(basename "$dir")
# 		LANGUAGES+=("$lang")
# 		rm -rf "$OUT_DIR/$lang"
# 		mv "$dir" "$OUT_DIR/$lang"
# 	fi
# done

# rm -rf nvim-treesitter

# # Download WASM grammar for each language
# for lang in "${LANGUAGES[@]}"; do
# 	url=$(printf "$WASM_URL_TEMPLATE" "$lang" "$lang")
# 	dest="$OUT_DIR/$lang/tree-sitter-$lang.wasm"

# 	if curl -sSfL "$url" -o "$dest"; then
# 		echo "✅ Downloaded to $dest"
# 	else
# 		# echo "❌ No prebuilt WASM found for '$lang' at:"
# 		# echo "	 $url"
# 		rm -rf "$OUT_DIR/$lang"	# cleanup broken file just in case
# 	fi
# done

# echo "✅ Downloaded ${#LANGUAGES[@]} languages"

# # # List of languages to build packages for
# # LANGUAGES=("go" "lua" "javascript" "typescript" "python" "rust" "json")

# # # For each language, copy query files from nvim-treesitter into local grammar folder
# # for lang in "${LANGUAGES[@]}"; do
# #	 echo "→ Packaging $lang"

# #	 mkdir -p "$OUT_DIR/$lang"

# #	 QUERY_SRC="$TS_QUERY_DIR/$lang"
# #	 if [[ -d "$QUERY_SRC" ]]; then
# #		 cp "$QUERY_SRC"/*.scm "$OUT_DIR/$lang/" || echo "⚠️	No .scm files for $lang"
# #	 else
# #		 echo "⚠️	Skipping $lang (no query dir found)"
# #	 fi
# # done

# # echo "✅ Neovim queries copied to $OUT_DIR"
