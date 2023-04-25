vacproxy is a simple http proxy with basic authentication support.

Installing from source
----------------------

To install, run

    $ go get github.com/vacuityv/vacproxy

You will now find a `vacproxy` binary in your `$GOPATH/bin` directory.

Usage
-----

Start proxy

    $ ./vacproxy

Run `vacproxy -help` for more information.

    $ vacproxy -help
    Usage of ./vacproxy:
        -bind string
            proxy bind address (default "0.0.0.0:7777")
        -config string
            config file (default "./config.yml")
        -log string
            the log file path (default "./vacproxy.log")
        -pid string
            the pid file path (default "./vacproxy.pid")
        -q  
            quit proxy
        -s string
            Send signal to the daemon:
                stop — shutdown, same as -q
                reload — reloading the configuration file

    

