# Translations README

This is a basic sketch of the workflow needed to add translations:

# Adding/Updating Translations

## New languages
Create `./translations/runr/<language>/LC_MESSAGES/k8s.po`. 
There's no need to update `translations/test/...` which is only used for unit tests.

Once you've added a new language, you'll need to register it in
`./i18n.go` by adding it to the `knownTranslations` map.

## Wrapping strings
There is a simple script in `./hack/extract.py` that performs
simple regular expression based wrapping of strings. It can always
use improvements to understand additional strings.

## Extracting strings
Once the strings are wrapped, you can extract strings from go files using
the `go-xgettext` command which can be installed with:

```console
go get github.com/gosexy/gettext/go-xgettext
```

Once that's installed you can run `./hack/update-translations.sh`, which
will extract and sort any new strings.

## Adding new translations
Edit the appropriate `k8s.po` file, `poedit` is a popular open source tool
for translations. You can load the `./translations/runr/template.pot` file
to find messages that might be missing.

Once you are done with your `k8s.po` file, generate the corresponding `k8s.mo`
file. `poedit` does this automatically on save, but you can also run
`./hack/update-translations.sh` to perform the `po` to `mo` translation.

We use the English translation as the `msgid`.

## Extracting strings

There is a script in `./hack/extract.py` that knows how to do some
simple extraction. It needs a lot of work.

# Using translations

To use translations, you simply need to add:
```go
import pkg/i18n
...
// Get a translated string
translated := i18n.T("Your message in english here")

// Get a translated plural string
translated := i18n.T("You had % items", items)

// Translated error
return i18n.Error("Something bad happened")

// Translated plural error
return i18n.Error("%d bad things happened")
```
