// Browser port of Muster's replay engine (engine/replay.go), lab search
// (lab/search.go) and robustness check (lab/robust.go).
//
// The Go engine is the source of truth. web/parity.mjs replays every choice
// subset of every lab against Go-generated golden output and fails the build
// on any difference, so this file must track the Go semantics exactly.
(function (root) {
  "use strict";

  const PHASES = ["observe", "block", "review", "critical", "escalate", "respond"];

  function sat(cond, state) {
    if (cond.Op === "present") return state.has(cond.Fact);
    if (cond.Op === "absent") return !state.has(cond.Fact);
    return false;
  }

  function applyEffect(effect, state) {
    if (effect.Op === "add") state.add(effect.Fact);
    else if (effect.Op === "remove") state.delete(effect.Fact);
  }

  function sorted(state) {
    // Go sorts facts with sort.Strings: byte order, not locale order.
    return Array.from(state).sort((a, b) => (a < b ? -1 : a > b ? 1 : 0));
  }

  function reason(cond, suppressed) {
    if (cond.Op === "present") return suppressed ? cond.Fact : "missing:" + cond.Fact;
    if (cond.Op === "absent") return suppressed ? "missing:" + cond.Fact : "present:" + cond.Fact;
    return (suppressed ? "suppressed:" : "unsatisfied:") + cond.Fact;
  }

  function evaluate(control, event, state) {
    const result = {
      id: control.ID,
      role: control.Role,
      action: control.Action,
      matched: (control.EventIDs || []).includes(event.ID),
      disposition: "unmatched",
      reason: "",
      acted: false,
    };
    if (!result.matched) return result;

    result.disposition = "waiting";
    for (const cond of control.Requires || []) {
      if (!sat(cond, state)) {
        result.reason = reason(cond, false);
        return result;
      }
    }
    for (const cond of control.SuppressedBy || []) {
      if (sat(cond, state)) {
        result.disposition = "suppressed";
        result.reason = reason(cond, true);
        return result;
      }
    }
    result.disposition = "ready";
    return result;
  }

  // Evaluate every matched control whose phase is `action`; ready controls act
  // immediately, so later controls in the same phase see earlier effects.
  function evaluatePhase(event, state, controls, action) {
    const out = [];
    for (const control of controls) {
      const result = evaluate(control, event, state);
      if (!result.matched || control.Action !== action) continue;
      if (result.disposition !== "ready") {
        out.push(result);
        continue;
      }
      result.acted = true;
      for (const effect of control.Effects || []) applyEffect(effect, state);
      out.push(result);
    }
    return out;
  }

  function replay(scenario, controls) {
    const state = new Set(scenario.InitialFacts || []);
    const trace = [];

    for (const event of scenario.Events || []) {
      const entry = {
        id: event.ID,
        before: sorted(state),
        status: "",
        unsatisfied: [],
        controls: new Map(),
        observed: [],
        blocked: [],
        responded: [],
      };

      entry.unsatisfied = (event.Preconditions || []).filter((c) => !sat(c, state));
      if (entry.unsatisfied.length > 0) {
        entry.status = "skipped";
        entry.controls = [];
        entry.after = sorted(state);
        trace.push(entry);
        continue;
      }

      for (const control of controls) entry.controls.set(control.ID, evaluate(control, event, state));

      let blocked = false;
      for (const phase of PHASES) {
        if (phase === "review") {
          for (const effect of event.Effects || []) applyEffect(effect, state);
          entry.status = "applied";
        }
        for (const result of evaluatePhase(event, state, controls, phase)) {
          entry.controls.set(result.id, result);
          if (!result.acted) continue;
          if (phase === "observe") entry.observed.push(result.id);
          if (phase === "block") entry.blocked.push(result.id);
          if (phase === "respond") entry.responded.push(result.id);
        }
        if (phase === "block" && entry.blocked.length > 0) {
          blocked = true;
          break;
        }
      }

      if (blocked) entry.status = "blocked";
      entry.after = sorted(state);
      entry.controls = Array.from(entry.controls.values());
      trace.push(entry);
    }

    return { trace, terminal: sorted(state) };
  }

  // ---- Labs -------------------------------------------------------------

  class Lab {
    constructor(bundleLab) {
      this.manifest = bundleLab.manifest;
      this.traces = bundleLab.compiled;
      this.keys = this.manifest.choices.map((c) => c.key);
    }

    trace(id) {
      return this.traces.find((t) => t.id === (id || "baseline")) || this.traces[0];
    }

    expand(trace, keys) {
      const wanted = new Set();
      for (const key of keys) {
        const choice = this.manifest.choices.find((c) => c.key === key);
        if (!choice) throw new Error("unknown choice " + key);
        for (const id of choice.controls) wanted.add(id);
      }
      return trace.controls.filter((c) => wanted.has(c.ID));
    }

    run(traceID, keys) {
      const trace = this.trace(traceID);
      const result = replay(trace.scenario, this.expand(trace, keys));
      result.harm = this.harm(result.terminal);
      return result;
    }

    harm(terminal) {
      const set = new Set(terminal);
      return this.manifest.adverse.map((o) => o.fact).filter((f) => set.has(f));
    }

    harmBits(terminal) {
      const set = new Set(terminal);
      let bits = 0;
      this.manifest.adverse.forEach((o, i) => {
        if (set.has(o.fact)) bits |= 1 << i;
      });
      return bits;
    }

    factsOf(bits) {
      return this.manifest.adverse.map((o) => o.fact).filter((_, i) => bits & (1 << i));
    }

    maskKeys(mask) {
      return this.keys.filter((_, i) => mask & (1 << i));
    }

    table(traceID) {
      const n = this.keys.length;
      const trace = this.trace(traceID);
      const harm = new Array(1 << n);
      const runs = new Array(1 << n);
      for (let mask = 0; mask < 1 << n; mask++) {
        const r = replay(trace.scenario, this.expand(trace, this.maskKeys(mask)));
        runs[mask] = r;
        harm[mask] = this.harmBits(r.terminal);
      }
      return { harm, runs };
    }

    search(traceID) {
      const n = this.keys.length;
      const t = this.table(traceID);
      const found = [];
      for (let mask = 0; mask < 1 << n; mask++) {
        for (let c = 0; c < n; c++) {
          const bit = 1 << c;
          if (mask & bit) continue;
          const before = t.harm[mask];
          const after = t.harm[mask | bit];
          const gained = after & ~before;
          if (!gained) continue;
          const kind = before & ~after ? "tradeoff" : "worse";
          found.push({ kind, ctx: mask, added: c, gained });
        }
      }
      // Stable sort by context size, as Go's sort.SliceStable does.
      found.sort((a, b) => popcount(a.ctx) - popcount(b.ctx));
      const kept = [];
      for (const f of found) {
        const redundant = kept.some(
          (k) => k.kind === f.kind && k.added === f.added && k.gained === f.gained && (k.ctx & f.ctx) === k.ctx
        );
        if (!redundant) kept.push(f);
      }
      const findings = kept.map((k) => {
        const withAdded = k.ctx | (1 << k.added);
        return {
          kind: k.kind,
          context: this.maskKeys(k.ctx),
          added: this.keys[k.added],
          before: this.factsOf(t.harm[k.ctx]),
          after: this.factsOf(t.harm[withAdded]),
          new_harm: this.factsOf(k.gained),
          divergence: firstDivergence(t.runs[k.ctx], t.runs[withAdded]),
        };
      });
      findings.sort((a, b) => {
        if (a.kind !== b.kind) return a.kind === "worse" ? -1 : 1;
        return a.context.length - b.context.length;
      });
      return { choices: n, watchlines: 1 << n, findings };
    }

    holds(traceID, f) {
      const before = this.run(traceID, f.context);
      const after = this.run(traceID, f.context.concat([f.added]));
      const gained = this.harmBits(after.terminal) & ~this.harmBits(before.terminal);
      let want = 0;
      this.manifest.adverse.forEach((o, i) => {
        if (f.new_harm.includes(o.fact)) want |= 1 << i;
      });
      return (gained & want) === want;
    }

    robust() {
      const presets = [{ label: "No added defense", choices: [] }].concat(this.manifest.presets || []);
      const rows = presets.map((p) => {
        let base = 0;
        const cells = this.traces.map((trace, i) => {
          const r = this.run(trace.id, p.choices);
          const bits = this.harmBits(r.terminal);
          if (i === 0) base = bits;
          return { trace_id: trace.id, harm: this.factsOf(bits), change: change(base, bits) };
        });
        return { label: p.label, choices: p.choices, cells };
      });

      const baseline = this.search("baseline");
      const variants = this.traces.slice(1);
      const findings = baseline.findings.map((f) => {
        const holds = {};
        for (const v of variants) holds[v.id] = this.holds(v.id, f);
        return { finding: f, holds };
      });

      const emergent = {};
      for (const v of variants) {
        for (const f of this.search(v.id).findings) {
          if (!baseline.findings.some((g) => sameFinding(g, f))) (emergent[v.id] = emergent[v.id] || []).push(f);
        }
      }
      return { rows, findings, emergent };
    }
  }

  function popcount(x) {
    let n = 0;
    while (x) {
      n += x & 1;
      x >>>= 1;
    }
    return n;
  }

  function change(base, harm) {
    if (harm === base) return "same";
    if (!(harm & ~base)) return "better";
    if (!(base & ~harm)) return "worse";
    return "different";
  }

  function sameList(a, b) {
    return a.length === b.length && a.every((x, i) => x === b[i]);
  }

  function sameFinding(a, b) {
    return a.kind === b.kind && a.added === b.added && sameList(a.context, b.context) && sameList(a.new_harm, b.new_harm);
  }

  function acted(entry) {
    return entry.controls.filter((c) => c.acted).map((c) => c.id).sort();
  }

  function firstDivergence(a, b) {
    for (let i = 0; i < a.trace.length; i++) {
      if (i >= b.trace.length) return a.trace[i].id;
      const ea = a.trace[i];
      const eb = b.trace[i];
      if (ea.status !== eb.status || !sameList(ea.after, eb.after) || !sameList(acted(ea), acted(eb))) return ea.id;
    }
    return "";
  }

  const api = { replay, Lab, firstDivergence };
  root.MusterEngine = api;
  if (typeof module !== "undefined" && module.exports) module.exports = api;
})(typeof globalThis !== "undefined" ? globalThis : this);
