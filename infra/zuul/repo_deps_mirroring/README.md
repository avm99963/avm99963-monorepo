# Repository dependencies mirroring

Zuul depends on some external repositories. As of now, just
https://opendev.org/zuul/zuul-jobs/, which defines base jobs and roles which
are very useful to us and upon which our jobs are based.

However, on April 7th 2026, there was an outage in their git server that didn't
allow us to run jobs.

Thus, this package was created: Copybara migrations that allow us to mirror
third-party repos to our Gerrit server, so we aren't affected by their outages.
Using Copybara also allows us to apply patches should this need arise (e.g., we
already have some forked versions of upstream roles in our own
[zuul-jobs][our-zuul-jobs] repo; now these could be converted into patches).

## Initializing the mirrors

Initialize each mirror by running locally:

```bash
bazel run //infra/zuul/repo_deps_mirroring:mirror -- WORKFLOW_NAME --init-history --force
```

## Running a migration

Migrations (mirroring) will run automatically in CI in the future. However, it
can also be manually run:

```bash
bazel run //infra/zuul/repo_deps_mirroring:mirror -- WORKFLOW_NAME
```

If running locally, you will need to use a HTTP password or set the following
git configuration so that you can authenticate via SSH:

```bash
git config --global url."ssh://gerrit.avm99963.com:29418/".insteadOf "https://gerrit.avm99963.com/a/"
```

## Testing a migration

Running the following command will output the result of the migration to a
temporary folder locally:

```bash
bazel run //infra/zuul/repo_deps_mirroring:mirror -- WORKFLOW_NAME --to-folder --squash
```

[our-zuul-jobs]: https://gerrit.avm99963.com/plugins/gitiles/zuul/jobs/
