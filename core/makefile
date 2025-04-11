# From: https://tech.davis-hansson.com/p/make/
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
# Tabs are annoying
.RECIPEPREFIX = >
# End From



downloadTestFiles: removeTestFiles
> mkdir -p test_data
> git clone https://github.com/BSVino/docs.gl.git test_data/docs.gl
> git clone https://github.com/tldr-pages/tldr.git test_data/tldr
> git clone https://github.com/neovim/doc.git test_data/neovim -b gh-pages
> git clone git@github.com:mozilla-l10n/documentation.git test_data/mozilla-l10n

removeTestFiles:
> rm -rf test_data
