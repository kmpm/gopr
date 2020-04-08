APPNAME=gopr
PKGBASE=github.com/kmpm/gopr
APPVERSION?=v0.0.2

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

GIT_HASH?=$(shell git rev-parse HEAD)
GIT_BRANCH?=$(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_TAG?=$(shell git describe --abbrev=0 --tags --always --match "v*")
GIT_DATE?=$(shell git log -1 --format="%ad" --date="format:%Y%m%d-%H%M%S")
# ifeq ($(GIT_BRANCH), master) 
# GIT_BRANCH:=
# else
# APPVERSION:=$(APPVERSION)-$(GIT_BRANCH)
# endif

LDFLAGS="-X '$(PKGBASE)/cmd.appVersion=$(APPVERSION)' -X '$(PKGBASE)/cmd.gitHash=$(GIT_HASH)' -X '$(PKGBASE)/cmd.gitBranch=$(GIT_BRANCH)'"

DISTDIR?=$(call FixPath,./dist)
BINFILE=$(call FixPath,$(DISTDIR)/$(APPNAME)$(BINEXT))
BINARCHIVE=$(call FixPath,$(DISTDIR)/$(APPNAME)-$(APPVERSION)$(BINEXT))
BUILDCMD=go build -ldflags $(LDFLAGS)

.PHONY: build clean help version

build: $(DISTDIR)
	$(BUILDCMD) -o $(BINFILE) .
	$(CP) $(BINFILE) $(BINARCHIVE)

version: 
	go run -ldflags $(LDFLAGS) . version

$(DISTDIR):
	$(MKDIR) $@

clean:
	$(RMDIR) $(DISTDIR)

help:
	@echo APPNAME=$(APPNAME)
	@echo APPVERSION=$(APPVERSION)
	@echo GIT_BRANCH=$(GIT_BRANCH)
	@echo GIT_TAG=$(GIT_TAG)
