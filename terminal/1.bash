$ go build -o rabbit
$ ./rabbit 
# Usage:
#   rabbit [flags]

# Flags:
#       --config string           config file (default is /etc/.talks.meshcon.23.pito.yaml) (default "/etc/.talks.meshcon.23.pito.yaml")
#       --go-max-procs int        How many processes to assign the Go runtime (default 1)
#   -h, --help                    help for rabbit
#       --listen-address string   The address on which the server should listen (default "localhost:80")
$ echo $?
# 1