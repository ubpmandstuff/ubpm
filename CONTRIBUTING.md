# contributing to ubpm

first of all, thanks for considering contributing to this project. however,
there are a few things to know first.

## workflow

we use a rather simple feature-branch workflow. instead of pushing to `main`
directly, you should create a new branch for your feature. we advise you follow
this format:

`<feat|fix>/<description>`

where `<feat|fix>` is the type, typically `feat` for a feature or `fix` for a
bugfix, and `<description>` is the description.

## commits

i would prefer you use a simple format for commits: `<scope>: <description>`.

 alternatively, you can use 
 [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/).

the general format for conventional commits is
`<type>[optional scope]: <description>`.

typically, the following types are enough:
- `feat`: a new feature (e.g. `feat: add list subcommand`)
- `fix`: a bugfix (e.g. `fix: stop generating new nonce on decrypt attempt`)
- `chore`: maintenance tasks (e.g. `chore: update version`)
- `refactor`: refactoring, formatting, etc. (e.g. `refactor: change run function name`)
- `docs`: documentation (e.g. `docs: add documentation for InitVault`)

