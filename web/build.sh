#!/usr/bin/env sh
# Build the browser lab into web/dist/index.html.
#
# The page is one self-contained file: the browser engine (web/engine.js) and
# every examples/**/lab.yaml, inlined. Before writing it, web/parity.mjs checks
# the browser engine against the Go engine on every choice subset of every
# trace of every lab, and the build fails on any difference.
set -eu
cd "$(dirname "$0")/.."

mkdir -p web/dist
go run . bundle --root examples --out web/dist/labs.json --golden web/dist/golden.json
node web/parity.mjs web/dist/labs.json web/dist/golden.json | tee web/dist/parity.txt

node - <<'EOF'
const fs = require("fs");
const page = fs.readFileSync("web/index.html", "utf8");
const engine = fs.readFileSync("web/engine.js", "utf8");
const labs = fs.readFileSync("web/dist/labs.json", "utf8").trim();
const replays = JSON.parse(fs.readFileSync("web/dist/golden.json", "utf8")).runs.length;
// Keep "</script" out of inlined JSON.
const safe = labs.replace(/<\//g, "<\\/");
const out = page
  .replace("/*__ENGINE__*/", () => engine)
  .replace("/*__LABS__*/null", () => safe)
  .replace("/*__PARITY__*/0", () => String(replays));
fs.writeFileSync("web/dist/index.html", out);
console.log(`wrote web/dist/index.html (${(out.length / 1024).toFixed(0)} KB)`);
EOF
