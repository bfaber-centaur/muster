# CircleCI, January 2023: a compressed public-record example

This lab is a small **public-record abstraction** of the CircleCI security
incident disclosed on January 4, 2023. It is not a forensic reconstruction.

## Sources

Primary (opened and read):

- CircleCI, *CircleCI incident report for January 4, 2023 security incident*
  (Rob Zuber, CTO, 2023-01-13):
  https://circleci.com/blog/jan-4-2023-incident-report/ (cited as **[IR]**)
- CircleCI, *CircleCI security alert: Rotate any secrets stored in CircleCI*
  (2023-01-04, with dated updates through 2023-03-13):
  https://circleci.com/blog/january-4-2023-security-alert/ (cited as **[SA]**)

Secondary (press, opened and read):

- The Stack, *CircleCI warns users to "immediately rotate all secrets" amid
  credential abuse evidence* (2023-01-05):
  https://www.thestack.technology/circleci-security-incident-breach/ (cited as **[TS]**)

## Checking the brief against the record

The brief summary holds up against [IR], with these corrections and
additions:

- **Dates.** Laptop infected Dec 16, 2022. First unauthorized production
  access and reconnaissance Dec 19. Exfiltration Dec 22, the last recorded
  unauthorized production activity. Customer report of suspicious GitHub OAuth
  activity Dec 29. That customer's token confirmed compromised Dec 30. GitHub
  OAuth rotation for all customers started Dec 31 and finished Jan 7. Scope
  and entry point found Jan 4. Disclosure Jan 4 to 5 [IR][SA].
- **Escalation route.** The engineer "had privileges to generate production
  access tokens". The attacker reached "a subset of our production systems"
  and "a subset of databases and stores" [IR]. It was not a separate
  privilege-escalation exploit.
- **Decryption is "potential".** [IR]: the data was encrypted at rest, but
  "the third party extracted encryption keys from a running process",
  *potentially* enabling decryption. CircleCI told customers to assume the
  secrets had been accessed. The lab models decryption as happening, and says
  so.
- **"A customer noticed" is confirmed, not just reported.** [IR]: "we were
  alerted to suspicious GitHub OAuth activity by one of our customers." [IR]
  does not name the customer, and neither does this lab.
- **Antivirus.** "The malware was not detected by our antivirus software"
  [IR]. The session was "a valid, 2FA-backed SSO session" [IR].
- **Response details the brief left out.** Jan 4: the employee's access was
  shut down, production access was removed for nearly all employees, and
  exposed production hosts were rotated. Jan 5 to 6: project and personal API
  tokens were revoked and Bitbucket tokens were rotated with Atlassian. Jan 7
  to 12: AWS notified customers whose keys were potentially affected [IR][SA].
  "Despite working with GitHub to increase API rate limits, the rotation
  process took time" [IR].
- **What customers had to do.** Customers had to "rotate keys at the source":
  "Simply removing them from CircleCI is not enough" [SA]. On 2023-03-13 the
  secret-finding tool was updated to find secrets "stored in CircleCI but may
  not be visible in the UI", such as variables on renamed or deleted projects
  [SA].
- **Downstream harm.** "Fewer than 5 customers have informed us of
  unauthorized access to third-party systems" as of Jan 13 [IR].
- **Nothing in [IR] was blocked.** The report records no attacker action that
  was denied or failed. The one failed attempt this lab uses comes from the
  press. [TS] reports that customers spotted improper use of AWS credentials
  with **Canary Tokens**: decoy keys that give an attacker nothing and alert
  their owner. A researcher shared an attacker IP "for AWS credentials that
  were stored in CircleCI". [TS] gives no date for this. A Dec 27 date
  appears in search snippets of other coverage, but I could not open that
  article, so the lab does not use the date.

## What the scenario preserves

Twelve events (`scenario.yaml`; each one is commented with its source):

1. Malware lands on the laptop. The deployed antivirus saw nothing, so it is
   modelled as absent rather than as a control.
2. The malware steals a session that **already passed 2FA**.
3. *(Modelling device, not in the record.)* `password-sign-in` is the
   sign-in transition that 2FA guards. The attacker never needs it, so it is
   always skipped. The deployed 2FA Picket (P1) sits here. This is the "strong
   Picket on the wrong transition".
4. The attacker reuses the session "in a remote location".
5. The attacker mints production access as the engineer.
6. The attacker exfiltrates the encrypted customer data.
7. The attacker extracts encryption keys from a running process.
8. *(Inferred.)* The attacker decrypts the data, producing three kinds of
   stolen material:
   - tokens CircleCI can rotate itself (GitHub/Bitbucket OAuth, API tokens);
   - customer secrets that only the customer can rotate at the source;
   - "forgotten" secrets on renamed or deleted projects. That such secrets
     were in the stolen subset is my inference from the March 2023 tool
     update.
9. The attacker misuses a customer's GitHub OAuth token. This is the external
   Vedette's moment (Dec 29).
10. *(Press, [TS].)* The attacker tries a customer's decoy AWS key. It gains
    nothing. Its position after step 9 is the conservative choice.
11. The attacker uses stolen secrets on customer systems (the "fewer than 5"
    reports). [IR] does not give timing. This step is placed **inside the
    rotation lag**.
12. *(Timeline marker.)* Customers rotate their secrets at the source, days
    after notice. Supporting evidence: the Jan 4 to 5 notice, AWS still
    notifying customers through Jan 12, and the March tool update.

CircleCI's response is in `controls.yaml`, not in the trace. Every signal
that reaches CircleCI runs the same staged chain: review, then critical, then
escalate (incident declared, customers told), then respond. Each Reserve
choice carries that chain, so a Vedette with no Reserve gives an alert and
nothing else.

**Omitted:** how the malware was delivered, exact timestamps, which
databases were reached, the attacker's infrastructure and IOCs, the Dec 29 to
Jan 4 investigation gap (the response here is instantaneous), and per-provider
rotation details beyond "CircleCI-controlled" versus "customer-controlled".

## Choices

| Key | Role | Choice | Status |
|---|---|---|---|
| V1 | Vedette | Customer reports odd GitHub activity | as deployed (Dec 29) |
| V2 | Vedette | Customer decoy key raises the alarm | press-reported. Whether it fed CircleCI's response is unknown |
| V3 | Vedette | Alert on unusual sessions in production | counterfactual ([IR] plans "expanded alerting") |
| V4 | Vedette | Endpoint detection catches the stealer | counterfactual (the antivirus missed it) |
| P1 | Picket | Two-factor at sign-in | as deployed |
| P2 | Picket | Sessions tied to the device | counterfactual ([IR] plans "reducing session trust") |
| P3 | Picket | Short-lived customer credentials (OIDC) | counterfactual ([SA] recommends OIDC) |
| R1 | Reserve | CircleCI cuts access and rotates its tokens | as deployed (Dec 31 to Jan 7) |
| R2 | Reserve | Customers rotate secrets at the source | as deployed, with lag |
| R3 | Reserve | Find hidden secrets too | late in the record (Mar 13, 2023) |

## Variants

All variants reuse `controls.yaml`.

- **`decoy-was-real`** (`failed-attempt`). The decoy key from [TS] was a
  real, long-lived cloud key. The branch that the decoy closed opens: the
  attacker reaches a customer cloud account (O9), and the tripwire signal
  disappears.
- **`on-device-session`** (`alternative-branch`, not in the record). The
  malware uses the stolen session on the engineer's own laptop instead of
  replaying it remotely. This tests which defenses depended on the session
  looking out of place. The unusual-network alert (V3) is suppressed in this
  branch, and device binding (P2) does not apply.
- **`rotation-beats-attacker`** (`alternative-branch`). Customers finish
  rotating before the attacker uses their secrets. The record does not settle
  this ordering, so the variant tests how much the outcome depends on it.

Event trace for the "What happened" preset:

```
lab: circleci-2023  trace: baseline  choices: [V1,P1,R1,R2]
  applied  Malware lands on an engineer's laptop
  applied  Malware steals the signed-in session
  skipped  Attacker signs in with password and two-factor
  applied  Attacker reuses the session from elsewhere
  applied  Attacker creates production access as the engineer
  applied  Attacker copies customer secrets out
  applied  Attacker takes encryption keys from a running program
  applied  Attacker decrypts the copied secrets
  applied  Attacker uses a stolen token on a customer's GitHub 🔭customer-spots-github-activity 🔧provider-cuts-access-and-rotates
  applied  Attacker tries a customer's decoy key 🔧provider-cuts-access-and-rotates
  applied  Attacker uses stolen secrets on customer systems
  applied  Customers rotate secrets, days after notice 🔧customers-rotate-at-source
adverse: customer secrets copied out (encrypted); keys to decrypt them taken; forgotten secrets (old or renamed projects) still valid; a customer's GitHub reached with a stolen token; customer systems reached with stolen secrets
```

## `muster robust`

```
# Robustness: CircleCI stolen session (January 2023)

Adverse outcomes reached by each Watchline, as recorded and under each alternative trace.
↑ = worse than as recorded, ↓ = better, ≠ = different.

| Watchline | As recorded | The decoy key was a real key | Session used from the infected laptop | Customers rotate before the attacker uses their secrets |
| --- | --- | --- | --- | --- |
| No added defense | O1, O2, O3, O4, O5, O6, O7, O8 | O1, O2, O3, O4, O5, O6, O7, O8, O9 ↑ | O1, O2, O3, O4, O5, O6, O7, O8 | O1, O2, O3, O4, O5, O6, O7, O8 |
| What happened | O1, O2, O6, O7, O8 | O1, O2, O6, O7, O8, O9 ↑ | O1, O2, O6, O7, O8 | O1, O2, O6, O7 ↓ |
| Two-factor alone | O1, O2, O3, O4, O5, O6, O7, O8 | O1, O2, O3, O4, O5, O6, O7, O8, O9 ↑ | O1, O2, O3, O4, O5, O6, O7, O8 | O1, O2, O3, O4, O5, O6, O7, O8 |
| Stop the copied session | none | none | O1, O2, O3, O4, O5, O6, O7, O8 ↑ | none |
| Catch it on the laptop | none | none | none | none |
| Tripwire instead of luck | O1, O2, O6, O7, O8 | O1, O2, O3, O4, O5, O6, O7, O8, O9 ↑ | O1, O2, O6, O7, O8 | O1, O2, O6, O7 ↓ |
| Clean-up with expiring keys | O1, O2, O7 | O1, O2, O7 | O1, O2, O7 | O1, O2, O7 |

Outcomes:
- O1 = customer secrets copied out (encrypted)
- O2 = keys to decrypt them taken
- O3 = attacker still inside production
- O4 = stolen CircleCI-issued tokens still valid
- O5 = stolen customer secrets still valid
- O6 = forgotten secrets (old or renamed projects) still valid
- O7 = a customer's GitHub reached with a stolen token
- O8 = customer systems reached with stolen secrets
- O9 = a customer cloud account reached with a stolen key
```

## `muster search` (baseline and every variant, with `--tradeoffs`)

```
# Search: circleci-2023 (baseline)

10 choices, 1024 Watchlines replayed.

Strictly worse after adding one defense: 0 minimal case(s)
Tradeoffs (fixes one outcome, opens another): 0

# Search: circleci-2023 (decoy-was-real)

10 choices, 1024 Watchlines replayed.

Strictly worse after adding one defense: 0 minimal case(s)
Tradeoffs (fixes one outcome, opens another): 0

# Search: circleci-2023 (on-device-session)

10 choices, 1024 Watchlines replayed.

Strictly worse after adding one defense: 0 minimal case(s)
Tradeoffs (fixes one outcome, opens another): 0

# Search: circleci-2023 (rotation-beats-attacker)

10 choices, 1024 Watchlines replayed.

Strictly worse after adding one defense: 0 minimal case(s)
Tradeoffs (fixes one outcome, opens another): 0
```

## Interpretation

- **2FA was present and irrelevant.** "Two-factor alone" gives the same
  result as no defense in every trace. The attacker never crosses the
  transition 2FA guards. This is a property of the modelling, but it is
  exactly what [IR] describes.
- **The real detection came after the harm.** V1 fires *on* the misuse of a
  customer's GitHub token, so O7 (GitHub reached) is present in every
  Watchline that relies on it. External Vedettes see only harm that has
  already happened downstream. The decoy key (V2) helps only when the attacker
  happens to touch it: in `decoy-was-real`, "Tripwire instead of luck" falls
  back to the no-defense result (O1 to O9), because no signal ever arrives.
- **The response that happened left two things open.** "What happened" still
  reaches O8 (customer systems reached during the rotation lag) and O6
  (forgotten secrets still valid). When customers rotate first
  (`rotation-beats-attacker`), O8 disappears. That is how much the result
  depends on downstream parties acting quickly. Only P3 (credentials that
  expire on their own) and R3 (the March sweep) close these regardless of
  ordering. "Clean-up with expiring keys" is the one response-side Watchline
  that gives the same result in every trace, but it still concedes exfiltration,
  the keys, and the first GitHub misuse.
- **Upstream fixes depend on the attacker's route.** Device-bound sessions
  (P2) give a clean result in three traces and fail completely in
  `on-device-session`. The same holds for the unusual-network alert (V3, not
  in a preset): `--choices V3,R1,R2 --variant on-device-session` reaches
  O1 to O8. Catching the stealer on the laptop (V4 with R1 and R2) is the only
  preset with no adverse outcome in any trace. That depends on detection the
  record says did not exist, and on an instant response.
- **No non-monotone interaction.** `search` finds no case where adding a
  defense makes things worse, and no tradeoffs, in any of the four traces.
  That is the honest result for this wiring. The incident's lesson is about
  placement and timing, not about interference between defenses.

## Caveats

- `password-sign-in` is a modelling device, not a recorded step. Without it
  the deployed 2FA would have nowhere to sit. It is always skipped.
- Decryption (step 8) is CircleCI's assumption ("potentially"), not a
  confirmed fact. The three-way split of stolen material is mine.
- The decoy-key step and its variant rest on one press article. Its timing is
  not established.
- The timing of the third-party accesses is not in [IR]. The baseline places
  them in the rotation lag and the variant reverses that. Neither ordering is
  confirmed.
- The response chain is instantaneous. In reality the Dec 29 report led to
  the Jan 4 shutdown six days later, and GitHub OAuth rotation took a week. The
  attacker's last recorded production activity was Dec 22, so this does not
  change the production outcomes here. It would matter for any faster
  attacker.
- "Delete the secret from CircleCI" is not a choice. [SA] says it does not
  revoke the attacker's copy. In this model it would be a Reserve with no
  effect.
- The counterfactual controls (V3, V4, P2, P3, R3 at notice time) are
  illustrative. They are not claims that any of them would have stopped this
  attacker. P2 in particular does nothing against malware that uses the
  session on the device itself, and the `on-device-session` variant shows
  this.
- Downstream harm is reduced to a few yes/no facts. "Fewer than 5" reports,
  and AWS saying it saw "no indication" of account access, are both coarser
  than this model.
