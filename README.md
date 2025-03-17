# GitHub Mailbot

Reusable GitHub Action for sending an email everytime someone pushes to a repo.

TODO:
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
