# Bangladesh Bank SWIFT heist, February 2016: compressed public-record example

This is a small **public-record abstraction**, not a forensic reconstruction.
It compresses the published record of the February 2016 theft from
Bangladesh Bank's account at the Federal Reserve Bank of New York into 17
transitions. It then asks which defenses would have kept the money from
leaving, and how much of the real save was luck.

## Sources

Primary or near-primary:

- *Bangladesh Bank v. Rizal Commercial Banking Corp.*, No. 19 Civ. 983 (LGS),
  Opinion & Order (S.D.N.Y., 20 Mar 2020). The background section restates the
  **complaint's allegations**, "accepted as true for the purpose of this motion
  only". It is the main source for the Philippine money trail.
  https://www.steptoe.com/a/web/4MFfJt5fYdF4STxtQDvnsm/3ZFvT9/19-civ-00983-20200320-opinion-dismissing-hacking-rico-case.pdf
- *United States v. Park Jin Hyok*, criminal complaint and FBI affidavit
  (C.D. Cal., filed 8 Jun 2018), paras. 144–166 (intrusion path, SWIFTLIVE,
  the Sri Lanka stop, recovery status).
  https://www.justice.gov/opa/press-release/file/1092091/download
- BAE Systems Threat Research, *Two bytes to $951m* (25 Apr 2016), on how the
  malware tampered with SWIFT Alliance Access.
  https://baesystemsai.blogspot.com/2016/04/two-bytes-to-951m.html

Reuters reporting (read through republished copies; reuters.com could not be fetched):

- *How the New York Fed fumbled over the Bangladesh Bank heist* (Jul 2016).
  [Business Insurance](https://www.businessinsurance.com/new-york-fed-bangladesh-bank-cyber-security-heist-fireeye-jupiter/),
  [Asia Times](https://asiatimes.com/2016/07/special-report-how-the-new-york-fed-fumbled-over-the-bangladesh-bank-heist/),
  [The Daily Star](https://www.thedailystar.net/backpage/how-ny-fed-fumbled-over-bb-fund-heist-1258042)
- Typo report, 10 Mar 2016, read through the
  [OCCRP summary](https://www.occrp.org/en/news/bangladesh-typo-triggers-alarm-prevents-billion-dollar-bank-heist)
- *Sri Lankan in Bangladesh cyber heist says she was set up by friend*
  (1 Apr 2016), [The Daily Star copy](https://www.thedailystar.net/business/banking/sri-lankan-bangladesh-cyber-heist-says-she-was-set-friend-1202785)
- SWIFT's mandatory Alliance update (27 Apr 2016),
  [Carrier Management copy](https://www.carriermanagement.com/news/2016/04/27/153709.htm)
- The Daily Star, *Hackers chose weekend to break into BB systems* (8 Mar 2016),
  https://www.thedailystar.net/frontpage/hackers-chose-weekend-break-bb-systems-788368

### Corrections to the commonly told version

Checking the sources changed several parts of the story as it is usually told:

- **"Jupiter" was in the address, not in an order to a sanctioned party.** It
  was part of the Philippine branch address (RCBC's Jupiter Street branch,
  Makati). It matched a sanctioned Iranian-linked tanker and shipping company,
  which had nothing to do with the heist. Reuters reports that **four of the
  five orders the Fed had already paid** contained "Jupiter". So the sanctions
  hit did not block those payments. It set off a *full manual review* that
  held the other 30 orders. A source told Reuters the save was a "total fluke".
- **The formatting problem came first, and it did not stop anything.** All 35
  first-batch messages were rejected for incorrect formatting. The attackers
  fixed them and resent 35. The SDNY complaint counts 70 orders in all. The
  five paid orders "were properly formatted, SWIFT authenticated and went
  through automatically."
- **The misspelling is reported two ways.** Reuters' March 2016 sources said
  "fandation". Its April 2016 report says "Fundation". **Who stopped the
  $20M is also disputed.** Reuters (March) credits a query from the routing
  bank, Deutsche Bank. Reuters (April) says the money reached Pan Asia Bank,
  which "refused to release the funds as the amount was unusually large". The
  FBI affidavit says it "was stopped by the recipient bank". The model keeps
  both as Pickets in series.
- **The Fed did query Dhaka.** On 4 Feb it flagged 12 orders, and it sent more
  SWIFT queries on 5 Feb. Those messages never reached anyone at Bangladesh
  Bank, and the Fed "did not reach out to Dhaka in any other way". On the other
  side, Bangladesh Bank's Saturday emails and fax went to Fed contacts that
  were monitored only on weekdays.
- **The printer tampering left a visible gap.** A Bangladesh Bank official
  found the confirmations missing at about 10:30 a.m. on Friday 5 Feb. He took
  it for a routine printer fault, and staff left for the Friday holy day.
- **The Philippine money was mostly still in place when the stop-payment
  arrived.** This is per the lawsuit's allegations. On 5 Feb a first tranche
  left the fake accounts. RCBC put a hold on the remitter accounts that
  evening and lifted it 45 minutes later. RCBC was logged off SWIFT from Sunday
  7 Feb to Tuesday 9 Feb at 9:11 a.m., when it received the STOPPAY. Between
  10:24 and 11:37 a.m. that day, $42.9M and $15.2M were moved out. RCBC froze
  the fake accounts at 3:31 p.m., by which time they were nearly empty.
- **Recovery:** the FBI affidavit (June 2018) says "the majority of the
  $81,000,000 has not been recovered". I did not verify later recovery figures
  or the current status of the New York litigation against a primary source.

## What the scenario preserves

1. Network foothold, then the SWIFTLIVE server, then malware that tampers with
   the Alliance Access records and printouts (DOJ, BAE).
2. Orders sent Thursday night Dhaka time: format rejection, then resubmission.
3. The Fed's automatic payment of five orders ($81M to four RCBC accounts via US
   correspondents, $20M to Sri Lanka), while the Fed's SWIFT queries go unread.
4. The two Sri Lankan barriers, in series.
5. The 5 Feb first tranche, which left *before* Dhaka staff were at work. Manila
   late morning is about 9–10 a.m. in Dhaka.
6. The visible printer gap at 10:30 a.m., not followed up because of the weekend.
7. The "Jupiter" hit, then the manual review, then the remaining orders held.
8. The recorded stop-payment reaching RCBC (Tuesday 9:11 a.m.), the bulk drain
   regardless, and dispersal through a remitter into casinos.

Two events are **counterfactual branches**: `fed-pays-remaining-orders` and
`sri-lanka-payout`. Defenses that held kept them closed in 2016. They sit in
every trace, and in the baseline they are skipped.

Omitted: the 2015 spear-phishing detail, the malware internals, exact amounts
per order, the US correspondent hops, the account-level laundering (Go-Centurytex,
Abba, Beacon), the casino junkets, Kim Wong's partial repayment, and every
attribution question. Actor attribution does not affect this replay.

## Watchline choices

| Key | Role | Defense | Status in 2016 |
|---|---|---|---|
| V1 | Vedette | printed SWIFT confirmations | existed. `suppressed_by` the print tampering |
| V2 | Vedette | reconcile Fed debits against a record kept off the SWIFT server | counterfactual |
| V3 | Vedette | Fed phones Dhaka about its queries | counterfactual. Reuters: no other channel was used |
| P1 | Picket | Fed real-time anomaly hold | counterfactual. Reuters: review is "usually the day after" |
| P2 | Picket | RCBC freezes on receipt of stop-payment (`requires: stoppay:at-rcbc`) | counterfactual |
| P3 | Picket | RCBC keeps its 5 Feb hold | the hold existed and was lifted after 45 minutes (alleged) |
| R1 | Reserve | business-hours follow-up and stop-payment | existed. `suppressed_by` the Fri–Sat weekend |
| R2 | Reserve | 24/7 duty officer and stop-payment | counterfactual |

## Variants

| id | kind | change |
|---|---|---|
| `no-jupiter-hit` | failed-attempt | No sanctions-like name, so no manual review, and the Fed pays the rest |
| `no-misspelling` | failed-attempt | No routing-bank query. The recipient bank's check remains |
| `sri-lanka-checks-fail` | failed-attempt | Both Sri Lankan barriers fail |
| `weekday-dhaka` | alternative-branch | No Bangladesh weekend: the printer gap is chased and the recorded stop-payment goes the same day |
| `no-ph-holiday` | alternative-branch | Launderers drain the fake accounts before the recorded stop-payment arrives |
| `formatted-first-time` | perturbation | No format rejection or resubmission. Should not matter |

The failed-attempt variants differ from the baseline only in initial facts.
The alternative branches move one event each. The perturbation changes one
effect.

## Results (frozen `muster-bin`, rebuilt with the engine fix in which waiting or suppressed Pickets do not block)

```
$ muster robust --file examples/bangladesh-bank-2016/lab.yaml
# Robustness: Bangladesh Bank SWIFT heist (2016)

Adverse outcomes reached by each Watchline, as recorded and under each alternative trace.
↑ = worse than as recorded, ↓ = better, ≠ = different.

| Watchline | As recorded | No "Jupiter" sanctions hit | Sri Lankan payee spelled correctly | Both Sri Lankan checks fail | Sent on a Bangladesh working night | No Philippine long weekend | Orders formatted correctly first time |
| --- | --- | --- | --- | --- | --- | --- | --- |
| No added defense | O3, O4, O5 | O1, O3, O4, O5 ↑ | O3, O4, O5 | O2, O3, O4, O5 ↑ | O3, O4, O5 | O3, O4, O5 | O3, O4, O5 |
| As in 2016 | O3, O4, O5 | O1, O3, O4, O5 ↑ | O3, O4, O5 | O2, O3, O4, O5 ↑ | O3, O4, O5 | O3, O4, O5 | O3, O4, O5 |
| RCBC acts on the stop-payment | O3, O5 | O1, O3, O5 ↑ | O3, O5 | O2, O3, O5 ↑ | O3, O5 | O3, O4, O5 ↑ | O3, O5 |
| Weekend duty officer chasing the printer | O3, O5 | O1, O3, O5 ↑ | O3, O5 | O2, O3, O5 ↑ | O3, O5 | O3, O5 | O3, O5 |
| Fed real-time hold | none | none | none | none | none | none | none |
| RCBC keeps its own hold | O3, O5 | O1, O3, O5 ↑ | O3, O5 | O2, O3, O5 ↑ | O3, O5 | O3, O5 | O3, O5 |
| Independent reconciliation, duty officer, RCBC freeze | none | O1 ↑ | none | O2 ↑ | none | none | none |

Outcomes:
- O1 = Fed pays the ~30 remaining orders (~$850M)
- O2 = $20M released to the Sri Lankan payee
- O3 = first Philippine tranche leaves the fake accounts (5 Feb)
- O4 = rest of the $81M drained from the fake accounts (9 Feb on)
- O5 = money laundered through casinos, trail goes cold
```

```
$ muster search --file examples/bangladesh-bank-2016/lab.yaml --tradeoffs
# Search: bangladesh-bank-2016 (baseline)

8 choices, 256 Watchlines replayed.

Strictly worse after adding one defense: 0 minimal case(s)
Tradeoffs (fixes one outcome, opens another): 0
```

`search --variant <id> --tradeoffs` also reports **0 worse and 0 tradeoff**
cases for every one of the six variants.

An extra brute-force pass replayed every Watchline of up to three choices
against all seven traces. Every Watchline that ends with no harm in every
trace contains **P1**. The only Watchlines that are clean in the baseline
without P1 are `V2,P2,R2` and `V3,P2,R2`. Both fail in the same two
failed-attempt variants: O1 if the Jupiter hit does not fire, and O2 if both
Sri Lankan checks fail.

## Interpretation

- **The 2016 Watchline did nothing.** The "As in 2016" row matches "No added
  defense" in every trace. The printed confirmations were suppressed by the
  malware. The business-hours follow-up was suppressed by the weekend.
  Everything that held in 2016 sits in the trace as an existing barrier. None
  of it was a control Bangladesh Bank operated.
- **The largest save was luck.** Without the Jupiter hit, every Watchline that
  lacks a Fed-side Picket also loses the ~$850M (the O1 column). Downstream
  defenses at Bangladesh Bank and RCBC cannot touch money that the Fed has
  already paid.
- **The typo was not load-bearing in this model.** Remove the misspelling and
  nothing changes, because the recipient bank's own check on the unusually
  large credit still holds. That check is what the FBI affidavit credits. The
  $20M is lost only when both barriers fail.
- **The weekend made no difference to any outcome here, and that is a
  finding, not a null.** On the `weekday-dhaka` branch, R1 does act. Trace:
  `bb-office-chases-print-failure` is `ready`, not `suppressed`. But the
  first tranche had already left before Dhaka staff arrived, and the recorded
  stop-payment already reached RCBC *before* the bulk drain. What decides O4
  is whether RCBC acts on the stop-payment (P2), not how fast Dhaka sends it.
- **The Philippine holiday cut both ways.** RCBC's stop-payment freeze (P2)
  works only because the long weekend held the launderers until after the
  stop-payment arrived. In `no-ph-holiday`, P2 alone loses O4. RCBC keeping
  its own hold (P3) does not depend on that timing, and neither does an early
  stop-payment from a duty officer (R2).
- **Only detection at payment time beats the first tranche.** V1 with R2 (a
  duty officer chasing the printer gap) is too late for 5 Feb. V2 (independent
  reconciliation) or V3 (the Fed phoning) with R2 and P2 gets a stop-payment
  to RCBC before the first tranche, and clears all Philippine harm.
- **No non-monotone interaction.** Search finds no case where adding a defense
  makes things worse. This wiring has no Reserve that removes an attacker
  capability that later evidence depends on, so none was expected. None was
  engineered.

## Caveats

- **The Philippine steps rest on allegations.** They come from the complaint,
  as a US court restated it on a motion to dismiss. The court made no
  findings of fact, and RCBC and the casinos dispute liability.
- **The `no-jupiter-hit` branch takes the Reuters source's "total fluke"
  at face value.** The record does not explain why the 30 held orders were not
  among the five auto-cleared. Fed staff had already flagged 12 orders, and
  that alone might have stopped some of them.
- **The two alternative branches are plausible orderings, not documented
  ones.** In particular, `no-ph-holiday` assumes the launderers would have
  moved as soon as the bank reopened.
- **R1 and R2 assume a stop-payment channel that does not depend on the broken
  SWIFT server.** In 2016 Bangladesh Bank's SWIFT terminal would not start on
  Saturday, and its email and fax went to weekday-only Fed contacts.
- **A known trace gives the matched Pickets an advantage defenders would not
  have had in advance.** P1, P2 and P3 are written with hindsight about
  exactly which transitions mattered.
- **Amounts are labels, not model quantities.** The engine has no amounts, so
  "first tranche" and "bulk" are coarse and only approximate the sums in the
  complaint.
- **Engine limits hit:** there is no OR in preconditions, so the two Sri Lankan
  barriers and the two R1/R2 triggers are separate events and controls. There
  is no clock, so timing lives entirely in event order, and weekday versus
  weekend has to be expressed as a suppression fact plus reordering.
- **Nothing here is operational.** The malware description is limited to what
  BAE and the FBI published.
