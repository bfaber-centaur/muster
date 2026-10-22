# Hugging Face July 2026 — compressed public-record example

This example is a deliberately small **public-record abstraction**, not a
forensic reconstruction of the July 2026 Hugging Face intrusion.

Primary source:

- Hugging Face, *Anatomy of a Frontier Lab Agent Intrusion: A Technical
  Timeline of the July 2026 Incident* (2026-07-27):
  https://huggingface.co/blog/agent-intrusion-technical-timeline

Useful independent defensive mapping:

- Elastic Security Labs, *Exploring the Hugging Face Breach: mapping AI agent
  tactics to Elastic Defend* (2026-07-31):
  https://www.elastic.co/security-labs/threat-command/ai-agent-attack-detection-hugging-face-breach

## What the scenario preserves

The ten semantic transitions preserve several public-record relationships that
matter for defensive counterfactuals:

1. Malicious dataset configuration produced local-file disclosure from the
   production worker, exposing environment credentials and worker source.
2. A separate template-injection vector produced arbitrary code execution in
   the worker pod.
3. The worker foothold enabled Kubernetes and cloud discovery.
4. A privileged pod with a host filesystem mount enabled node root. Hugging
   Face states that no admission policy rejected privileged/hostPath pods.
5. Node root enabled cluster-secret theft and then mesh-network enrollment.
6. The compromised worker could also resolve and reach an internal service
   connector directly, a different route from the mesh pivot. The credential it
   replayed there came from the same node-root secret read as the mesh key
   ("That single read yielded both the mesh-VPN auth key used minutes later and
   the access-broker credential used an hour later"). A shared connector
   credential bound to `system:masters` yielded cluster-admin across multiple
   clusters. **So the connector branch is a separate route, not an independent
   path: in this model it still depends on node root.**
7. Hugging Face reports that security signals were correlated into a coherent
   attack signal but failed to receive the correct criticality and on-call
   escalation.

The model intentionally omits action counts, exact timestamps, C2 mechanics,
exfiltration, supply-chain exploration, denied probes, persistence, and many
other details. Those details are important to the incident but unnecessary for
this demonstration.

## The runs

`observed-no-page.yaml` abstracts the disclosed detection/escalation failure: a
Vedette records worker compromise, but no response follows. It does **not**
claim that the real signal came from the exact event modeled here.

The other Watchlines are explicitly counterfactual:

- `admission-policy.yaml` blocks the privileged/hostPath pod. Because the
  connector credential came from the node-root secret read, this also removes
  the cluster-admin branch. It does nothing about credentials the dataset
  file-read exposed earlier.
- `connector-scope.yaml` blocks cross-cluster replay of the shared connector
  identity. It stops cluster-admin while leaving node-root and mesh compromise
  possible.
- `page-and-isolate.yaml` models a correctly prioritized alert followed by
  worker isolation. It cuts every downstream branch after the worker RCE, but
  the earlier file-read exposure has already happened.
- `hdf5-hardening.yaml` blocks the dataset local-file read. It protects the
  worker credentials and leaves the rest of the intrusion intact.

> **Correction (October 2026).** Earlier versions of this README and the
> top-level README described the connector path as independent of node root
> and showed the admission policy leaving cluster-admin reachable. The
> scenario never encoded that. The connector replay has always required
> `attacker:cluster-secrets`, which matches Hugging Face's timeline. The table
> below is the engine's actual output, and matches the sprint report's Table 1.
> The independent-path reading is now an explicit variant
> (`variants/connector-cred-on-worker.yaml`) so the two can be compared.

## Run it

From the repository root:

```bash
go run . experiment --file experiments/hf-july-2026.yaml
```

For machine-readable output:

```bash
go run . experiment --file experiments/hf-july-2026.yaml --format json \
  > runs/hf-july-2026.json

go run . experiment --file experiments/hf-july-2026.yaml --format csv \
  > runs/hf-july-2026.csv
```

Output (`go run . experiment --file experiments/hf-july-2026.yaml`):

| Watchline | node root | internal network | cluster admin | worker credentials exposed |
|---|---:|---:|---:|---:|
| observed signal, no page | yes | yes | yes | yes |
| privileged/hostPath admission policy | no | no | no | yes |
| cluster-scoped connector identity | yes | yes | no | yes |
| critical page + worker isolation | no | no | no | yes |
| dataset local-file-read hardening | yes | yes | yes | no |

These are outputs of the **compressed model**, not empirical claims that any
single intervention would certainly have produced the same outcome in the
real incident.

## Variants: what if a defense that held had failed?

Hugging Face's timeline lists several attempts that **failed**. These are
Pickets that held in the real incident, which an after-the-fact replay can
easily take for granted. `lab.yaml` turns them into alternative traces:

| Variant | Kind | What changes |
|---|---|---|
| `private-link-fails` | failed attempt | Hub production DB connections "never completed" because only allow-listed private-link traffic was accepted. Here the allowlist does not hold. |
| `cloud-mutation-succeeds` | failed attempt | "Every mutating call was denied by the role's policy", and 84 credential-store operations were denied. Here the role is over-permissioned. |
| `connector-cred-on-worker` | alternative branch | The shared connector credential is reachable from the worker's own environment, so cluster-admin no longer needs node root. |

```bash
go run . robust --file examples/hf-july-2026/lab.yaml
```

| Watchline | As recorded | Private-link allowlist fails | Cloud role over-permissioned | Connector credential on the worker |
| --- | --- | --- | --- | --- |
| No added defense | O1, O2, O3, O4 | O1, O2, O3, O4, O5 ↑ | O1, O2, O3, O4, O6 ↑ | O1, O2, O3, O4 |
| What happened (signal, no page) | O1, O2, O3, O4 | O1, O2, O3, O4, O5 ↑ | O1, O2, O3, O4, O6 ↑ | O1, O2, O3, O4 |
| Admission policy | O4 | O4, O5 ↑ | O4, O6 ↑ | O3, O4 ↑ |
| Connector scoping | O1, O2, O4 | O1, O2, O4, O5 ↑ | O1, O2, O4, O6 ↑ | O1, O2, O4 |
| Critical page + isolate | O4 | O4 | O4 | O4 |
| Dataset local-file-read hardening | O1, O2, O3 | O1, O2, O3, O5 ↑ | O1, O2, O3, O6 ↑ | O1, O2 ↓ |
| Layered prevention | none | O5 ↑ | none | none |

O1 node root · O2 internal network · O3 cluster-admin · O4 worker credentials
exposed · O5 Hub production DB · O6 cloud credential store.

What this says, within the model:

- **Early page + isolation is the only Watchline whose result does not move.**
  It acts at the first worker foothold, so every branch that needs the worker
  is gone, including the branches that only open when another defense fails.
  It cannot undo the file-read exposure that came before the alert.
- **Prevention is sensitive to which Pickets happened to hold.** Even the
  layered all-Picket Watchline lets the attacker reach the Hub production
  database if the private-link allowlist fails, because nothing else covers
  that branch.
- **The admission policy's apparent coverage of cluster-admin depends on where
  the connector credential lived.** If it was on the worker, the admission
  policy no longer covers it, and file-read hardening suddenly does.
- `muster search` finds **no** non-monotone interaction in this lab, in the
  baseline or any variant. As in the sprint report, the HF case is the
  baseline where ordinary dependency reasoning is enough.

These variants are counterfactuals about defenses that held, not claims about
what the attackers would have done next.
