# guntainer

A minimal rootless container implementation on Linux. ~~Copied from~~ Inspired by Liz Rice's amazing [talk](https://youtu.be/8fi7uSYlOdc) on implementing containers from scratch.

## Features

 - *Rootless*: never use `sudo` to run a container
 - *Images*: guntainer can build container images, with a `Gunfile` to define the structure.

## Getting started

Install using (Go 1.16+ recommended):
```
go install github.com/Samyak2/guntainer@latest
```
(ensure `GOBIN` is in path)

To run a container:
```
guntainer run <archive_of_root_FS> <command_to_run> <args_for_command>...
```

 - `archive_of_root_FS` is an archive (tar, zip, etc.) of a root filesystem which will be `chroot`ed into
 - `command_to_run` is a program *inside* the container root FS that is run (note that `PATH` is not set unless you execute a shell)
 - `args_for_command` are arguments for the program (everything after the `command_to_run` is passed directly to the program)


To build a new image, first make a [`Gunfile`](#building-images) and then use:
```
guntainer build Gunfile <path_to_new_image>
```

Where `Gunfile` can be replaced with the path to the Gunfile and `path_to_new_image` is the path where the generated image is saved (as a tar file)

## Building images

guntainer can build new images out of existing ones, in a similar way to `docker build`. Images are described using a `Gunfile`, whose format is very much a work-in-progress. Following is the structure of a Gunfile (at least the things that currently work):
```Gunfile
Using "<archive_of_root_FS>"

Exec "some command"
```

More examples can be found [here](./examples/).

To build the image from [example_02](./examples/02_alpine_vim/), we can use:
```
guntainer build examples/02_alpine_vim/Gunfile example_02.tar
```

This will generate an `example_02.tar` which is the newly built image with `vim` installed. Run it using:
```
guntainer run example_02.tar /bin/sh
```

You should be able to use `vim` inside the container now.

## `run` Examples

These are examples of running an existing Linux distro inside guntainer.

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

 - [x] Better CLI using [cobra](https://github.com/spf13/cobra)
 - [ ] Use cgroups for resource limits
 - [ ] (a bit ambitious) be OCI compliant
 - [ ] Fix Ubuntu issues
 - [x] Dockerfile equivalent
 - [ ] Download images from URL, like go's package management
 - [ ] Figure out how to store metadata along with the built image. Could do it similar to docker images (OCI) or do something more hacky.
 - [ ] Optional logging - `-v` flag should enable more logs.

## Resources

If you're looking to implement your own container runtime, these links are great to start with:
 - [Mythili Vutukuru's lecture](https://youtu.be/4BG-hE72r_I) on containers - provides a good overview of the Linux concepts behind containers (namespaces and cgroups)
 - [Containers From Scratch by Liz Rice](https://youtu.be/8fi7uSYlOdc) - most of the code is from this talk. Liz Rice implements it live on the stage while explaining how it works.
    - The corresponding [code respository](https://github.com/lizrice/containers-from-scratch) is a good reference and also links to another slide deck for implementing rootless containers

## Implementation details

 - The root FS archive is extracted in a temporary directory and `chroot`ed into. The directory is cleaned up once the container exits.
 - Rootless is implemented by mapping the current user's UID and GID to 0 (root) inside the container. This means that inside the container you are root while the same user outside the container is your user.
 - The Gunfile is implemented [here](./gunfile/). It uses [participle](https://github.com/alecthomas/participle) to parse the Gunfile and build an AST.
 - To save the built image I had to implement a custom tar wrapper to handle in-container symlinks. I called it [guntar](./guntar/).

## License

MIT
