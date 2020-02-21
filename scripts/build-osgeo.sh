#!/usr/bin/env bash
##
# docker run -v $PWD:/app -it --env arch=cflinuxfs3 cloudfoundry/cflinuxfs3 /app/build.sh
##
set -eo pipefail

version="0.0.1"
sandbox=/tmp/sandbox
osgeolib=$sandbox/osgeolib
pythonsp=$osgeolib/python-osgeolib

apt-get update && apt-get install -y build-essential \
                   libsqlite3-dev \
                   sqlite3 \
                   python-dev

mkdir -p $osgeolib
cd $sandbox

# compile proj
wget http://download.osgeo.org/proj/proj-6.3.1.tar.gz
tar -xf proj-6.3.1.tar.gz && cd proj-6.3.1/
./configure --prefix=$osgeolib --enable-static=no
make install
cd $sandbox


# compile gdal
wget http://download.osgeo.org/gdal/3.0.4/gdal-3.0.4.tar.gz
tar xf gdal-3.0.4.tar.gz && cd gdal-3.0.4/
 ./configure --prefix=$osgeolib --enable-static=no --with-proj=$osgeolib \
    --with-libz=internal \
    --with-curl \
    --with-expat \
    --with-threads
    # --with-python
make install
cd $sandbox


# cleanup large boost headers
rm -fr $osgeolib/include/boost

# tar up directory
cd /tmp/sandbox/osgeolib
tar -czf /app/osgeolib-${version}-${arch}-linux-x64.tar.gz *