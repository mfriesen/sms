
## SMS - Service Monitoring System

**Service Monitoring System** is an application that makes it easy to monitor the status of services across different systems and platforms.

### Tested Platforms

| From  | To | Note |
| ------------- | ------------- | ------------- |
| Windows  | Windows  | | 
| Windows  | Linux  | * connects to Linux via SSH |
| Linux  | Windows  | * requires SAMBA 'net' execuable |

### Usage
```
  sms [options] [user@]<host>[:port] <servicename> start
  sms [options] [user@]<host>[:port] <servicename> status
  sms [options] [user@]<host>[:port] <servicename> stop

 Options:
  --user=userid  userid
  --password=password  password
  --sudo=sudopw  sudo password
  -h, --help     show help
  -v, --verbose  show debug info
```

 ### Examples

 Get the status of a Linux Service (requires the Linux Server is running SSH)

 * sms myuser@myhost myservice status 

#### Get the status of a Linux Service that requires a "sudo" password (requires the Linux Server is running SSH)

```
sms --sudo=MYSUDOPASSWORD myuser@myhost myservice status 
or 
sms --sudo= myuser@myhost myservice status (will prompt for a SUDO password)
```

#### Get the status of a Windows Service

* sms myhost myservice status

### Contributing

We love contributions! If you'd like to contribute please submit a pull request via Github.

### LICENSE

This library is distributed under the **MIT Open Source License**.