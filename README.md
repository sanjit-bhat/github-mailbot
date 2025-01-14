# GitHub Mailbot

Reusable GitHub Action for sending an email everytime someone pushes to a repo.

TODO:
1. Delta doesn't have Rocq file support. It thinks it's Verilog.
Perhaps we can extend it to work or use some other syntax highlighting differ.
Or just remove syntax highlighting, although it's very useful.
1. Find out how to get rid of "---" in `git show` output.
1. Document how to use this.
1. Add email screenshots.
1. Document code.
1. Add tests for the parts of the pipeline that are testable.
1. Explain why we used reusable workflow vs. composable workflow.
1. Cache at least specific versions of the brew pkgs.
1. Delta removes color in the `git show` header.
Maybe we can only pipe the diff part to delta.
