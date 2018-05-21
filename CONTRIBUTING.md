# Contributing

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

We love pull requests from everyone. By participating in this project, you
agree to abide by the thoughtbot [code of conduct].

[code of conduct]: https://github.com/fossapps/golang_starter/blob/master/CODE_OF_CONDUCT.md

Before contribution, please make sure there's a issue created (if there's none, create one!),
and wait for response, unless a issue has a "ready" label, it won't be accepted)

If a issue isn't accepted add your opinion and make sure your bug report or feature request is actually ready to be worked on.

Once the issue is labelled as "accepted", it should already contain enough information,
if there's none present, one can always ask by commenting, you can start working on it after you mention so in the issue.
This is to eliminate conflicting and duplicate work.

If you're new to this repository, looking for issue labeled with "Good First Issue" will be a good idea.


To begin contributing:

[Fork] this repository

[Fork]: https://github.com/fossapps/golang_starter/fork

Get the code

    git clone git@github.com:<your username>/golang_starter.git starter

Make sure you clone inside `$GOPATH/src/github.com/<username>/starter` (change name to from golang_starter to starter)

Get dependencies (we use dep for this project):

    dep ensure

Copy example env:

    cp .env.example .env

Start Docker

    docker-compose up // can pass -d if you don't want to use new terminal later)

Run migrations

    make migrate

Make sure the unit tests pass:

    make test

Make sure integration tests pass:

    make test-integration

Make your change (follow TDD whereever possible). Make the tests & lint pass:

    make test && make test-integration && make lint

Push to your fork and [submit a pull request][pr].

[pr]: https://github.com/cyberhck/pushy/compare/

At this point you're waiting for review which hopefully doesn't take long.

Some things that will increase the chance that your pull request is accepted:

* Write tests.
* Go style commit messages.

Commit message should follow the following convention:
```
[package name]: summary

[description]
```
package name shouldn't include github.com/fossapps part for obvious reasons.

summary should be quite short and consise.

there should be an empty line after first line, after which description should be written.
