#
Package main - provides a CLI interface for testing a metrics interface
according to it specification as provided by client. To use this program to
test all default test cases:

RUN: ./ciohw1 on the command line after you have built the said executable,
which with no other options will exercise the "default test cases". [to show
that list see -tsnameslist -tsnamelistrunseq options] For a multitude of
options invoke with option -h for more info. One immediate tip increasing X
in -verblvl=X option will provide more "inside" info including key items and
the why of failures. To see what tests are running use >0 for X with
-verblvl=X i.e. ./ciohw1 -verblvl=1 If all the requested tests pass then the
program exit value will be zero, for other possible values ... grep for
os.Exit and log.Fatal.

for examples: see file: EXAMPLES

INSTALL: via normal golang setup environment:

    go get github.com/phcurtis/ciohw1
    go get github.com/phcurtis/fn
    Then go build while in github.com/phcurtis/ciohw1
    Currently this has only been alpha tested in ubuntu.
