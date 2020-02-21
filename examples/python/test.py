import time

import sys
from osgeo import gdal

version_num = int(gdal.VersionInfo('VERSION_NUM'))
if version_num < 1100000:
   print ("version: %s" %(version_num))

time.sleep(5000)
