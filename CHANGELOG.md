# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased][]

## [0.4.0][] - 2018-11-02

### Added

- Field `Download.dsaSignature` with getter and setter to hold the DSA signature
of the file
- Field `Release.minimumSystemVersion` with getter and setter to hold the
minimum required system version
- Field `Release.releaseNotesLink` with getter and setter to hold the release
notes link
- Interface `Downloader` implemented by `Download`
- Support for the "Sparkle RSS Feed" new `<enclosure />` attributes:
`dsaSignature`, `md5sum`, `minimumSystemVersion` and `releaseNotesLink`
- Support passing `dsaSignature` and `md5` as parameters in the `NewDownload`
function

### Changed

- All `Download` fields to become unexported
- Client-specific stuff to be in the separate `client` package
- Field `Download.Type` in favour of `Download.filetype`
- Field `Download.URL` in favour of `Download.url`
- Field `PublishedDateTime.time` to become a `*time.Time` type
- Function `NewClient` to become the `New` in the `client` package
- Method `GitHubAtomFeedAppcast.UnmarshalReleases` to return an error when no
source or unmarshalling failure
- Method `PublishedDateTime.String` to return an empty string for the `nil`
value
- Method `SourceForgeRSSFeedAppcast.UnmarshalReleases` to return an error when
no source or unmarshalling failure
- Method `SparkleRSSFeedAppcast.UnmarshalReleases` to return an error when no
source or unmarshalling failure
- Unmarshalling structs for "SourceForge RSS Feed" and "GitHub Atom Feed" to
become unexported

### Fixed

- Method `GitHubAtomFeedAppcast.UnmarshalReleases` to set the
`GitHubAtomFeedAppcast.source.appcast`
- Method `SourceForgeRSSFeedAppcast.UnmarshalReleases` to set the
`SourceForgeRSSFeedAppcast.source.appcast`
- Method `SparkleRSSFeedAppcast.UnmarshalReleases` to set the
`SparkleRSSFeedAppcast.source.appcast`

## [0.3.0][] - 2018-10-27

### Added

- Extendable `Output` with `Outputer` interface for creating use-case specific
outputs
- Field `Output.appcast` with getter and setter to hold the provider-specific
appcast after marshalling
- Field `Source.appcast` with getter and setter to hold the provider-specific
appcast after unmarshalling
- Source `LocalOutput` with `LocalOutputer` interface to save an appcast to the
local file by path
- Struct `PublishedDateTime` to use as the `Release.publishedDateTime` type in
the `release` package
- Unmarshalling support for the "Sparkle RSS Feed" channel as the
`SparkleRSSFeedAppcast.channel`

### Changed

- Code coverage service from "Coveralls" to "Codecov"
- Dependencies versions to match the latest ones
- Field `Release.publishedDateTime` type to become the new `PublishedDateTime`
- Function `NewRelease` to become the `New` in the `release` package
- Method `Appcast.LoadFromLocalSource` to also return the provider-specific
appcast
- Method `Appcast.LoadFromRemoteSource` to also return the provider-specific
appcast
- Method `Appcast.UnmarshalReleases` to also return the provider-specific
appcast
- Method `GitHubAtomFeedAppcast.UnmarshalReleases` to also return the
provider-specific appcast
- Method `SourceForgeRSSFeedAppcast.UnmarshalReleases` to also return the
provider-specific appcast
- Method `SparkleRSSFeedAppcast.UnmarshalReleases` to also return the
provider-specific appcast
- Release to store the original time and format in the `Release.publishedDateTime`
- Release-specific stuff to be in the separate `release` package
- Testdata published releases dates to match the real ones in the past
- Unmarshalling structs for the "Sparkle RSS Feed" to become unexported

### Removed

- Deprecated `Appcast` methods
- Deprecated `GuessProviderFromContent` and `GuessProviderFromURL` functions
- Deprecated `Release` methods

## [0.2.0][] - 2018-08-09

### Added

- Convenience method `Appcast.LoadSource` to call the `Appcast.Source.Load`
methods chain
- Extendable `Source` with `Sourcer` interface for creating use-case specific
sources
- Field `Appcast.source` to hold the source-specific data instead of the removed
`Appcast` fields
- Function `GuessProviderByContent` for the `[]byte` type similar to the
deprecated `GuessProviderFromContent`
- Getters and setters for all unexported `Appcast` fields
- Getters and setters for all unexported `Release` fields
- Getters for all unexported `Checksum` fields
- Interface `Appcaster` implemented by `Appcast`
- Interface `GitHubAtomFeedAppcaster` implemented by `GitHubAtomFeedAppcast`
- Interface `Releaser` implemented by `Release`
- Interface `SourceForgeRSSFeedAppcaster` implemented by `SourceForgeRSSFeedAppcast`
- Source `LocalSource` with `LocalSourcer` interface to load an appcast from the
local file by path
- Source `RemoteSource` with `RemoteSourcer` interface to load an appcast from
the remote location by URL
- Source support for the `New` function
- This CHANGELOG.md

### Changed

- All `Appcast` fields to become unexported
- All `Checksum` fields to become unexported
- All `Release` fields to become unexported
- Always extract releases when calling `Appcast.LoadFrom...` methods
- Always generate SHA256 checksum when calling `Appcast.LoadFrom...` methods
- Field types for `Appcast.releases` and `Appcast.originalReleases` in favour of
`Releaser` interface
- Method `Release.SetVersion` to `Release.SetVersionString`
- Struct `BaseAppcast` name to `Appcast`

### Deprecated

- Function `GuessProviderFromContent` in favour of `GuessProviderByContentString`
- Function `GuessProviderFromURL` in favour of `GuessProviderByUrl`
- Method `Appcast.ExtractReleases` in favour of `Appcast.UnmarshalReleases`
- Method `Appcast.GenerateChecksum` in favour of `Appcast.GenerateSourceChecksum`
- Method `Appcast.GetChecksum` in favour of `Appcast.Source.Checksum` methods
chain
- Method `Appcast.GetFirstRelease` in favour of `Appcast.FirstRelease`
- Method `Appcast.GetProvider` in favour of `Appcast.Source.Provider` methods
chain
- Method `Appcast.GetReleasesLength` in favour of `len(Appcast.Releases)`
- Method `Appcast.LoadFromFile` in favour of `Appcast.LoadFromLocalSource`
- Method `Appcast.LoadFromURL` in favour of `Appcast.LoadFromRemoteSource`
- Method `Release.GetBuildString` in favour of `Release.Build`
- Method `Release.GetDownloads` in favour of `Release.Downloads`
- Method `Release.GetVersionOrBuildString` in favour of `Release.VersionOrBuildString`
- Method `Release.GetVersionString` in favour of `Release.Version.String`
methods chain

### Removed

- Source-specific `Appcast` fields in favour of `Appcast.source` field

### Fixed

- Method `Appcast.Uncomment` to correctly handle the return of the `Unknown`
provider
- Method `SparkleRSSFeedAppcast.Uncomment` to match the `Appcaster` interface

## 0.1.0 - 2018-08-04

First release.

[unreleased]: https://github.com/victorpopkov/go-appcast/compare/v0.4.0...HEAD
[0.4.0]: https://github.com/victorpopkov/go-appcast/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/victorpopkov/go-appcast/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/victorpopkov/go-appcast/compare/v0.1.0...v0.2.0
