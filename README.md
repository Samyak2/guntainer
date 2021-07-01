# guntainer

A minimal rootless container implementation on Linux. ~~Copied from~~ Inspired by Liz Rice's amazing [talk](https://youtu.be/8fi7uSYlOdc) on implementing containers from scratch.

## Getting started

Install using (Go 1.16+ recommended):
```
go install github.com/Samyak2/guntainer@latest
```
(ensure `GOBIN` is in path)

Usage:
```
guntainer run <archive_of_root_FS> <command_to_run> <args_for_command>...
```

 - `archive_of_root_FS` is an archive (tar, zip, etc.) of a root filesystem which will be `chroot`ed into
 - `command_to_run` is a program *inside* the container root FS that is run (note that `PATH` is not set unless you execute a shell)
 - `args_for_command` are arguments for the program (everything after the `command_to_run` is passed directly to the program)

## Examples

### Alpine

Get the Alpine "mini root filesystem" from [here](https://alpinelinux.org/downloads/) ([direct link](https://dl-cdn.alpinelinux.org/alpine/v3.14/releases/x86_64/alpine-minirootfs-3.14.0-x86_64.tar.gz) of the specific version this was tested with).

Run the container with (replace the archive path if necessary) (you can also use `sh` instead of `ash`):

```sh
guntainer run alpine-minirootfs-3.14.0-x86_64.tar.gz /bin/ash
```

Most programs will work. Running `hostname` should say `guntainer`.

Internet will not work out of the box as no DNS servers are configured. Use the following to access internet (replace the IP address as necessary):
```sh
echo "nameserver 8.8.8.8" > /etc/resolv.conf
```

### Ubuntu

Get the Ubuntu Base image from [here](http://cdimage.ubuntu.com/ubuntu-base/releases/) ([direct link](http://cdimage.ubuntu.com/ubuntu-base/releases/21.04/release/ubuntu-base-21.04-base-amd64.tar.gz) to the specific version used).

Run using:
```sh
guntainer run ubuntu-base-21.04-base-amd64.tar.gz /bin/bash
```

Running `hostname` should say `guntainer`.

Issues:
 - `apt` will not work out of the box. Refer [this](https://github.com/opencontainers/runc/issues/2517#issuecomment-657139999) thread for [details](https://github.com/opencontainers/runc/issues/2517#issuecomment-764163674).

    Workaround (pls don't run this outside a container):
    ```sh
    sed -i '/_apt/d' /etc/passwd
    ```
 - Internet will not work out of the box. Workaround:
    ```sh
    echo "nameserver 8.8.8.8" > /etc/resolv.conf
    ```

## TODO

 - [ ] Better CLI using [cobra](https://github.com/spf13/cobra)
 - [ ] Use cgroups for resource limits
 - [ ] (a bit ambitious) be OCI compliant
 - [ ] Fix Ubuntu issues

## Resources

If you're looking to implement your own container runtime, these links are great to start with:
 - [Mythili Vutukuru's lecture](https://youtu.be/4BG-hE72r_I) on containers - provides a good overview of the Linux concepts behind containers (namespaces and cgroups)
 - [Containers From Scratch by Liz Rice](https://youtu.be/8fi7uSYlOdc) - most of the code is from this talk. Liz Rice implements it live on the stage while explaining how it works.
    - The corresponding [code respository](https://github.com/lizrice/containers-from-scratch) is a good reference and also links to another slide deck for implementing rootless containers

## Implementation details

 - The root FS archive is extracted in a temporary directory and `chroot`ed into. The directory is cleaned up once the container exits.
 - Rootless is implemented by mapping the current user's UID and GID to 0 (root) inside the container. This means that inside the container you are root while the same user outside the container is your user.

## License

MIT
