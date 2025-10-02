# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0](https://github.com/jycamier/wrapper/compare/v1.0.0...v1.1.0) (2025-10-02)


### Features

* add list command ([8e3a196](https://github.com/jycamier/wrapper/commit/8e3a1961b178576b6b4ab2b646b46fc92eb726a8))
* finalize ([b38e940](https://github.com/jycamier/wrapper/commit/b38e940998545cad0e5dafe4c5f9ae963b909070))
* mvp ([0c85342](https://github.com/jycamier/wrapper/commit/0c85342c0fdc7a9dac7403c4b64dbdc6036dc18c))
* releasing ([150a8d3](https://github.com/jycamier/wrapper/commit/150a8d3e7d6c2b81f8966bb69e8ab94237e4b1bf))

## 1.0.0 (2025-10-02)


### Features

* add list command ([8e3a196](https://github.com/jycamier/wrapper/commit/8e3a1961b178576b6b4ab2b646b46fc92eb726a8))
* finalize ([b38e940](https://github.com/jycamier/wrapper/commit/b38e940998545cad0e5dafe4c5f9ae963b909070))
* mvp ([0c85342](https://github.com/jycamier/wrapper/commit/0c85342c0fdc7a9dac7403c4b64dbdc6036dc18c))
* releasing ([150a8d3](https://github.com/jycamier/wrapper/commit/150a8d3e7d6c2b81f8966bb69e8ab94237e4b1bf))

## [Unreleased]

### Added
- Initial release of wrapper
- DDD architecture with clean separation of concerns
- Profile management for any CLI binary
- Commands: `list`, `alias`, `profile`, `profile create`, `profile set`
- Color-coded profile listing with emojis
- Shell alias generation for bash, zsh, and fish
- Binary execution with environment profiles
- Filesystem-based profile storage in `~/.config/wrapper/`
