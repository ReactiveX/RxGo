# Contributions Guidelines
Contributions are always welcome. In fact, yours can shape the project this young. However, to make this a smooth collaboration experience for everyone and to maintain the quality of the code, here is a few things to consider before and after making a pull request:

## Is it a feature or a patch?
Before contributing, ask yourself if the change you're working on is a feature or a patch. A feature (like a new operator, function, etc.) needs a unit test with it, while a patch (like fixing a typo in `README.md`) may not. For small typos, it is always a good idea to just [email me](sendto:jo.chasinga@gmailcom) about it so we don't clutter the branch with super small changes.

## Include a unit test
Any change in the code other than small patches won't be accepted without a unit test. We use [testify/assert](https://github.com/stretchr/testify#assert-package) for the test.

## Read these great tips

https://github.com/erlang/otp/wiki/Writing-good-commit-messages    
https://help.github.com/articles/closing-issues-via-commit-messages    
https://github.com/blog/1506-closing-issues-via-pull-requests    

There is no rigid rule or anything, just remember that a brief, succinct commit summary line can go a long way. In general, please don't include issue number the commit intends to resolve (see the above link for tip on how to do that), alphanumeric or any symbols, or past tense. The commit summary line should reads like a short directive like:

> Add feature A to class B

If you would like to be more elaborate, please see the first link on how to add that to the commit description.

## Open an issue
This is to encourage discussions and reach a sound design decision before implementing an additional feature or fixing a bug. If you're proposing a new feature, make it obvious in the subject.

## Create a feature branch
Always create a branch dedicated to the feature or fix you intend to be merged into the base repository and keep it clean and commits minimal. 

## Squash your commits
It is asked that you [squash your commits](https://github.com/blog/2141-squash-your-commits) down to a few, or one, discreet change sets before submitting a pull request. Fixing a bug will usually only need one commit, while a larger feature might contain a couple of separate improvements that is easier to track through different commits.

## Include some comments
It is not mandatory but rather nice to have a short comment accompanying a pull request. We're all humans here and sometime we lose tracks of things.

## Write nice code
Try to write idiomatic code according to [Go style guide](https://github.com/golang/go/wiki/CodeReviewComments). Also, see this project style guide for project-specific idioms (when in doubt, consult the first).
