# GOPF: GNU Online Player Framework.

GOPF is a framework for making available your media online.  If you have audio *and* video files to display, then you have to run two instances of this site.

GOPF works on a LAMP server, and takes advantage of HTML5's `<audio>` and `<video>` tags.  Only Chrome is actively supported, because it came out with support for those tags long before the other major browsers.

Please don't hesitate to change, distribute, and make this program available.  It's free software!  Let it be free!

## Installing and Running

Use the [Go](https://golang.org) build system!  From the root directory of your project, type:
```sh
$ export GOPATH=$GOPATH:`pwd`
$ go install gopf
$ bin/gopf
```

You should see it print out:
```
Running server on :8079
```
If you see that, you'll know everything worked!

## Support

If you'd like to show your appreciation, you can donate to any one of these addresses:

```
 Bitcoin (BTC):  1DYeoForwthTgfnhUrkdwCRWL3cSBQwUje
Litecoin (LTC):  LcSBUQdjj3nbgPFBejfP7rExYivHobxR2a
Dogecoin (DOGE): DUHUB6m1DstTEbrur3xAMNvCEv5nCnRWct
```

## Bugs

To submit a feature request or to report a bug, go to http://github.com/gopf and click on "Issues".  Enter your feature or bug in the space provided, and please be sure to include any relevant information -- e.g., how your server is configured, what you did, and what you expected to happen, etc.  Also, labeling your issue with a relevant label will ensure that it gets addressed as quickly as possible.
