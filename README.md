# gotftpd

gotftpd is a standalone, single-binary tftp server.

It came about from my occasional need to transfer something over tftp (usually OS images or configs with Cisco equipment), but I don't really want to set up a permanent tftp system service running all the time with a dedicated folder.

## Installation/usage

On Linux, in order to bind to the default tftp port 69 as normal user, you'll need to use `setcap`:
```
setcap cap_net_bind_service=+ep /path/to/gotftpd
```

Now, as long as it's in your PATH, you can simply run it from any folder:
```
[alex@laptop ~/stuff]$ gotftpd
Starting server on port 69
```

You can use the `-p <number>` option to specify the port that the server runs on, however most clients (especially embedded devices) will only try to connect to port 69.