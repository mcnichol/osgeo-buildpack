#!/usr/bin/env bash
##
set -eo pipefail

version="0.0.1"
sandbox=/tmp/sandbox
osgeolib=$sandbox/osgeolib
pythonsp=$osgeolib/python-osgeolib

#need to install sqllite3 latest manually for fs2
if [[ ${arch} = "cflinuxfs2" ]]
then
    echo "running for cflinuxfs2"
    apt-get update && apt-get install -y build-essential python-dev
    apt-get remove sqlite3 libsqlite3-dev -y
    cd /tmp
    wget https://sqlite.org/2020/sqlite-autoconf-3310100.tar.gz
    tar -xvzf sqlite-autoconf-3310100.tar.gz

    cd /tmp/sqlite-autoconf-3310100
    ./configure --prefix=/usr --disable-static CFLAGS="-g"
    make
    make install
else
    echo "running for cflinuxfs3"
    apt-get update && apt-get install -y build-essential \
                   libsqlite3-dev \
                   sqlite3 \
                   python-dev
fi



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
    --with-png=internal \
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