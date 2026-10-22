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
incident. The mechanism is the point:

> **Containment is not always causally independent of investigation.**

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

Open [`notebooks/muster_lab.ipynb`](notebooks/muster_lab.ipynb).

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

Deliberate limitations:

* event order and attack relationships are supplied by the scenario;
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
