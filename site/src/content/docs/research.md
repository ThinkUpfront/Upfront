---
title: Research
description: Empirical research backing Upfront's design decisions.
---

## The skill erosion problem

**Anthropic (2026):** Randomized controlled trial — 52 junior engineers. The AI group scored 17% lower on comprehension. The largest gap was in debugging. Developers who used AI for conceptual questions scored 65%+, while those who delegated code generation scored below 40%. The interaction pattern matters more than whether AI was used.

**METR (2025):** AI tooling increased completion time by 19% for experienced open-source developers on their own repositories. The productivity narrative assumes AI helps everyone — this study shows it can hurt experienced developers on familiar codebases.

**Harvard/BCG (2023):** 758 consultants using GPT-4 — quality jumped 40% inside the AI's capability frontier, but dropped 19 percentage points below the no-AI group on tasks outside it. Judgment atrophies without exercise.

**DORA (2025):** AI adoption increases throughput but also increases delivery instability. Speed without verification is not a gain.

**Microsoft Research (2024):** Developers using AI assistants experienced a "false sense of confidence" — they believed their code was more secure than it actually was.

## The review collapse

**SmartBear/Cisco:** Review effectiveness collapses above 400 lines of code. This is why `/upfront:plan` targets ~400 LOC phases.

**Faros AI (2024):** 10,000 developers across 1,255 teams — PR volume increased 98% after AI adoption, but net throughput showed zero improvement. Volume increase absorbed by review overhead and rework.

**GitClear (2024):** Across 211 million changed lines, code churn doubled while refactoring collapsed. More code, less of it survives production.

## Why forcing functions work

**Nagappan & Ball (2005):** Code churn predicts defects with 89% accuracy. Rework rate is one of the strongest signals of software quality.

**Capers Jones:** Defect removal efficiency above 95% is adequate quality, across 12,000+ projects.

**Mantyla & Lassenius (2009):** 75% of defects found in code review are evolvability issues, not functional bugs. This is why `/upfront:build`'s review checks architecture and spec compliance, not just correctness.

## The specification gap

**Montgomery et al.:** Ambiguous requirements are the single largest source of downstream defects. Clear intent — what `/upfront:feature` forces — reduces rework more than any other intervention.

**Jellyfish (2025):** 60% of engineering leaders cite "lack of clear metrics" as their biggest AI challenge. Only 20% measure actual impact on delivery outcomes.
