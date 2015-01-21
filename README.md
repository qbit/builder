builder
=======

Set of tools for doing CI type builds for OpenBSD ports

### Plan

This project should result in a three component app, server, cli and
daemon.

### bserver

Will keep track of jobs and their associated status (stati?). New jobs
will be added via a POST req that includes: title, desc, port and a
diff to be tested.

### bdaemon

This bit will fire requests to bserver asking for grabable jobs, pull
down said job.. apply the diff.. and run a dpb against the specific port.

Once a build is done, the status will be reported back to bserver.

### bcli

This app will be used to register new jobs:

    bcli -diff node-0.10.35.diff -title "node 0.10.32 -> 0.10.35" -desc "Bring node to the latest version" -port "lang/node"

It will post the data to bserver and the job will default to being "Grabable"