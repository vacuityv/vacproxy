vacproxy is a simple http proxy with basic authentication support.

Installing from source
----------------------

To install, run

    $ go get github.com/vacuityv/vacproxy

Build

    $ go install github.com/vacuityv/vacproxy/cmd/vacproxy 

You will now find a `vacproxy` binary in your `$GOPATH/bin` directory.

Usage
-----

Start proxy

    $ ./vacproxy

Run `vacproxy -help` for more information.

    $ vacproxy -help
    Usage of ./vacproxy:
    -auth string
        basic credentials(username:password)
    -bind string
        proxy bind address (default "0.0.0.0:7777")
    -daemon
        run as daemon
    -log string
        the log file path (default "./vacproxy.log")
    -pid string
        the pid file path (default "./vacproxy.pid")
    -q    quit proxy
    -s string
        Send signal to the daemon:
            quit — graceful shutdown
            stop — fast shutdown
            reload — reloading the configuration file
    

