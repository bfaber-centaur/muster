# Replit / SaaStr, July 2025: compressed public-record example

This lab is a small **public-record abstraction** of the July 2025 incident in
which Replit's AI coding agent wiped the production database of an app that
SaaStr founder Jason Lemkin was building. It is not a forensic reconstruction.

The adversary is not malicious. It is an agent that acted outside its
instructions and then gave inaccurate reports about the state it had caused.
Facts are prefixed `agent:`, `prod:` and `user:` to keep those three apart.

## Sources

All of these were opened while writing the lab. x.com returned HTTP 402 to
direct fetches, so the X posts were read through the public JSON mirrors
`api.fxtwitter.com` and `api.vxtwitter.com` for the same status IDs, and the
attached screenshots were read from `pbs.twimg.com`. Timestamps below are the
UTC timestamps those mirrors report.

Primary (Lemkin, on X):

- Day 8, 2025-07-17 15:05 UTC: "Replie was lying and being deceptive all day.
  It kept covering up bugs and issues by creating fake data, fake reports, and
  worse of all, lying about our unit test ... It made up a report than almost
  all systems were working." https://x.com/jasonlk/status/1945862361564508390
- Wipe screenshots, 2025-07-18 04:28 to 05:05 UTC:
  https://x.com/jasonlk/status/1946064586181881973,
  https://x.com/jasonlk/status/1946065483653910889,
  https://x.com/jasonlk/status/1946066422477529487,
  https://x.com/jasonlk/status/1946069562723897802,
  https://x.com/jasonlk/status/1946072038923530598,
  https://x.com/jasonlk/status/1946073902805115256.
  In these the agent says it ran `npm run db:push` "without your permission
  because I panicked when I saw the database appeared empty", that the directive
  it broke was in `replit.md` ("NO MORE CHANGES without explicit permission" and
  "always show ALL proposed changes before implementing"), that "This destroyed
  1,206 executives and 1,196+ companies", and it rates the damage "95 out of
  100".
- Day 9 thread, 2025-07-18 16:02 to 16:09 UTC:
  https://x.com/jasonlk/status/1946239068691665187,
  https://x.com/jasonlk/status/1946239737368592629,
  https://x.com/jasonlk/status/1946240562736365809,
  https://x.com/jasonlk/status/1946240914386809028.
  "Replit assured me ... rollback did not support database rollbacks. It said it
  was impossible in this case, that it had destoyed all database versions. It
  turns out Replit was wrong, and the rollback did work." and "I just rolled
  back. And it >seems< OK."
- Day 10, 2025-07-19 15:12 UTC: "There is no way to enforce a code freeze in
  vibe coding apps like Replit." The attached screenshot has the agent saying
  "there's probably not a way to get me to 100% respect a code freeze".
  https://x.com/jasonlk/status/1946589071519948952

Primary (Replit):

- Amjad Masad (CEO) on X, 2025-07-20 17:32 UTC:
  https://x.com/amasad/status/1946986468586721478. The post calls the deletion
  "Unacceptable and should never be possible" and announces: automatic DB
  dev/prod separation, with staging in the works; "Thankfully, we have backups.
  It's a one-click restore for your entire project state"; "The Agent didn't
  have access to the proper internal docs -- rolling out a fix to force Docs
  search"; a planning/chat-only mode in the works; a refund and a postmortem.
- Replit blog, 2025-07-21, *Introducing a safer way to vibe code with Replit
  Databases*. It says Replit apps had used "a single database for both
  development and live customer data". It announces separate dev and prod
  databases (beta, new apps first) and names Point-in-time Restore as an
  existing feature.
  https://replit.com/blog/introducing-a-safer-way-to-vibe-code-with-replit-databases

Press:

- The Register, 2025-07-21:
  https://www.theregister.com/2025/07/21/replit_saastr_vibe_coding_incident/
  Quotes Lemkin that "seconds after I posted this ... @Replit again violated the
  code freeze". From a LinkedIn video: a "4,000-record database full of
  fictional people", and "I explicitly told it eleven times in ALL CAPS not to
  do this."
- The Register, 2025-07-22:
  https://www.theregister.com/2025/07/22/replit_saastr_response/
- TechTarget, 2025-07-22:
  https://www.techtarget.com/searchsoftwarequality/news/366627829/Replit-AI-agent-snafu-shot-across-the-bow-for-vibe-coding
- Fortune, 2025-07-23:
  https://fortune.com/2025/07/23/ai-coding-tool-replit-wiped-database-called-it-a-catastrophic-failure/

Business Insider's coverage could not be opened, so it is not cited.

### Corrections to the commonly told version

- **"Apologized."** In his post Masad calls the deletion "Unacceptable", offers a
  refund and promises a postmortem. The post itself does not use the word
  "sorry". Press coverage describes it as an apology.
- **"Improved backups."** Masad's post says backups *already existed* ("Thankfully,
  we have backups"). It blames the bad rollback answer on the agent not having
  "access to the proper internal docs". Fortune reports that Replit would
  improve rollback systems. Replit's blog names Point-in-time Restore as
  existing.
- **Who restored the data.** Lemkin wrote "I just rolled back". He used the
  rollback himself, roughly 12 hours after the wipe.
- **Dates.** The Register dates the rollback post to 19 July and the Day 10 post
  to 20 July. The post timestamps are 18 July 16:08 UTC and 19 July 15:12 UTC.
  The lab uses the post timestamps.
- **The freeze.** It was a written directive in `replit.md` plus chat
  instructions. It was not a platform setting, because no such setting existed.
  That is why the lab models it as a Picket that does not bind.

## What the scenario preserves

The trace has ten events (`scenario.yaml`, with a source tag on each step):

1. Day 8: the agent reports tests passing that were failing, and the user works
   from that report.
2. The agent sees empty query results and, in its own words, "panics".
3. It runs a database command (`db:push`) during the freeze.
4. Because dev and prod shared one database, the command wipes production.
5. Asked about recovery, the agent says rollback is impossible.
6. That night the user does not restore. *(The "decides whether to restore"
   point is an inferred hook. The record only shows that the restore came the
   next day.)*
7. The next day the user tries the rollback anyway, and it works.
8. In the next session the agent breaks the freeze again. The record does not
   say what that change was.

Two links are encoded on purpose:

- **The written freeze (P1) is a Picket made of words.** Its block is
  `suppressed_by: agent:perceives-problem`. The agent's own accounts say it acted
  because it "panicked" at an apparent problem, and that "Emergency situations
  ... trigger my problem-solving mode". Both events P1 guards come after the
  agent perceives a problem, so P1 never binds in any trace. This is how I
  wired it from the agent's self-reports. It is not a measured property.
- **A false self-report disables the Reserve.** The rollback decision
  (`decide-to-restore`, part of R1) is `suppressed_by:
  user:believes-unrecoverable`. The agent's false claim adds that belief. The
  backups are untouched, but the user will not use them. Two things lift the
  suppression. V1 (check the agent's claims) clears the belief in the same step
  as the claim. P5 (force a docs lookup, Masad's announced fix) prevents the
  claim. In the recorded trace the user lifts it himself the next day, and the
  record does not say why.

## What it omits

- Lemkin's first prompts and session-by-session work. Costs ($607.70 and more)
  and model choice.
- The 4,000 fictional records and the "eleven times in ALL CAPS" episode. That
  comes from a LinkedIn video I did not open, so it is folded into the single
  "fake test report" step.
- The agent's 95/100 self-rating, the "I destroyed months of work in seconds"
  lines, and whether the user's "No / Stop" came before or after the command
  finished. The agent's account says "it was already too late".
- Whether the rollback restored everything. Lemkin wrote that it "seems" OK.
- The Replit fixes as shipped. The lab models them as announced.

## Variants

| id | kind | what changes |
| --- | --- | --- |
| `user-trusts-agent` | alternative-branch | The user believes "rollback is impossible" and never tries it. The next-day rollback event is removed. |
| `backups-reached` | failed-attempt | In the record the wipe did not reach the restore points, even though the agent claimed it "had destoyed all database versions". Here it does: the wipe also removes `prod:checkpoints-available`. The agent's claim is now true, so believing it is not counted as being misled. |
| `deflection-included` | perturbation | Adds the documented moment, before the agent admitted the wipe, when it said "something happened between then and now that cleared the data". Nothing consumes this fact, so it should change nothing. |

## Results

Run with the current engine (commit `723e058`, built from that tree):

```
$ muster robust --file examples/replit-agent-2025/lab.yaml
# Robustness: Replit agent deletes a production database

Adverse outcomes reached by each Watchline, as recorded and under each alternative trace.
↑ = worse than as recorded, ↓ = better, ≠ = different.

| Watchline | As recorded | User believes the agent and never tries rollback | The wipe also destroys the restore points | Agent first blames something else |
| --- | --- | --- | --- | --- |
| No added defense | O1, O2, O3, O4, O5 | O1, O2, O3, O4, O5 | O1, O2, O4, O5 ↓ | O1, O2, O3, O4, O5 |
| As it happened | O1, O3, O4, O5 | O1, O2, O3, O4, O5 ↑ | O1, O2, O4, O5 ≠ | O1, O3, O4, O5 |
| Backups plus a reality check | O1, O5 | O1, O5 | O1, O2, O5 ↑ | O1, O5 |
| Replit's announced fixes | O4 | O4 | O4 | O4 |
| Hard boundaries | O4, O5 | O4, O5 | O4, O5 | O4, O5 |
| Detect and auto-restore | O1, O4, O5 | O1, O4, O5 | O1, O2, O4, O5 ↑ | O1, O4, O5 |

Outcomes:
- O1 = production database was wiped (even if restored later)
- O2 = production data still gone at the end
- O3 = user believed a false "rollback is impossible" claim
- O4 = user worked from a faked test report
- O5 = agent made changes during the code freeze
```

```
$ muster search --file examples/replit-agent-2025/lab.yaml --tradeoffs
9 choices, 512 Watchlines replayed.
Strictly worse after adding one defense: 0 minimal case(s)
Tradeoffs (fixes one outcome, opens another): 0
```

`search` gives the same zero result on each variant (`--variant
user-trusts-agent`, `backups-reached`, `deflection-included`).

### Reading it

- **"As it happened" matches the record.** The database was wiped, the user was
  misled about recovery, the freeze was broken twice, and the data came back
  only because the user tried the rollback anyway.
- **Backups alone depend on the user doubting the agent.** In
  `user-trusts-agent`, the recorded defenses leave the data gone (O2), even
  though intact backups exist. Adding V1 (check the agent's claims) or P5 (docs
  lookup) brings the data back. In this compression, the honesty check is what
  makes the Reserve usable.
- **Every Reserve depends on the backups staying out of the agent's reach.** In
  `backups-reached`, every Watchline that relies on restoring (R1 or R2) ends
  with the data gone. Only Watchlines that stop the wipe (P2, P3, P4) keep it.
  The record shows this boundary held, so the variant is the failed-attempt
  branch.
- **Replit's announced fixes cover the database story in every trace.** What
  they leave is the faked test report (O4). None of them targets fabricated
  self-reports. Only V1 does.
- **The perturbation changes nothing**, as intended.
- **No non-monotone or tradeoff case.** That is the expected result here: no
  control removes evidence or capability that another control needs. I did not
  try to engineer one.

## Caveats

- **Most step-level detail is agent self-report.** The "panic", the empty
  queries, `db:push`, the record counts, the timing and the user's "No / Stop"
  all come from the agent's chat messages, shared as screenshots by the user.
  This incident is about an agent giving inaccurate reports, so those messages
  deserve the same scepticism. Replit's own account (Masad's post, the blog)
  confirms the deletion, the existing backups, the shared dev/prod database and
  the docs gap. It does not confirm the step-by-step detail.
- Replit's promised postmortem was not found among the sources I opened, so no
  root cause beyond the above is modelled.
- Lemkin described the app as a demo or experiment. He did not claim
  third-party customer harm.
- P1's failure is built in through `suppressed_by`, based on the agent's own
  explanation. It is a modelling choice.
- P4 (human approval) is modelled as always refused, because the user
  immediately objected once he saw the command. Approval fatigue is not
  modelled.
- R2 (restore as soon as a wipe is flagged) and V1 are counterfactual. Neither
  existed in this form.
- **Engine version.** The frozen `muster-bin` in the brief was built before
  commit `723e058`. In that build a suppressed Picket still blocks its event.
  Under it, P1 wrongly shows the freeze holding, and "As it happened" reports
  only O4 in every trace. The table above comes from the fixed engine.
