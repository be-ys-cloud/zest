# Zest

A lightweight and powerful Debian cache repository.

## What is Zest ?

Zest is a web server who acts as a proxy for Debian repositories. When a user asks for a package, Zest try to serve it
from his disk. If it can't, he goes search it online, and stores it. The main difference with other tools
like [apt-cacher-ng](https://wiki.debian-fr.xyz/Apt-cacher-ng) is that Zest is able to re-create repository indexes to
add old versions to newer indexes. It enables people to keep packages up to date AND still download old versions for
specific usage.

An important point to know is that Zest **is not** a mirror of your repositories ! Zest will only store the packages you
need. If you never downloaded a package before his deletion on distant repositories, you will not find him here.

Zest works without the needs of an administrator: if a repository doesn't exist, it will be automatically created and
imported when a user will request it. Cache expiration is also automatically handled.

It is important to know that Zest needs a PGP key to work: as we are modifying indexes files to add old versions,
we must recreate Packages files. A "pass through" mechanism is available, but Zest will not be able to serve old
packages via this option.

The main differences between Zest and other package management are provided in the chart below :

| /                                                                         | Zest | apt-cacher-ng | Aptly |
|---------------------------------------------------------------------------|------|---------------|-------|
| Store multiple versions of a package and serve all of them simultaneously | Yes  | No            | No    |
| Mirrors a whole repository                                                | No   | No            | Yes   |
| Needs to sign a new PGP key on client machines                            | Yes  | No            | Yes   |
| Can handle both HTTP and HTTPS repositories                               | Yes  | No            | Yes   | 
| Can acts as a simple proxy, without re-signing files                      | Yes  | Yes           | No    |

In order to keep performances, a lot of tasks are done in subprocesses. You can check number of subprocesses through
the `/metrics` endpoints.

If you need to stop a Zest server, we strongly encourage you to use the `/stop` endpoint, witch will ensure that all
subprocesses ended before killing the server.

Please note that Zest does not natively handle HTTPS handshake for client connections. If you need to encrypt the
network between your machines and your Zest server, we encourage you to use a proxy server, like Traefik.

## Running a Zest server

**There is some system dependencies that you must take care about. To be functional, you need to have :**

* bash
* grep
* sed

We will remove these system dependencies in the future.

_____________

Running a Zest server is conceived to be the simplest possible.

* First, you will need a PGP key. You can create one through OpenGPG CLI, or a tool like Kleopatra.
    * Export your public key and give it to your clients (for details, see Chapter "Configuring Clients")
    * Export your private key, it will be useful for the server
* Secondly, you must define the following environment variables :
    * `ZEST_KEY_PASSWORD` : The password of the Private Key you will use
    * You also can redefine the following environment variables :
        * `ZEST_PORT` : The port where the app will start on (default: `80`)
        * `ZEST_TEMP_STORAGE` : The place where Zest will store his working garbage (default: `tmp/`)
        * `ZEST_DATA_STORAGE` : The place where Zest will store his important data (default: `data/`)
        * `ZEST_KEY_FILE` : The location of the Private Key that will be used to sign repositories (default: `key.asc`)
        * `ZEST_MAX_ROUTINES` : The maximum number of subprocesses that can be launched by the program (default: `10`)
        * `ZEST_ADMIN_PASSWORD` : The password that will be used for admin (default: `12345`)
        * `ZEST_DATABASE_FILE` : The SQLite file name (default: `database.sql`)
        * `ZEST_MAX_RETENTION_TIME` : The amount of days that a file will be kept after his last download (
          default: `60`)
        * `ZEST_FREE_SPACE_THRESHOLD` : The percentage of free space that is required to automatically start a purge.
          For example, 10% means that purge will begin only if there is less than 10% of space available on disk. (
          default: `10`)
* Then, you just have to build & run : `go run ./cmd`

Note : you can also use Docker or a built binary.

## Configuring Clients

* Add the public key to your keyring : `cat pubkey.asc | apt-key add -`
* Update your `/etc/apt/sources.list` file (and files in `/etc/apt/sources.list.d/`) :
    * Remove all the `[trusted=yes]` and the `[signed-by=.........]` fields that you may encounter
    * Replace every URL to the one using proxy. `http://` and `https://` becomes `http/` and `https/`, and you must add
      your host at the beginning of each line.
        * For example, if your Zest server is up on `http://192.168.10.0/`:
            * Your line `deb http://ftp.de.debian.org/debian buster-backports main` will become
            * `deb http://192.168.10.0/http/ftp.de.debian.org/debian buster-backports main`
    * If you want to use proxy but not custom indexes (witch means you will not have access to old packages that are no
      longer in upstream), you can use `pass/http/` and `pass/https` :
        * For example, if your Zest server is up on `http://192.168.10.0/`:
            * Your line `deb http://ftp.de.debian.org/debian buster-backports main` will become
            * `deb http://192.168.10.0/pass/http/ftp.de.debian.org/debian buster-backports main`
* Run `apt update`, and you're done !

## Endpoints & Authentication.

Zest offers several endpoints for its own configuration and management.

* `GET /metrics` (No authentication) : Retrieve statistics in a Prometheus-readable format.
* `GET /admin/stop` (Authentication required) : Gracefully stops the server (wait for indexes to be up-to-date, and all
  goroutines done. During this time, some clients may encounter an error because server is closing).
* `GET /admin/cleanup` (Authentication required) : Force server to cleanup old files.
* `DELETE /admin/packages/<package_name>` (Authentication required) : Asks server to permanently delete a package. It
  may be useful, eg. in case of security advisory on this package. For example :  ̀DELETE
  /admin/packages/deb.debian.org/debian/pool/main/libc/libcap2/libcap2_2.44-1_i386.deb` will delete this file and its
  indexes.

## Zest & Security considerations

As Zest changes indexes files and PGP signatures of repositories, you must ensure that, as a client, you deeply trust
the owner of this server.

In the following chapter, we will indicate what does Zest do with repositories, and why it does so.
We will also tell you what's Zest does not do, and why you should be suspicious in some cases.

### What does Zest do ?

The main feature that concerns security is Release & InRelease modification & signing. As Zest keeps old packages in
indexes,
it must add some part of indexes in Packages files. These parts can be found in folders ending with `_partials/`,
in `dists` (sub-)folders.
When we add these files to the Packages files, we are changing his size, and his checksums (MD5/SHA1/SHA256/SHA512).
These hashes are referenced
in a `Release`, or `InRelease` file, depending on the repository. When checksums and/or size does not match, apt clients
stops using this repository.
To be able to provide you a repository with old versions, we so need to edit these Release files. And as these files are
signed via PGP, we also need to remove the old signature,
and add ours. That's why we have to sign repositories with a custom key.

**What should make you suspicious ?** :

* GPG errors on client : if you retrieve a BADSIG error on client, it means that someone probably compromised the
  Release/InRelease file and is trying to serve you a bad one.
* You should also be suspicious if you already configured your apt client to accept your key, and still receiving a
  NO_PUBKEY error
* If you are using the pass through option and still viewing a custom GPG Key

### What does Zest DO NOT DO ?

Zest will NEVER :

* Add packages to a repository that does not hold it (merging/cross-publishing packages)
* Change .deb files, checksum, or signature
* Force you to use an old version of a package