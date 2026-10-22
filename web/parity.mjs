// Check the browser engine against the Go engine.
//
// usage: node web/parity.mjs LABS.json GOLDEN.json
// Both files come from `muster bundle --out LABS.json --golden GOLDEN.json`.
import { readFileSync } from "node:fs";
import { createRequire } from "node:module";

const require = createRequire(import.meta.url);
const { Lab } = require("./engine.js");

const [labsPath, goldenPath] = process.argv.slice(2);
const bundle = JSON.parse(readFileSync(labsPath, "utf8"));
const golden = JSON.parse(readFileSync(goldenPath, "utf8"));

const labs = new Map(bundle.labs.map((b) => [b.manifest.id, new Lab(b)]));
// Go encodes map keys sorted and nil slices as null; compare canonically.
const canon = (v) =>
  Array.isArray(v)
    ? v.map(canon)
    : v && typeof v === "object"
      ? Object.fromEntries(Object.keys(v).sort().map((k) => [k, canon(v[k])]))
      : v;
const same = (a, b) => JSON.stringify(canon(a)) === JSON.stringify(canon(b));
let failures = 0;
const fail = (msg) => {
  if (failures++ < 20) console.error("MISMATCH " + msg);
};

for (const g of golden.runs) {
  const lab = labs.get(g.lab);
  const r = lab.run(g.trace, g.choices);
  const where = `${g.lab}/${g.trace}/[${g.choices}]`;
  if (!same(r.harm, g.harm)) fail(`${where} harm ${r.harm} != ${g.harm}`);
  if (!same(r.terminal, g.terminal)) fail(`${where} terminal`);
  g.events.forEach((ge, i) => {
    const e = r.trace[i];
    const actedInOrder = e.controls.filter((c) => c.acted).map((c) => c.id);
    const got = { id: e.id, status: e.status, acted: actedInOrder, observed: e.observed, blocked: e.blocked, responded: e.responded };
    if (!same(got, ge)) fail(`${where} event ${ge.id}: ${JSON.stringify(got)} != ${JSON.stringify(ge)}`);
  });
}

for (const g of golden.searches) {
  const got = labs.get(g.lab).search(g.trace);
  const want = g.result.findings.map((f) => ({ ...f, divergence: f.divergence || "" }));
  if (!same(got.findings, want)) fail(`search ${g.lab}/${g.trace}`);
}

for (const g of golden.robust) {
  const got = labs.get(g.lab).robust();
  const rows = g.result.rows.map((r) => ({ label: r.label, choices: r.choices, cells: r.cells }));
  if (!same(got.rows, rows)) fail(`robust rows ${g.lab}`);
  const holds = (g.result.findings || []).map((f) => f.holds);
  if (!same(got.findings.map((f) => f.holds), holds)) fail(`robust holds ${g.lab}`);
  const emergent = Object.fromEntries(
    Object.entries(g.result.emergent || {}).map(([k, v]) => [k, v.map((f) => ({ ...f, divergence: f.divergence || "" }))])
  );
  if (!same(got.emergent, emergent)) fail(`robust emergent ${g.lab}`);
}

const runs = golden.runs.length;
if (failures) {
  console.error(`parity: ${failures} mismatch(es) across ${runs} replays`);
  process.exit(1);
}
console.log(`parity: browser engine matches Go on ${runs} replays, ${golden.searches.length} searches, ${golden.robust.length} robustness checks`);
