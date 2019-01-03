## Contributing Guidelines

Thank you for considering to contribute to this application!

The following text lists guidelines for contributions.
These guidelines don't have legal status, so use them as a reference and common sense - and feel free to update them as well!


### "I just want to know..."

For questions, or general usage-information about HackEd, please refer to the [systemshock.org](https://systemshock.org) forums (`Engineering` subforum for modding in general, `Deck 13` for source development, including engine code), or refer to the [Wiki](https://github.com/inkyblackness/hacked/wiki).

You might also find information on the data files in the [ss-specs](https://github.com/inkyblackness/ss-specs) project.


### Scope

This is an editor for mods and fan-missions of the video game System Shock (1994, and compatible), written in Go.

If there are extensions to the data files, please be sure that these were discussed and/or documented in `ss-specs`. Usability features are always welcome, I only ask to give them a second thought if there might be a more generic solution (i.e., don't immediately jump to the implementation of "I need X", rather think how a slightly different feature might both give you X and also Y in even a simpler form...)

Please also consider engine compatibility. If there are engine features not supported by all (especially the classic version), please make sure these are indicated to the user. Also refer to the documentation on the Wiki for such cases.


### Code Style

Please make sure code is formatted according to `go fmt`, and use the following linter: [golangci-lint](https://github.com/golangci/golangci-lint).

> If there are linter errors that you didn't introduce, you don't have to clean them up - I might have missed them and will be handling them separately.


### Testing Style

The code under `ss1` should be well tested, please try to keep up coverage. The code under `editor` is tricky to test - if you figure out how to test ImGui based stuff, I'd be happy for a pull request! At least utilities should have tests.

> I am aiming for a coverage of over 90% in there.

For testing use either the vanilla `testing` package, or augment it with the `github.com/stretchr/testify` framework. See existing examples.
