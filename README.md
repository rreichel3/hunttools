# hunttools
## Iterate, quickly.

HuntTools (ht) is a CLI designed to help you perform repetitive operations easily.  
The expected input type is always a newline separated file and the default output type will return the same.   
Using this, you should spend less time trying to remember how to write a loop in bash and more time actually hunting

### Quick Download
<!-- This should point to the latest binary -->
#### Linux

Run the following and you should be good to go! 

```bash
wget -O ht https://github.com/rreichel3/hunttools/releases/download/0.1beta/linux_ht && chmod +x ./ht && ./ht
```

#### MacOS

Run the following and you should be good to go! 

```bash
wget -O ht https://github.com/rreichel3/hunttools/releases/download/0.1beta/mac_ht && chmod +x ./ht && ./ht
```

#### Windows
_Note: This is untested_
```bash
wget -O ht https://github.com/rreichel3/hunttools/releases/download/0.1beta/ht.exe 
```



### How to Use
The input expected is always a newline file. For example, the `ht ping -i infile.txt` command expects `infile.txt` to contain something like:
```
rj3.me
test.rj3.me
127.0.0.1
```
The default output will simply write the lines that are "positive" or in this case, responding to ping.  If you want more information on each one check out the `--verbose` flag.  

This enables us to easily chain commands together, writing the output of the ping to a file then running a subsequent command to get all of the hostnames.  

#### Main Usage
```
$ ./ht 
HuntTools (ht) is a CLI designed to help you perform repetitive operations easily.  
The expected input type is always a newline separated file and the default output type will return the same.   
Using this, you should spend less time trying to remember how to write a loop in bash and more time actually hunting

Usage:
  ht [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  hostnames   Get hostnames for a list of IP Addresses
  ping        Ping a list of hosts

Flags:
  -h, --help      help for ht
  -v, --verbose   Output full results, not newline formatted

Use "ht [command] --help" for more information about a command.
```

#### Ping
```
$ ./ht ping -h  
Runs the ping command against a list of hosts

Usage:
  ht ping [flags]

Flags:
  -h, --help            help for ping
  -i, --infile string   Newline delimited file of ping destinations
```
#### Hostname
```
./ht hostnames -h
Gets the hostnames for a list of IP Addresses.

Usage:
  ht hostnames [flags]

Flags:
  -h, --help            help for hostnames
  -i, --infile string   Newline delimited file of IPs for which to fetch their hostname
```


