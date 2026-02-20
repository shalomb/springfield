--- Sandbox Capability Audit ---
Date: Thu 19 Feb 14:23:58 CET 2026
User: unop

Checking: docker
  Status: NOT FOUND

Checking: podman
  Status: FOUND
  Path: /usr/bin/podman
-rwxr-xr-x 1 root root 42410968 Dec 21 17:42 /usr/bin/podman
podman version 5.4.2

Checking: bwrap
  Status: NOT FOUND

Checking: nsenter
  Status: FOUND
  Path: /usr/bin/nsenter
-rwxr-xr-x 1 root root 133696 May 10  2025 /usr/bin/nsenter
nsenter from util-linux 2.41

Checking: chroot
  Status: NOT FOUND

--- System Info ---
Linux idun 6.1.0-9-arm64 #1 SMP Debian 6.1.27-1 (2023-05-08) aarch64 GNU/Linux
uid=1001(unop) gid=1001(unop) groups=1001(unop),24(cdrom),25(floppy),27(sudo),29(audio),30(dip),46(plugdev),108(netdev),113(bluetooth),115(docker),117(lpadmin),120(scanner),993(ollama),1000(parallels)
