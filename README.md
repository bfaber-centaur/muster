# Muster

**A small incident-replay harness for asking where defensive interventions actually change an attack.**

> “If he sends reinforcements everywhere, he will everywhere be weak.”
>
> — Sun Tzu, *The Art of War*, VI.17

Muster takes a compressed attack trajectory, applies different defensive controls, and shows which downstream attacker capabilities remain reachable.

It was built alongside **Watchline**, a role-based model of layered defense.

```text
                  ATTACK TRAJECTORY
                         │
                         ▼
             ─────── transition ───────
                         │
             🔭 VEDETTE  │ observe
             🛡️ PICKET   │ intervene
             🔧 RESERVE  │ recover
                         │
                         ▼
             ─────── transition ───────
                         │
                        ...
```

> **You cannot defend every transition. Where do you place your watchline?**

**Try it in your browser:** [Muster Lab](https://bfaber-centaur.github.io/muster/)
replays nine scenarios (cyber intrusions, AI-agent incidents, an offshore
blowout, a bank heist, sepsis care, and fire investigation) on the same engine,
with a search for defenses that backfire and a check of how each defense holds
up when events go differently. To build it locally, run `./web/build.sh` and
open `web/dist/index.html`.

## Why?

Post-incident analysis often asks:

* What happened?
* What failed?
* What would have stopped it?

That last question is easier in hindsight than the prospective problem:
**what will the next incident look like?**

One control may stop a single attack path while leaving another open. Excellent detection may produce no effective response. A late intervention may still prevent the outcome that matters.

Muster makes those counterfactuals executable.

## A counterintuitive result

The toy scenario contains a deliberately constructed example of
**non-monotone defense**.

One Watchline detects later Kubernetes discovery and uses that evidence to revoke
an already-stolen credential. That composition prevents node access.

Add an otherwise useful early worker-isolation response, however, and the result
gets worse: isolation removes the worker access needed for the later discovery
event. The evidence needed to trigger credential revocation is never generated,
the stolen credential survives, and the attacker reaches the node.

This is a **synthetic existence proof**, not a claim about the Hugging Face
incident, and not evidence that early isolation generally makes incidents
worse. The mechanism is the point:

> **Containment is not always causally independent of investigation.**

Incident responders know the tension: NIST SP 800-61 lists "need for
evidence preservation" among the criteria for choosing a containment
strategy. Muster turns that judgment call into something you can check
against a specific trace: which response loses which evidence, and what that
leaves reachable.

The same mechanism is written into professional guidance in other fields.
Fire-investigation guidance warns that overhaul (tearing out smouldering
material so a fire cannot rekindle) done before an investigator looks can
destroy the evidence of how a fire started. Sepsis guidelines ask for blood
cultures before the first antibiotic dose where that causes no real delay,
because antibiotics first lower the chance of ever identifying the organism.
Both are modeled as labs (`examples/fire-overhaul-evidence`,
`examples/sepsis-cultures`), and both show the same shape, along with the
tradeoff that makes it hard: the early action buys something real.

You can reproduce the comparison directly:

```bash
go run . experiment --file experiments/toy-jackpot.yaml
```

## Watchline

Watchline groups defenses by role rather than product or infrastructure layer:

**🔭 Vedette — observe**
Detect, correlate, classify, or surface evidence.

**🛡️ Picket — intervene**
Prevent or interrupt attacker transitions.

**🔧 Reserve — recover**
Contain compromise, revoke access, rebuild, or restore state.

These roles are **interleaved**, not a pipeline. A system may have all three at many different layers.

Roughly, they map to familiar **Detect / Protect / Respond** activities. The
point of the taxonomy is causal role, not new names for existing products.

Muster lets you test those placements against an incident path.

> **Well-placed Pickets for what you expect. Broad Vedettes for what you might
> notice. Robust Reserves for what you got wrong.**

## Try Muster

The easiest way to understand Muster is to **build a Watchline and replay an incident**.

The quickest route is the [browser lab](https://bfaber-centaur.github.io/muster/).
Nothing to install.

The original Jupyter lab is still here too: open [`notebooks/muster_lab.ipynb`](notebooks/muster_lab.ipynb).

Start with the **Toy incident**. Pick defenses, run the replay, and see what the attacker can still reach. Change your Watchline and try again.

No Go or Python knowledge is required to use the lab.

### Running the notebook

You’ll need **Go 1.22+**, **Python 3**, and Jupyter.

<details>
<summary><strong>Need to install Go or Python?</strong></summary>

### macOS

With [Homebrew](https://brew.sh/):

```bash
brew install go python
```

### Windows

With `winget`:

```powershell
winget install GoLang.Go
winget install Python.Python.3
```

### Ubuntu / Debian

```bash
sudo apt update
sudo apt install golang-go python3 python3-pip
```

Check installation:

```bash
go version
python3 --version
```

On Windows, use `python` if `python3` is unavailable.

</details>

Clone Muster:

```bash
git clone https://github.com/bubskee/muster
cd muster
```

Install Jupyter and launch the notebook:

```bash
python3 -m pip install jupyter ipywidgets
python3 -m jupyter notebook
```

Navigate to:

```text
notebooks/muster_lab.ipynb
```

Then choose **Run All**.

Once the toy model feels intuitive, switch to the **HF-inspired public-record abstraction**.

## Scenarios

Each scenario is a lab: a compressed trace, a pool of defenses, and
alternative versions of events. Public-record labs trace each step to a
published account. Illustrative labs model a mechanism documented in a
field's own guidance, not one specific event. Each folder's README lists its
sources, what the compression keeps and drops, and its caveats.

| Scenario | Field | Kind | What it shows |
| --- | --- | --- | --- |
| [Toy incident (JACKPOT)](examples/toy-jackpot/) | cyber | synthetic | Early isolation removes the evidence a later credential revocation needs. |
| [Hugging Face intrusion, July 2026](examples/hf-july-2026/) | cyber | public record | Only early paging plus worker isolation holds in every version; prevention depends on which Pickets happened to hold. |
| [Claude Mythos 5 / PyPI incident, 2026](examples/anthropic-pypi-2026/) | AI agents | public record | The evaluation operator's defenses hold everywhere; a vendor's own defenses protect only that vendor. |
| [CircleCI stolen session, 2023](examples/circleci-2023/) | cyber | public record | The 2FA in place did nothing: the stolen session had already passed it. The real detection came from a customer, after harm. |
| [Replit agent deletes a production database, 2025](examples/replit-agent-2025/) | AI agents | public record | Backups only help if the user doubts the agent's claim that rollback is impossible. |
| [Bangladesh Bank SWIFT heist, 2016](examples/bangladesh-bank-2016/) | finance | public record | The biggest save was luck: a sanctions hit on a street name. Without it, only a real-time hold at the Fed saves the bulk of the money. |
| [Deepwater Horizon / Macondo, 2010](examples/deepwater-horizon-2010/) | engineering | public record | The barrier interaction the CSB documents never changes an outcome; only the early exits survive every version. |
| [Sepsis: cultures and the first dose](examples/sepsis-cultures/) | medicine | illustrative | Antibiotics first trades a late dose for a lost organism, and becomes a strict backfire when broad therapy can never be narrowed. |
| [Fire overhaul vs. arson evidence](examples/fire-overhaul-evidence/) | fire | illustrative | Routine overhaul destroys the evidence an investigator needs; holding overhaul without a fire watch adds a rekindle. |

The defenses in each lab are counterfactuals unless the README says they were
deployed. Scenarios outside security use the same three roles: a
negative-pressure test is a Vedette, a shoe-track cement barrier is a Picket,
diverting flow overboard is a Reserve.

## Hugging Face example

The main real-world example is a compressed abstraction of the **July 2026 Hugging Face intrusion**, based on Hugging Face's published technical timeline.

It is **not a forensic reconstruction**. It preserves enough structure to compare defensive counterfactuals across several attack branches. Those branches share an early foothold, and the cluster-admin branch also depends on the node-root secret read.

| Watchline                            | Node root | Internal network | Cluster admin | Worker credentials |
| ------------------------------------ | --------: | ---------------: | ------------: | -----------------: |
| observed signal, no page             |       yes |              yes |           yes |                yes |
| privileged/hostPath admission policy |        no |               no |            no |                yes |
| cluster-scoped connector identity    |       yes |              yes |            no |                yes |
| critical page + worker isolation     |        no |               no |            no |                yes |
| dataset local-file-read hardening    |       yes |              yes |           yes |                 no |

Different controls buy different pieces of the incident.

Blocking privileged workloads removes node root and everything downstream of
it, including cluster-admin, because the connector credential came from the
node-root secret read. Scoping the connector credential removes only
cluster-admin. Early escalation and worker isolation cuts every downstream
branch. None of these touch the credentials exposed by the dataset file-read
before any alert fired; only hardening that read does, and it leaves the rest
of the intrusion intact.

See [`examples/hf-july-2026/`](examples/hf-july-2026/) for provenance and modeling caveats.

## Anthropic PyPI example

A second public-record abstraction follows Anthropic's **2026 Claude Mythos 5 /
PyPI incident** across three defensive boundaries: an evaluation environment, a
public package registry, and a third-party security vendor.

The example asks where a Watchline could interrupt the progression from
unexpected public-internet access, to package publication, to third-party
execution, credential exposure, and live-database access.

See [`examples/anthropic-pypi-2026/`](examples/anthropic-pypi-2026/) for the
compressed scenario, controls, provenance, and caveats.

```bash
go run . experiment --file experiments/anthropic-pypi-2026.yaml
```

## Backfire search and robustness

Each scenario has a **lab manifest** (`lab.yaml`) beside its scenario and
control pool. The manifest names the adverse outcomes, groups engine controls
into human-facing choices, defines presets, and lists **alternative versions
of events**. It only selects; replay semantics stay in the scenario YAML and
the Go engine.

```bash
go run . lab    --file examples/toy-jackpot/lab.yaml --preset "Strong but brittle"
go run . search --file examples/toy-jackpot/lab.yaml
go run . robust --file examples/hf-july-2026/lab.yaml
go run . survey
```

**`search`** replays every subset of a lab's choices and reports the smallest
cases where adding one defense makes the outcome strictly worse: the adverse
set after is a strict superset of the set before. Nobody has to think of the
JACKPOT pair first. On the toy, the search finds it unprompted, along with a
second form of it: once the early-isolation response is in place, *adding the
detector that triggers it* also makes things worse. `--tradeoffs` also lists
defenses that remove one harm and open another.

**`robust`** replays every preset against the recorded trace and each
alternative version, and checks whether each backfire survives. Alternatives
come in three kinds:

* **failed attempt**: the record shows something the adversary tried that a
  defense stopped. The variant asks what happens if that defense had not held.
  A replay of a known incident takes these defenses for granted.
* **alternative branch**: a plausible different route or ordering.
* **perturbation**: a change that should not matter, as a control check.

For example, Hugging Face's timeline lists attempts that failed: Hub
production database connections that "never completed" behind a private-link
allowlist, and cloud-role calls that were "denied by the role's policy". Taking
those away changes which Watchlines still hold. Only early paging plus worker
isolation gives the same result in every version, because it acts at the
first foothold. Even layered prevention fails once the private-link allowlist
does. See [`examples/hf-july-2026/`](examples/hf-july-2026/).

**`survey`** runs both over every lab:

<!-- survey:start -->
| Lab | Domain | Kind | Watchlines | Backfires (recorded / any variant) | Tradeoffs (recorded) | Fragile presets |
| --- | --- | --- | ---: | ---: | ---: | --- |
| Toy incident (JACKPOT) | cyber | synthetic | 64 | 6 / 8 | 0 | 1 of 3 |
| Hugging Face intrusion, July 2026 | cyber | public-record | 64 | 0 / 0 | 0 | 5 of 6 |
| Claude Mythos 5 / PyPI incident, 2026 | ai-agents | public-record | 128 | 0 / 0 | 0 | 4 of 6 |
| CircleCI stolen session (January 2023) | cyber | public-record | 1024 | 0 / 0 | 0 | 4 of 6 |
| Replit agent deletes a production database | ai-agents | public-record | 512 | 0 / 0 | 0 | 3 of 5 |
| Bangladesh Bank SWIFT heist (2016) | finance | public-record | 256 | 0 / 0 | 0 | 5 of 6 |
| Deepwater Horizon / Macondo blowout (2010) | engineering | public-record | 1024 | 0 / 0 | 0 | 3 of 6 |
| Sepsis, cultures and the first dose | medicine | illustrative | 64 | 0 / 2 | 8 | 5 of 5 |
| Fire overhaul vs. arson evidence | fire | illustrative | 64 | 5 / 7 | 1 | 5 of 5 |
<!-- survey:end -->

A few things this table says:

* **No public-record incident here contains a backfire.** On the recorded
  traces and every alternative, adding a defense never made things strictly
  worse. In these incidents, ordinary dependency reasoning gives the replay's
  answer.
* **Backfires appear exactly where a response destroys evidence a later
  response needs**: the synthetic toy, fire overhaul, and sepsis care. In the
  last two the mechanism comes from published guidance, not from the model's
  construction. In sepsis the recorded case is a *tradeoff* (antibiotics first
  avoid a late dose but lose the organism); it becomes a strict backfire when
  the organism turns out to be susceptible and broad therapy can never be
  narrowed.
* **Most presets are fragile.** In most labs, most Watchlines that work on the
  recorded incident give a different result in at least one alternative
  version. The recorded trace flatters defenses placed where the adversary
  happened to go.

## Model semantics and limits

Muster is intentionally small and **attack-graph-like**. The claim is not that
these incidents cannot be represented with attack graphs, state machines,
Petri nets, or other formalisms. The value of the harness is making defensive
state and counterfactual composition directly executable and inspectable.

The model is simple:

* **State** is a set of facts.
* **Events** have preconditions and add/remove effects.
* **Controls** are evaluated around events and may observe, block, review,
  escalate, or respond.
* **Reserves** can change persistent state, so defensive history can affect
  later reachability.
* Only controls that act count. A Picket that is suppressed or still waiting
  on a requirement does not block its event. (Before October 2026 it did; no
  published result depended on that, but one Gate 2 summary reported an
  observation by a suppressed Vedette.)

Deliberate limitations:

* event order and attack relationships are supplied by the scenario, and
  alternative versions are supplied by hand, not generated;
* replay is single-pass, and skipped events are not retried;
* there is no probability, continuous time, attacker adaptation, or automatic
  cost model;
* a known trace gives matched Pickets an ex-post advantage defenders do not
  have prospectively;
* Muster can test relationships represented in the model, but it cannot
  discover an attack relationship nobody modeled.

> **You can test whether your Watchline survives the attacks you imagined.
> You cannot replay the attack you failed to imagine.**

Think of Muster as a **counterfactual reasoning aid**, not a cyber range,
attacker emulator, or full threat model.

## CLI and reproducible experiments

The notebook sits on top of the Go replay engine. Technical users can run experiments directly:

```bash
go run . experiment --file experiments/hf-july-2026.yaml
```

JSON and CSV output are also available:

```bash
go run . experiment \
  --file experiments/hf-july-2026.yaml \
  --format json

go run . experiment \
  --file experiments/hf-july-2026.yaml \
  --format csv
```

Individual scenario/control pairs can be replayed directly:

```bash
go run . replay \
  --scenario examples/gate2-between.yaml \
  --controls examples/gate2-broad-correlated.yaml
```

Scenario configuration lives in YAML; replay semantics remain in Go.

Lab commands work on a lab manifest:

```bash
go run . lab    --file LAB.yaml [--choices V1,R1 | --preset LABEL] [--variant ID] [--json]
go run . search --file LAB.yaml [--variant ID] [--tradeoffs] [--json]
go run . robust --file LAB.yaml [--json]
go run . survey [--root examples]
go run . bundle [--root examples] [--out labs.json] [--golden golden.json]
```

The browser lab uses a JavaScript port of the replay engine. `./web/build.sh`
writes every replay of every choice subset of every lab from Go
(`bundle --golden`), replays the same in the browser engine, and fails if any
event status, acting control, terminal state, search finding or robustness
cell differs.

## Direction

Existing frameworks describe attacker techniques, defensive controls, and response playbooks.

Muster explores one layer between them:

```text
incident / playbook
        ↓
defensive actions
        ↓
Watchline roles
        ↓
counterfactual replay
        ↓
what did this defense actually buy us?
```

Longer term, the question is whether this helps compare defensive portfolios, allocate scarce attention, and decide where the next watchline belongs.

For now, it is a small executable argument.

---

> “Everything is very simple in war, but the simplest thing is difficult.”
>
> — Carl von Clausewitz, *On War*

---

Built for the **Apart Research Incident Response Sprint**.
