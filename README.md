# GitHub Mailbot

[Reusable GitHub Action](https://docs.github.com/en/actions/how-tos/reuse-automations/reuse-workflows)
for sending a commit email with the git diff everytime someone pushes to a branch.
GitHub provides commit emails for PRs only, but those don't include the diff.

## Example workflow file

```
name: Commit mailbot
on: push

jobs:
  mailbot:
    if: github.actor != 'dependabot[bot]'
    uses: sanjit-bhat/github-mailbot/.github/workflows/mailbot.yml@main
    with:
      host: smtp.gmail.com
      port: 587
      from: alice@gmail.com
      to: bob@gmail.com,charlie@gmail.com
    secrets:
      password: ${{ secrets.MAILBOT_PASSWORD }}
```

This assumes a
[GitHub secret](https://docs.github.com/en/actions/how-tos/write-workflows/choose-what-workflows-do/use-secrets)
called `MAILBOT_PASSWORD` with the SMTP password for `alice@gmail.com`.
It doesn't run on commits from `dependabot[bot]` since those
don't have secret access, resulting in a mailbot error.

## TODO

1. Find out how to get rid of "---" in `git show` output.
1. Document how to use this.
1. Add email screenshots.
1. Document code.
1. Add tests for the parts of the pipeline that are testable.
1. Explain why we used reusable workflow vs. composable workflow.
1. Cache at least specific versions of the brew pkgs.
1. Delta removes color in the `git show` header.
Maybe we can only pipe the diff part to delta.
1. Set email time to the commit time.
1. Document testing process.
E2E testing on forked mailbot can use `mailbot_repo` flag.
1. Document looking at GH action output log to see GH env json, which is
passed to Go mailbot script and contains information that can be useful
for debugging mailbot behavior and adding features.
1. Don't run the mailbot workflow unless it can access the secret.
E.g., dependabot PRs can't access the secret, so either they need
to be given access, or the workflow should not run.
