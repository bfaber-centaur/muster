# October 2026 update: responding to the sprint reviews

Three reviewers read the sprint submission. This note maps what they asked
for to what changed in the repository, and records what the new work found,
including where it cuts against the original framing.

## What the reviewers asked for

| Reviewer point | What changed |
|---|---|
| R1, R3: align the README with the implemented scenario; the paper and repository describe access differently | Fixed. The READMEs said the HF connector path was independent of node root and showed the admission policy leaving cluster-admin reachable. The scenario never encoded that, and Hugging Face's timeline says the connector credential came from the same node-root secret read as the mesh key. The engine output and the paper's Table 1 were right; the READMEs and notebook help text now match them. The independent-path reading is kept as an explicit variant so the two can be compared. |
| R1: test a few alternative attack paths, to show how much conclusions depend on the reconstruction | Every lab now has alternative versions of events, and `muster robust` replays each preset against all of them. |
| R2: an interactive website | Muster Lab runs in the browser (`web/`), on an engine checked against Go on every replay. |
| R2: enrich the story with failed attempts, denied actions, blocked requests | "Failed attempt" variants: the record shows something the adversary tried that a defense stopped, and the variant removes that defense. HF's timeline supplied several (private-link allowlist, denied cloud-role mutations, denied credential-store operations). |
| R2: use failed attempts to create counterfactual branches and check whether the Watchline still holds | That is what `robust` reports, preset by preset. |
| R3: treat the isolation result as a way to examine response decisions, not evidence that early isolation makes incidents worse | The README now says so directly, next to the result. |

## Engine fix

The replay loop credited every matched control in the observe, block and
respond phases, including controls waiting on a requirement or suppressed. A
suppressed Picket therefore still blocked its event. No published scenario had
a Picket with `requires` or `suppressed_by`, so no terminal state in the
report changes. One summary does: Gate 2 early with the broad correlated
controls reported a first observation from a Vedette suppressed by compromised
telemetry. Three of the new scenario authors hit the bug independently while
modeling, which is how it surfaced.

## New tools

* `lab.yaml` manifests: adverse outcomes, human-facing choices, presets, labels,
  and alternative traces. They select; they do not change semantics.
* `muster search`: replays every subset of choices and lists minimal cases
  where adding one defense strictly grows the adverse set.
* `muster robust`: presets × alternative traces, plus whether each backfire
  survives.
* `muster survey`: both, across every lab.

## New scenarios

<!-- labs:start -->
<!-- labs:end -->

## What the survey says

<!-- survey:start -->
<!-- survey:end -->

1. **Backfires are rare in the public-record incidents.** None of the
   public-record labs contains one, on the recorded trace or any alternative.
   That supports the sprint report's own position: in most incidents,
   stateful replay agrees with dependency reasoning, and replay earns its keep
   only when defensive history matters.
2. **Where backfires do appear, the mechanism is the one JACKPOT was built
   around**: a response that destroys evidence a later response needs. Fire
   overhaul and sepsis care are not constructed to produce it. Their
   professional guidance warns about it in so many words. That answers part of
   the "synthetic existence proof" objection: the mechanism is documented
   practice elsewhere, even though its prevalence in cyber incidents is still
   unknown.
3. **In the real-world analogues it shows up first as a tradeoff, not a strict
   reversal.** Early antibiotics and early overhaul each buy something real
   (time; no rekindle). The strict reversal appears in specific conditions:
   a susceptible organism, a scene already being preserved. The toy's strict
   reversal is the clean case; reality mostly offers the tradeoff.
4. **The recorded trace flatters prevention.** Most presets that work on the
   recorded incident change outcome in some alternative. The Watchlines that
   held everywhere tended to act at the first foothold (HF page + isolate; PyPI
   evaluation egress), or combined prevention with a Reserve that does not
   depend on the adversary's route.

## Open questions

* Alternatives are hand-written. Generating them from failed-attempt records,
  or from perturbations of preconditions, would make robustness less
  dependent on the author's imagination.
* The survey counts backfires by minimal context. Several findings in a lab
  are often one mechanism seen through different contexts; a grouping by
  mechanism would read better.
* Prior work. The security version of the tension is already in incident
  response guidance. NIST SP 800-61 Rev. 2, §3.3.1 "Choosing a Containment
  Strategy", lists "Need for evidence preservation" among the criteria for
  choosing a containment strategy, describes redirecting an attacker to a
  sandbox "to gather additional evidence", and warns that "the delayed
  containment strategy is dangerous because an attacker could escalate
  unauthorized access" (https://nvlpubs.nist.gov/nistpubs/SpecialPublications/NIST.SP.800-61r2.pdf;
  Rev. 2 has since been superseded by Rev. 3, which was not checked here). So
  JACKPOT's mechanism is not new to practitioners. The guidance treats it as
  a judgment call. Muster's contribution, if any, is making it a checkable
  property of a specific trace and Watchline: which response loses which
  evidence, and what that leaves reachable. The report should cite 800-61 and
  frame the result that way.
