APPNAME=gopr
PKGBASE=github.com/kmpm/gopr
APPVERSION?=v0.0.0

# try to be os agnostic
ifeq ($(OS),Windows_NT)
	# aliases
	RM = del /q
	RMDIR = rmdir /s /q
	DEVNUL := NUL
	WHICH := where
	CAT = type
	MV = move
	CPR = xcopy /s /e /y /i
	MKDIR = mkdir
	CP = copy
	ZIP = @echo Windows should learn how to zip
	NEWLINE = echo.
	# functions
	FixPath = $(subst /,\,$1)
	Sleep = ping 192.0.2.0 -n 1 -w $1000

	# variables
	CGO_windows ?= 1
	CGO_linux ?= 0
	CGO_darwin ?= 0
	BINEXT = .exe
else
	# aliases
	RM = rm -f
	RMDIR = rm -Rf
	CAT = cat
	DEVNUL = /dev/null
	WHICH = which
	MV = mv
	MKDIR = mkdir -p
	CPR = cp -r
	CP = cp
	ZIP = zip -rq 
	NEWLINE = printf "\n"
	# functions
	FixPath = $1
	Sleep = sleep $1

	BINEXT = 
	# os specific
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		CGO_windows ?= 1
		CGO_linux ?= 1
		CGO_darwin ?= 0
	endif
	ifeq ($(UNAME_S),Darwin)
        CGO_windows ?= 1
		CGO_linux ?= 1
		CGO_darwin ?= 1
    endif
endif

GIT_VERSION ?= $(shell git rev-parse HEAD)
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_TAG?=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_DATE?=$(shell git log -1 --format="%ad" --date="format:%Y%m%d-%H%M%S")

LDFLAGS="-X '$(PKGBASE)/cmd.gitVersion=$(GIT_VERSION)' -X '$(PKGBASE)/cmd.appVersion=$(APPVERSION)'"

DISTDIR?=$(call FixPath,./dist)
BINFILE=$(call FixPath,$(DISTDIR)/$(APPNAME)$(BINEXT))
BINARCHIVE=$(call FixPath,$(DISTDIR)/$(APPNAME)-$(APPVERSION)$(BINEXT))
BUILDCMD=go build -ldflags $(LDFLAGS)

.PHONY: build clean

build: $(DISTDIR)
	$(BUILDCMD) -o $(BINFILE) .
	$(CP) $(BINFILE) $(BINARCHIVE)

$(DISTDIR):
	$(MKDIR) $@


clean:
	$(RMDIR) $(DISTDIR)
