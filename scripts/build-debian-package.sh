#!/usr/bin/env bash
set -eux

BASE_DIR=$(cd $(dirname $(readlink -f $0)) && cd .. && pwd)
cd ${BASE_DIR}

# Generate changelog
git_describe="$(git describe --tags)"
VERSION=${git_describe:1}.$(TZ=JST-9 date +%Y%m%d)+$(lsb_release -cs)
DATE=$(LC_ALL=C TZ=JST-9 date '+%a, %d %b %Y %H:%M:%S %z')

cat <<EOF > "${BASE_DIR}/debian/changelog"
git-caddy (${VERSION}) unstable; urgency=medium

  * This is an automated build.

 -- Sigmonsays <noreply@example.net>  ${DATE}
EOF

apt-get update
apt-get install -y --no-install-recommends ca-certificates \
    curl

# Install go
# https://golang.org/dl/go1.16.linux-amd64.tar.gz
curl https://golang.org/dl/go1.16.linux-amd64.tar.gz | tar -C /usr/local -xzf -

# Install deps to build.
mk-build-deps --install --remove \
  --tool='apt-get -o Debug::pkgProblemResolver=yes --no-install-recommends --yes' \
  "${BASE_DIR}/debian/control"

fakeroot debian/rules clean
fakeroot debian/rules build
fakeroot debian/rules binary
