# Installing Singularity development-3.0

Since you are reading this from the Singularity source code, it will be assumed
that you are building/compiling.

For full instructions on installation, check out our
[installation guide](https://www.sylabs.io/guides/3.0/user-guide/installation.html).

## Install system dependencies

You must first install development and libraries to your host.
Assuming Ubuntu:

```bash
$ sudo apt-get update && \
  sudo apt-get install -y build-essential \
  libssl-dev uuid-dev libgpgme11-dev squashfs-tools libseccomp-dev pkg-config
```

On CentOS/RHEL:

```bash
$ sudo yum groupinstall -y 'Development Tools' && \
  sudo yum install -y openssl-devel libuuid-devel libseccomp-devel
```
Skip `libseccomp-devel` on CentOS/RHEL 6.

## Install golang

This is one of several ways to [install and configure golang](https://golang.org/doc/install).

First, visit the [golang download page](https://golang.org/dl/) and pick a
package archive to download. Copy the link address and download with `wget`.

Then extract the archive to `/usr/local` (or use other instructions on go
installation page).

```bash
$ export VERSION=1.11.4 OS=linux ARCH=amd64

$ cd /tmp  # or wherever you want to store golang

$ wget https://dl.google.com/go/go$VERSION.$OS-$ARCH.tar.gz && \
  sudo tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```

Finally, set up your environment for Go:

```bash
$ echo 'export GOPATH=${HOME}/go' >> ~/.bashrc && \
  echo 'export PATH=/usr/local/go/bin:${PATH}:${GOPATH}/bin' >> ~/.bashrc && \
  source ~/.bashrc
```

## Clone the repo

golang is a bit finicky about where things are placed. Here is the correct way
to build Singularity from source:

```bash
$ mkdir -p ${GOPATH}/src/github.com/sylabs && \
  cd ${GOPATH}/src/github.com/sylabs && \
  git clone https://github.com/sylabs/singularity.git && \
  cd singularity
```

## Compile the Singularity binary

Now you are ready to build Singularity. Dependencies will be automatically
downloaded. You can build Singularity using the following commands:

```bash
$ cd ${GOPATH}/src/github.com/sylabs/singularity && \
  ./mconfig && \
  cd ./builddir && \
  make && \
  sudo make install && \
```

And Thats it! Now you can check you Singularity version by running:

```bash
$ singularity version
```

<br>

Alternatively, to build an RPM on CentOS/RHEL use the following commands:

```bash
$ sudo yum install -y rpm-build wget

$ cd ${GOPATH}/src/github.com/sylabs/singularity && \
  ./mconfig && \
  make -C builddir rpm && \
```

Golang doesn't have to be installed to build an rpm because the rpm
build installs golang and all dependencies, but it is still recommended
for a complete development environment.

To build a stable version of Singularity, check out a [release tag](https://github.com/sylabs/singularity/tags) before compiling:

```bash
$ git checkout v3.0.0
```

To build in a different folder and to set the install prefix to a different path:

```bash
$ ./mconfig -p /usr/local -b ./buildtree
```
