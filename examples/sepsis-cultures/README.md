# Sepsis: blood cultures and the first antibiotic dose (illustrative)

> **Educational model, not clinical guidance.** This lab compresses a
> well-known mechanism from published studies and guidelines into a few binary
> facts so that Muster can replay it. It is not a description of a real
> patient. It leaves out almost everything that matters at a real bedside, and
> it must not be used to make or justify treatment decisions.

`kind: illustrative`. The mechanism comes from the literature. The patient,
the timeline and the event ordering are constructed, and no single published
case is being reconstructed.

## The mechanism

In suspected sepsis, two good things compete for the same few minutes:

- **Antibiotics fast.** In New York State's mandated sepsis-care data, a longer
  time to antibiotics was associated with higher risk-adjusted in-hospital
  mortality (Seymour et al. 2017).
- **Cultures first.** A blood culture drawn after antibiotics have started is
  much less likely to grow the organism (Cheng et al. 2019). With no culture
  result, nobody learns what the organism is or which drugs it responds to. That
  knowledge is what the later, targeted step needs: switching to a drug that
  works on a resistant organism, or narrowing broad antibiotics.

This is a clinical version of the pattern in Muster's toy JACKPOT scenario. An
early intervention that is good on its own terms (antibiotics now) removes the
evidence (a culture that will grow) that a later, targeted intervention (switch
or narrow therapy) depends on. The difference from the toy is that here the
early intervention also buys something real, namely time. The guideline answer
is to sequence the two steps, not to choose between them: the Surviving Sepsis
Campaign says to take cultures before antimicrobials when doing so causes no
substantial delay (< 45 min).

## Sources (each one opened)

- **[CHENG]** Cheng MP et al. *Blood Culture Results Before and After
  Antimicrobial Administration in Patients With Severe Manifestations of
  Sepsis: A Diagnostic Study.* Ann Intern Med 2019;171(8):547–554 (FABLED).
  PubMed record: https://pubmed.ncbi.nlm.nih.gov/31525774/. PubMed and the
  publisher page did not render for me (cookie wall / HTTP 403), so I read the
  numbers below from the abstract as reproduced by IVTEAM:
  https://ivteam.com/intravenous-literature/blood-culture-results-before-and-after-antimicrobial-administration/
  - Single-group diagnostic study in 7 North American emergency departments.
    Adults with severe manifestations of sepsis (systolic BP < 90 mm Hg or
    lactate ≥ 4 mmol/L). Post-antimicrobial cultures were drawn within 120
    minutes of starting treatment (median 70 min).
  - Pre-antimicrobial blood cultures were positive in 102/325 (31.4%) and
    post-antimicrobial cultures in 63/325 (19.4%), an absolute difference of
    12.0% (95% CI 5.4–18.6%). Sensitivity of post-antimicrobial culture was
    52.9% (95% CI 42.8–62.9%). With other microbiological cultures added,
    pathogens were found in 69/102 (67.6%).
  - Secondary results and limitations, from The Bottom Line's summary
    (https://www.thebottomline.org.uk/summaries/icm/fabled/): when the
    organism was **not** susceptible to the antibiotic given, the
    post-antibiotic culture was positive in 66.7%, against 51.7% when it was
    susceptible. The study had no comparator group, so it cannot measure how
    much drawing cultures first delays antibiotics. Only some screened patients
    were enrolled.
- **[SEYMOUR]** Seymour CW et al. *Time to Treatment and Mortality during
  Mandated Emergency Care for Sepsis.* N Engl J Med 2017;376:2235–2244.
  https://pmc.ncbi.nlm.nih.gov/articles/PMC5538258/ — a retrospective cohort
  of 49,331 patients at 149 New York hospitals. Each hour to antibiotic
  administration: OR 1.04 per hour (95% CI 1.03–1.06) for risk-adjusted
  in-hospital mortality. Each hour to completion of the fluid bolus: OR 1.01
  (0.99–1.02, P = 0.21). The study is observational.
- **[SSC]** Evans L et al. *Surviving Sepsis Campaign: International
  Guidelines for Management of Sepsis and Septic Shock 2021.* Intensive Care
  Med 2021. https://pmc.ncbi.nlm.nih.gov/articles/PMC8486643/
  - Rec 12: for possible septic shock or a high likelihood of sepsis, give
    antimicrobials immediately, ideally within 1 h of recognition. Rec 14: for
    possible sepsis without shock, within 3 h.
  - The rationale text says to obtain cultures before antimicrobials "if it
    results in no substantial delay in the start of antimicrobials (i.e. < 45
    min)".
  - Recs 17–20 make empiric MRSA and double gram-negative cover depend on the
    patient's risk. The rationale notes that failing to cover MRSA in a patient
    with MRSA may be harmful, and that unnecessary MRSA cover may also be
    harmful.
  - Rec 29 advises daily assessment for de-escalation. The text describes
    de-escalation "once both the pathogen(s) and susceptibilities are known",
    and says thoughtful de-escalation based on clinical improvement is
    appropriate even if cultures are negative.

## The compressed trace

Baseline patient: a bloodstream infection with an organism that the standard
first-line ("empiric") regimen does **not** cover. The baseline trace is
**usual workup without a sepsis fast-track**: cultures are drawn during the
workup and the first dose follows it, which the model treats as past the first
hour.

| # | Event | Source |
|---|---|---|
| 1 | Sepsis recognised | — |
| 2a | Fast-track first dose, after cultures (only if V1 drew them) | SSC Rec 12 |
| 2b | Fast-track first dose, before any culture | SSC Rec 12 |
| 3a | Blood cultures drawn, no antibiotic yet | CHENG |
| 3b | Blood cultures drawn after the first dose | CHENG (median 70 min) |
| 4 | First dose after the routine workup (adds *dose came late*) | SEYMOUR (binary stand-in) |
| 5 | Culture grows the organism (pre-antibiotic sample only) | CHENG |
| 6 | Switch to a drug the organism is known to respond to (also narrows broad cover) | SSC |
| 7 | Bedside reassessment after a day or two | SSC Rec 29 |
| 8 | Infection returns after narrowing blind (resistant organism only) | **inference** |

The paired events (2a/2b, 3a/3b) are mutually exclusive by precondition. This
is how the model expresses "which came first" in a single left-to-right replay.

Choices:

| Key | Role | Choice | Engine control |
|---|---|---|---|
| V1 | Vedette | Draw cultures at recognition (before any dose) | observe on 1 |
| P1 | Picket | Hold antibiotics until cultures are drawn (no time limit) | block on 2b |
| R1 | Reserve | Fast-track antibiotics within the hour | respond on 1 |
| R2 | Reserve | Broad cover from the first dose | respond on 2a/2b/4 |
| R3 | Reserve | Broaden at review if not improving and no culture | respond on 7 |
| R4 | Reserve | Narrow at review on improvement, even without a culture | respond on 7 |

Adverse outcomes: **O1** first dose came late. **O2** organism never
identified. **O3** infection not covered by the antibiotics given. **O4** still
on broad antibiotics, never narrowed.

### What the compression preserves

- Antibiotics before cultures lower the chance that the culture grows.
- Without a culture result there is no targeted switch and no
  susceptibility-guided narrowing.
- Strict cultures-first can delay antibiotics, which is why the guideline
  limits it to "no substantial delay".
- Broad up-front cover depends on cultures if it is ever to be narrowed safely.

### What it omits or distorts

- **Everything is binary and deterministic.** Real culture sensitivity after
  antibiotics is about 53%, not 0%. The baseline is the false-negative half;
  the variant `grows-despite-antibiotics` is the other half.
- **Delay is one bit.** O1 stands in for a continuous, population-level,
  observational association (OR 1.04 per hour). It is not a statement about
  this patient. For a resistant organism the first-line dose would not have
  covered the infection anyway, but the clinician cannot know that at the
  bedside.
- **The model applies only to the bacteremic minority.** In CHENG, about 69%
  of patients had negative cultures even before antibiotics.
- **Time to effective therapy is not scored.** R3 "rescues" coverage at about
  48 h. The model does not count that the patient spent those hours on
  ineffective therapy, so R3 looks better than it should.
- **R3's broadening and R2's broad cover are assumed to work.** R2 is applied
  unconditionally, whereas the guideline makes it depend on risk. Event 8
  (relapse after blind narrowing) is my inference about the mechanism, not a
  finding in the cited sources.
- **Other sites are left out.** Cultures from other sites (urine, sputum),
  rapid molecular tests, source control, fluids beyond the perturbation,
  allergies, toxicity, *C. difficile* and resistance selection are all
  omitted.

## Variants

| id | kind | What changes | Why |
|---|---|---|---|
| `susceptible-organism` | alternative-branch | First-line regimen covers the organism | If the reversal is driven by resistance, the inadequate-therapy harm should vanish |
| `grows-despite-antibiotics` | alternative-branch | Post-dose culture still grows | About half of post-dose cultures still grow in CHENG (66.7% when the organism was not susceptible to the drug given). Losing the evidence is not certain |
| `difficult-access` | alternative-branch | Cultures cannot be drawn quickly at recognition (suppresses V1) | The case the SSC "< 45 min" caveat exists for |
| `fluid-bolus` | perturbation | IV fluid bolus inserted after recognition | Nothing reads it. SEYMOUR found no mortality association with bolus timing. Results should not move |

There is no `failed-attempt` variant. This is an illustrative mechanism, not
an incident record, so there is no recorded attempt that was denied.

## Results (frozen build, `muster-bin`)

### `muster robust`

```
# Robustness: Sepsis, cultures and the first dose

Adverse outcomes reached by each Watchline, as recorded and under each alternative trace.
↑ = worse than as recorded, ↓ = better, ≠ = different.

| Watchline | As recorded | First-line antibiotics cover the organism | Culture grows even after the first dose | Cultures can't be drawn quickly | IV fluid bolus added |
| --- | --- | --- | --- | --- | --- |
| No added defense | O1 | O1 | O1 | O1 | O1 |
| Antibiotics first | O2, O3 | O2 ↓ | none ↓ | O2, O3 | O2, O3 |
| Cultures first, then antibiotics within the hour | none | none | none | O2, O3 ↑ | none |
| Strict cultures-first | none | none | none | O1 ↑ | none |
| Fast and broad | O2, O4 | O2, O4 | none ↓ | O2, O4 | O2, O4 |
| Fast, broad and reassess | O2, O3 | O2 ↓ | none ↓ | O2, O3 | O2, O3 |

Outcomes:
- O1 = first antibiotic dose came late
- O2 = organism never identified
- O3 = infection not covered by the antibiotics given
- O4 = still on broad antibiotics, never narrowed

## Do the baseline non-monotone findings survive?


## New non-monotone cases under “First-line antibiotics cover the organism”

- [R1 Fast-track antibiotics] + **R2 Broad cover from the first dose** → also still on broad antibiotics, never narrowed

## New non-monotone cases under “Culture grows even after the first dose”

- [R1 Fast-track antibiotics] + **P1 Hold antibiotics until cultures are drawn** → also first antibiotic dose came late
```

### `muster search` (key findings; full output with `--tradeoffs`)

- **Baseline: 0 strictly-worse cases, 8 tradeoffs.** The headline tradeoff is
  `[no defense] + R1 Fast-track antibiotics`: before, *first dose came late*;
  after, *organism never identified; infection not covered*. The first
  divergence is at "Sepsis recognised". The mirror image is
  `[R1] + P1 Hold antibiotics`, which turns that pair back into *dose came
  late*. Other tradeoffs: adding R2 or R3 to R1 buys coverage at the cost of
  *never narrowed*, and adding R4 to `[R1, R2]` narrows blind and the
  resistant infection returns.
- **`susceptible-organism`: 1 strictly-worse case.** `[R1] + R2 Broad cover`
  goes from *organism never identified* to *organism never identified; still
  on broad antibiotics, never narrowed*. Broad cover given before any culture
  has nothing to be narrowed by. Without R1, the culture is drawn first and
  event 6 narrows it, so R2 does no harm. The harm comes from the R1 × R2
  interaction.
- **`grows-despite-antibiotics`: 1 strictly-worse case.** `[R1] + P1 Hold
  antibiotics` goes from *none* to *dose came late*. When the evidence
  survives the first dose anyway, strict cultures-first costs time and buys
  nothing.
- **`difficult-access`** has the same search findings as the baseline. The
  robust table shows the change: the guideline ordering `[V1, R1]` collapses
  to the antibiotics-first outcome (O2, O3), and strict cultures-first
  collapses to the late dose (O1).
- **`fluid-bolus`** (perturbation): identical to the baseline in both robust
  and search, as intended.

### Interpretation

- **Does a JACKPOT appear? Partly, and honestly as a tradeoff.** In the
  baseline, adding "antibiotics first" never makes the result strictly worse,
  because it always removes *dose came late*. It swaps a timing harm for an
  evidence harm and the downstream inadequate-therapy harm. In this model that
  is the accurate statement: the early intervention is not free, and it is not
  useless either.
- **Sequencing removes the conflict.** `[V1, R1]` (cultures at recognition,
  then the dose within the hour) has no adverse outcome in the baseline. Under
  the model this is why the guideline asks for both. The conflict returns only
  when cultures cannot be taken quickly (`difficult-access`). Then neither
  strict rule is safe: one gives up the evidence, the other gives up the time.
- **The reversal behaves as the brief predicted for the right reason.** With a
  susceptible organism, *infection not covered* disappears from every
  Watchline. What is left of the evidence loss is *organism never identified*,
  plus the *never narrowed* harm if broad cover was used. The strictly-worse
  `R1 + R2` case there is grounded in the de-escalation logic of SSC. It
  matches the toy mechanism: evidence lost early disables a later corrective
  step.
- **The strict cultures-first Picket is the one control that becomes "worse"
  outright** (`grows-despite-antibiotics`). The guideline's 45-minute limit
  exists to prevent that pure-delay case.

## Caveats

- This is an educational model and must not be used for clinical decisions.
- The numbers are quoted from the CHENG abstract (as reproduced by IVTEAM), a
  secondary summary (The Bottom Line), the SEYMOUR abstract on PMC, and the
  SSC 2021 guideline on PMC. I could not read the CHENG full text.
- SEYMOUR is observational. It shows an association, not proven causation.
  CHENG had no comparator arm, so it cannot measure
  the delay that cultures-first causes.
- Binary harms cannot be weighed against each other. Muster reports
  "tradeoff" and does not say which side is clinically larger. That judgement
  depends on the patient's severity, their risk of resistant organisms, and
  how long cultures would really take.
- Results depend on modelling choices that are marked as inference: event 8,
  R3 always succeeding, R2 being unconditional, and the baseline dose coming
  after the first hour.
- Engine note: this lab was first run on a frozen build that predated commit
  723e058 ("Only credit controls that acted..."). In that build a Picket still
  waiting on its `requires` would still block. The lab now gives its only
  Picket (P1) no `requires`, so it does not depend on that behaviour. The
  output above comes from the rebuilt `muster-bin`, and robust and search
  output is byte-identical across the old frozen build, the rebuilt build and a
  build of the current source.

## Run it

```bash
BIN=muster   # or the frozen build
$BIN robust --file examples/sepsis-cultures/lab.yaml
$BIN search --file examples/sepsis-cultures/lab.yaml --tradeoffs
$BIN search --file examples/sepsis-cultures/lab.yaml --variant susceptible-organism
$BIN lab    --file examples/sepsis-cultures/lab.yaml --preset "Antibiotics first"
$BIN lab    --file examples/sepsis-cultures/lab.yaml --preset "Cultures first, then antibiotics within the hour" --variant difficult-access
```
