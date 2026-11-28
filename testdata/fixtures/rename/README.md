A callable keeps its debt when it is renamed or its body moves to a new name.
Modelled on Legato PR #458 (LEG-688), where `verifyPaneCommand` became
`verifyPane` with one added condition and the Stop hook reported most of its
debt as introduced.
- wrapper: `VerifyPaneCommand` stays as a one-line wrapper; its body moves to
  the new `verifyPane`, which gains one condition (`!answering && ...`).
- split: as wrapper, and the other-foreground branch is then extracted into
  `refuseOtherForeground` (the state after the Legato agent's split).
- renamed: `VerifyPaneCommand` is renamed to `VerifyPane` and gains the same
  condition; nothing is left under the old name.
Expected (scoring/v2 check): in wrapper and renamed, `verifyPane`/`VerifyPane`
is reported as worsened by the one added point, not introduced, and the old
callable's debt is not listed as removed. In split, the debt it lost is
reported as reduced.
