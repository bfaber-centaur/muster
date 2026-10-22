# Deepwater Horizon / Macondo blowout (April 20, 2010): compressed public-record example

This example is a small **public-record abstraction** of the Macondo well
blowout on the Deepwater Horizon. It is not a forensic reconstruction. The
adversary is a physical hazard (pressured reservoir hydrocarbons), so facts use
`hazard:`, `well:`, `rig:` and `monitor:` prefixes rather than `attacker:`.

## Sources (all opened and read for this lab)

- **[BP]** BP, *Deepwater Horizon Accident Investigation Report*, 8 September
  2010 (the "Bly report"), as filed with the SEC as Exhibit 99.3 (report body,
  without appendices):
  https://www.sec.gov/Archives/edgar/data/313807/000119312510216268/dex993.htm
- **[NC]** National Commission on the BP Deepwater Horizon Oil Spill and
  Offshore Drilling, *Deep Water: The Gulf Oil Disaster and the Future of
  Offshore Drilling*, Report to the President, January 2011 (chapter 4):
  https://www.govinfo.gov/content/pkg/GPO-OILCOMMISSION/pdf/GPO-OILCOMMISSION.pdf
- **[CC]** National Commission, *Macondo: The Gulf Oil Disaster, Chief
  Counsel's Report*, 2011 (chapters 4.6 negative test, 4.7 kick detection, 4.8
  kick response, 4.9 BOP):
  https://www.govinfo.gov/content/pkg/GOVPUB-PR-PURL-gpo4390/pdf/GOVPUB-PR-PURL-gpo4390.pdf
- **[CSB]** U.S. Chemical Safety and Hazard Investigation Board, *Investigation
  Report Volume 2: Explosion and Fire at the Macondo Well*, Report
  No. 2010-10-I-OS, 5 June 2014 (BOP failure analysis, barrier management):
  https://www.csb.gov/assets/1/7/vol_2_final_version.pdf

The CSB published further Macondo volumes later. I did not read them, and
nothing here relies on them. The DNV forensic BOP examination is cited only
through what CSB says about it ("previous incident investigation reports").

## BP's eight key findings, checked

The assignment's list of the eight findings is essentially right. Two points
needed correcting:

| # | BP key finding (wording from the report) | How it appears here |
|---|---|---|
| 1 | The annulus cement barrier did not isolate the hydrocarbons. | Scenario event `annulus-cement-fails`. Not offered as a choice: a cement Picket would end the trace at step one and teach nothing. |
| 2 | The shoe track barriers did not isolate the hydrocarbons. | Event `shoe-track-fails`; **Picket P2**. |
| 3 | The negative-pressure test was accepted although well integrity had not been established. | Event `negative-test-shows-flow`; **Vedette chain V1 → V2**, then **Reserve R1**. |
| 4 | Influx was not recognized until hydrocarbons were in the riser. BP: indications "discernable in real-time data from approximately 40 minutes before the rig crew took action". | Events `influx-begins`, `returns-sent-overboard`, `influx-grows`; **Vedettes V3, V4**, then **Picket P1**. |
| 5 | Well control response actions failed to regain control of the well. | **This finding already contains the MGS-not-overboard choice:** "routing the fluids exiting the riser to the ... mud gas separator (MGS) system rather than to the overboard diverter line". Annular closed ~21:41 without sealing; a VBR sealed ~21:47. **Picket P3** (pipe ram), **Reserve R2** (overboard). |
| 6 | Diversion to the mud gas separator resulted in gas venting onto the rig. | This is the *consequence* (the goosenecked 12 in. vent sent gas down onto the rig), not a second diversion decision. Event `mgs-vents-gas-onto-rig`. |
| 7 | The fire and gas system did not prevent hydrocarbon ignition. | Event `gas-ignites`. Not a choice (see Omissions). |
| 8 | The BOP emergency mode did not seal the well. | BP: EDS "very likely disabled" by the explosions; AMF likely prevented by a faulty yellow-pod solenoid and low blue-pod batteries; ROV autoshear "likely" closed the BSR 33 hours later, and it did not seal. **CSB (2014) disagrees on the AMF**: it likely *did* fire on April 20, but the shear ram closed on buckled, off-centre drill pipe. Modelled as **Reserve R3** following CSB. |

Corrections, in short:

1. KF5 and KF6 are not "BOP actions" and "diversion" as separate barriers.
   KF5 is the whole crew response, including the MGS routing. KF6 is what the
   MGS routing then did: it put gas on the rig.
2. KF8 is where the investigations disagree most. BP says the AMF probably did
   not work. CSB says it probably worked, but the pipe was already buckled.

## Watchline roles

- **Vedettes**: the negative-test readings (V1 observe+review, V2
  critical+escalate) and kick monitoring (V3 flow-out/pits, V4 drill-pipe
  pressure+review). In both cases the record shows the chain **reached
  review and stopped before "critical"**:
  - Negative test: the readings were discussed at length. The 1,400 psi
    drill-pipe pressure was explained away as a "bladder effect" [NC]. The test
    was declared good at 19:55 [BP]. CC: the readings "were not ambiguous".
  - Kick: at about 21:30 the driller and toolpusher saw and discussed an odd
    drill-pipe/kill-line pressure difference. "Neither ... performed a visual
    flow check or shut in the well" [NC].
- **Pickets**: shoe track (P2), shutting in the BOP before gas passes it (P1;
  BP 5C says any time before about 21:38), and the pipe ram, which did seal
  the annulus at about 21:47 (P3).
- **Reserves**: keep the well overbalanced and repair it (R1), divert
  overboard (R2), and the deadman/AMF blind shear ram after power is lost (R3).
  The emergency disconnect (EDS) is also a Reserve. It is left out as a choice
  because both BP and NC say the explosion probably disabled it, so it would
  always be suppressed (BP: the emergency methods "were not fully
  independent").

## Documented cases where one barrier action undermined another

I looked for these specifically. Only mechanisms the sources state are
modelled.

1. **Pipe ram seal → buckled pipe → shear ram cannot seal** (modelled). CSB
   KF3: "Large pressure differences were established between the inside and
   outside of the drillpipe when well control actions by the crew sealed the
   well ... This likely caused drillpipe in the BOP to buckle". CSB 3.2.2: the
   partial BSR closure "punctured the drillpipe and caused flow from the well
   to the environment to reestablish ... minutes after closure of the pipe
   rams had actually sealed around the drillpipe." CSB 5.4 also says the BP
   and Transocean manuals' advice to switch from annular to pipe ram may
   "encourage the likelihood of pipe buckling". The scenario adds
   `well:pipe-off-center` after a pipe-ram seal, and that fact suppresses R3.
2. **Sheen-test overboard lineup blinded the flow-out meter** (modelled). For
   the sheen test, returns went straight overboard, "bypassing the pits,
   Sperry-Sun flow-out meter, and both gas meters" [CC 4.7; BP]. The sheen test
   is an environmental check, not a defense against the blowout. So this is an
   operation in the trace that adds `monitor:flow-out-bypassed`, which
   suppresses V3. It is not a choice.
3. **MGS lineup, possibly set to avoid pollution, put gas on the rig**
   (modelled as initial state). CC 4.8: the valves had been set to the MGS
   beforehand. "The crew may have done this to avoid inadvertently discharging
   oil-based drilling mud ... in violation of environmental regulations." CC
   also says diverting to the MGS "stopped the flow of mud onto the rig floor
   within seconds", so it worked briefly, and then "made ignition all but
   inevitable".
4. **Spacer chosen to avoid hazardous-waste disposal may have confounded the
   negative test** (perturbation variant only). CC 4.6: BP used lost-circulation
   material as spacer "to avoid having to dispose of the material as hazardous
   waste". It "may have migrated into and clogged the open kill line". CC also
   says the failed test "should have been clear" anyway.
5. **Engine-room ventilation did not trip automatically on gas, to protect
   station-keeping power** (not modelled). BP KF7: this design "was probably
   selected so that false gas-detection trips would not interrupt the power
   supply to the thrusters".

Mechanisms 3–5 are trade-offs between the blowout and *another* hazard
(pollution, waste handling, keeping the rig on station). One hazard model
cannot score those trade-offs; see Caveats.

## What the scenario preserves

15 events, in [BP]'s chronology order:

1. The annulus cement fails, then the shoe track fails, so the well is open to
   the reservoir.
2. The negative-pressure test shows the leak and is accepted.
3. Riser displacement leaves the well underbalanced (~20:52). The influx
   starts (39 bbl gain 20:58–21:08; drill-pipe pressure rising from 21:01).
   Returns go overboard for the sheen test, bypassing the meter (21:08–21:14).
   The influx grows (~300 bbl by 21:31). Hydrocarbons enter the riser (~21:38).
4. Mud erupts onto the rig floor (~21:40). The MGS is overwhelmed and vents gas
   onto the rig (~21:44–21:47). Flow continues past the BOP unless a pipe ram
   seals (~21:47). After a seal, the drill pipe buckles [CSB].
5. Gas ignites (21:49). The deadman fires the shear ram. The well reopens after
   the rig is lost, or the blowout simply continues if nothing sealed.

## What it omits

Cement design and testing, centralizers, float-collar conversion, the
positive-pressure test, the annular preventer's partial seal and the riser
slip joint, crane operations and other activities during displacement, the
yellow/blue pod wiring details, EDS, ROV interventions, the 87 days of
response and containment, and every human and organisational factor beyond
"chain stopped at stage X". The fire-and-gas system (KF7) is not a choice.
BP says that once gas spread beyond the classified areas "multiple ignition
sources could have existed", so a ventilation-trip Picket on one ignition path
would overclaim.

## Variants

| id | kind | what changes | why |
|---|---|---|---|
| `hanger-seal-path` | failed-attempt | Hydrocarbons get in past the casing hanger seal instead of the shoe track. | BP considered this path and concluded the seal held. What if it hadn't? |
| `rams-overwhelmed` | failed-attempt | The pipe ram does not seal (event renamed, so P3 cannot block it). | In the record the VBR sealed at ~21:47. NC: "Flow rates at this point may have been too high for either the annular preventer or a variable bore ram to seal the well." |
| `residual-gas-path` | alternative-branch | Gas reaches the rig even with overboard diversion. | BP 5C (slip-joint seals, moon pool) and CC 4.8 (drill pipe; 2–4 knot wind could blow gas back) both give this caveat while still recommending overboard. |
| `trip-tank-masking` | alternative-branch | The 20:58–21:08 flow-out signal is hidden by trip-tank emptying. | BP: emptying the trip tank "could have masked this flow indication". |
| `spacer-below-bop` | perturbation | Adds the spacer leaking below the BOP before the test. | CC says the drill-pipe reading was unambiguous regardless, so nothing should change. |

The assignment suggested a "negative test read correctly" counterfactual and a
"kick detected earlier" alternative. In Muster those are defensive choices,
not different hazard traces, so they appear as presets (V1+V2+R1, V4+P1,
V3+P1).

## Run it

```bash
BIN=muster   # or: go run .
$BIN lab    --file examples/deepwater-horizon-2010/lab.yaml --preset "As recorded"
$BIN robust --file examples/deepwater-horizon-2010/lab.yaml
$BIN search --file examples/deepwater-horizon-2010/lab.yaml --tradeoffs
```

### As recorded

```text
lab: deepwater-horizon-2010  trace: baseline  choices: [V1,V4,P3,R3]
  applied  Cement around the casing fails to seal
  applied  Hydrocarbons get in through the shoe track
  applied  Negative-pressure test shows the well leaking 🔭npt-readings
  applied  Mud in the riser replaced with seawater
  applied  Hydrocarbons start flowing into the well 🔭kick-pressure-watch
  applied  Returns sent overboard, bypassing the flow meter
  applied  Influx grows unnoticed 🔭kick-pressure-watch
  applied  Hydrocarbons pass the blowout preventer
  applied  Mud and gas erupt onto the rig floor
  applied  Overwhelmed mud-gas separator vents gas onto the rig
  blocked  Flow continues past the blowout preventer 🛡️pipe-ram-seals
  applied  Drill pipe buckles off-centre in the preventer
  applied  Gas ignites; explosions and fire
  applied  Well reopens after the rig is lost
  skipped  Blowout continues unsealed
adverse: hydrocarbons reached the rig floor; gas cloud spread over the rig; gas ignited and the rig exploded; uncontrolled release from the well to the sea
```

The `--json` trace shows `deadman-closes-bsr` at `gas-ignites` with
disposition `suppressed`, reason `well:pipe-off-center`.

### `muster robust`

```text
# Robustness: Deepwater Horizon / Macondo blowout (2010)

Adverse outcomes reached by each Watchline, as recorded and under each alternative trace.
↑ = worse than as recorded, ↓ = better, ≠ = different.

| Watchline | As recorded | Leak past the casing hanger seal | Pipe ram cannot seal | Gas reaches the rig despite diversion | Early flow signal masked | Spacer left below the preventer |
| --- | --- | --- | --- | --- | --- | --- |
| No added defense | O1, O2, O3, O4 | O1, O2, O3, O4 | O1, O2, O3, O4 | O1, O2, O3, O4 | O1, O2, O3, O4 | O1, O2, O3, O4 |
| As recorded | O1, O2, O3, O4 | O1, O2, O3, O4 | O1, O2, O3, O4 | O1, O2, O3, O4 | O1, O2, O3, O4 | O1, O2, O3, O4 |
| Negative test read correctly | none | none | none | none | none | none |
| Kick caught on drill-pipe pressure | none | none | none | none | none | none |
| Kick caught on flow-out | none | none | none | none | O1, O2, O3, O4 ↑ | none |
| Late response, done better | O1 | O1 | O1, O4 ↑ | O1, O2, O3, O4 ↑ | O1 | O1 |
| Trust the shoe track | none | O1, O2, O3, O4 ↑ | none | none | none | none |

Outcomes:
- O1 = hydrocarbons reached the rig floor
- O2 = gas cloud spread over the rig
- O3 = gas ignited and the rig exploded
- O4 = uncontrolled release from the well to the sea
```

### `muster search` (baseline and every variant, with `--tradeoffs`)

```text
# Search: deepwater-horizon-2010 (baseline)

10 choices, 1024 Watchlines replayed.

Strictly worse after adding one defense: 0 minimal case(s)
Tradeoffs (fixes one outcome, opens another): 0
```

The result is the same for `hanger-seal-path`, `rams-overwhelmed`,
`residual-gas-path`, `trip-tank-masking` and `spacer-below-bop`: 0 worse, 0
trade-offs.

## Interpretation

- **No defense makes things worse in this model, and that is the honest
  result.** The CSB buckling mechanism really is a barrier action undermining
  a later barrier. It does not make adding the pipe ram worse, because the
  sources do not say the shear ram would have sealed *without* it. CSB KF1 says
  that in an active blowout, "well pressures and well flow can exceed the
  design capabilities of the BOP elements". So R3 is suppressed in both
  branches: by `well:pipe-off-center` after a pipe-ram seal, and by
  `hazard:flow-through-bop` without one. **R3 never changes an outcome in any
  trace.** That is the model's version of CSB's conclusion: by 21:49 the last
  Reserve had already been defeated by what came before it.
- **Only the early exits are robust.** Reading the negative test correctly
  (V1+V2+R1), or acting on drill-pipe pressure before 21:38 (V4+P1), gives no
  harm in all six traces. These are the stages where the record shows the
  Vedette chain stopping at "critical".
- **Which instrument you watch matters.** V3+P1 (flow-out) works only because
  the 20:58–21:08 gain comes before the sheen-test lineup blinds the meter.
  Under `trip-tank-masking` it fails completely. V4+P1 (drill-pipe pressure)
  does not depend on either.
- **The late response is fragile.** Adding overboard diversion to the recorded
  Watchline cuts harm to "hydrocarbons reached the rig floor" in the baseline,
  which matches BP and CC: it would "probably" have kept most gas off the rig.
  But it depends on two things that held or might have held: the pipe ram
  sealing (`rams-overwhelmed` → release to sea) and the gas actually going
  overboard (`residual-gas-path` → everything). The good baseline result also
  rests on a model inference: with no ignition, the 21:47 pipe-ram seal keeps
  holding. CSB describes the well as "essentially static" just before the
  explosion, but no report claims that counterfactual.
- **A well-placed Picket is placed for the path you expect.** Trusting the
  shoe track is perfect against the path BP found, and useless against the
  hanger-seal path BP ruled out. The negative test covers both paths.

## Caveats

- **Kind is `public-record`**: every event traces to a published account.
  The counterfactual controls are labelled as counterfactual in
  `controls.yaml`. Facts are a single-pass compression of minutes of
  continuous, interacting flow physics.
- **Investigations disagree.** BP (2010) and CSB (2014) disagree about whether
  the AMF fired and about why the pipe was off-centre. CSB lists other
  candidate causes (weight of pipe above the BOP after the explosion; flow
  drag) and judges effective compression most likely. The model follows CSB on
  the BOP and BP's chronology for timings.
- **No time, no probability.** "More time to respond", which is BP's and CC's
  actual claim for overboard diversion, cannot be expressed. Diversion
  is modelled as keeping the gas off the rig (their central estimate), and the
  `residual-gas-path` variant covers the case where it does not.
- **The early flow-out signal is treated as detectable** in the baseline,
  because BP and CC computed a gain that volume accounting would have shown.
  BP says the trip-tank emptying could have masked it, which is the
  `trip-tank-masking` variant.
- **Chain stages are bundled into choices.** For example, P1 contains review,
  critical, escalate and the BOP closure. Choice roles name the decisive
  action.
- **Cross-hazard trade-offs are out of scope.** Three documented
  interactions (MGS lineup vs. pollution, LCM spacer vs. hazardous-waste
  rules, ventilation trips vs. station-keeping power) trade the blowout
  against another harm. A single-hazard adverse set cannot score them.
- **Engine notes.** The engine has no OR in `requires`, so release is reached
  by two events (`well-reopens-after-rig-loss`, `blowout-continues`). A
  variant can switch off a control that held in the record only by renaming
  the event it attaches to (`hanger-seal-fails`, `flow-overwhelms-rams`).
- This is an educational model of a past accident. It is not engineering or
  well-control guidance.
