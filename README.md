## OSGEO Buildpack

this buildpack can be used as an additional buildpack in a multi-buildpack scenario to install the osgeo packages into the droplet. This will install the low level OSGEO packages but not any language specific bindings, those can be installed with a lanuguage specific buildpack. see the python example below. this is a minimal buildpack and should be used in conjunction with another buildpack. most of the heavy lifting is done in the [build script](/scripts/osgeo-build.sh). 

inspiration taken from this [buildpack](https://github.com/planetfederal/osgeolib-buildpack)

### Included Packages 

* proj-6.3.1
* gdal-3.0.4

### Added ENV Vars

the following vars are added to the runtime environment of the droplet

`GDAL_DATA`
`PROJ_LIB`
`LDFLAGS`
`CPLUS_INCLUDE_PATH`
`C_INCLUDE_PATH`


### Usage

this buildpack will work with both `cflinuxfs2` or `cflinuxfs3`. you will need to use a different build for each though.

the latest releases can be downloaded from the releases page in this repo. there are cached and unchached versions.

1. download the latest release of your choice for `cflinuxfs3` or `cflinuxfs2`
2. upload the buildpack to your foundation
   ```bash
    cf create-buildpack [BUILDPACK_NAME] [BUILDPACK_ZIP_FILE_PATH] [POSITION]
   ```
3. create a manifest file to use the buildpack with another builpack

```yml
applications:
- name: test-gdal
  buildpacks:
  - osgeo_buildpack
  - python_buildpack 
```

4. push your application, depending on if your buildpack is cached or not it may go to the internet to download the vendored osgeo libs.

### Building the vendored osgeo libs
**only needed to be done when updating versions of osgeo libs**

becuase osgeo needs to be installed in an a user space and/or offline inside of the droplet we need to vendor the packages.

1. execute the build script from the root of this repo

```bash
docker run -v scripts:/app -it --env arch=cflinuxfs3 cloudfoundry/cflinuxfs3 /app/build-osgeo.sh
```
2. this will output a `tar.gz` that will then be uploaded to github releases
3. update the dependencies in the `manifest.yml` to reflect the new version



### Building the Buildpack

To build this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

    ```bash
    source .envrc
    ```

    To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Install buildpack-packager

    ```bash
    ./scripts/install_tools.sh
    ```

1. Build the buildpack

    ```bash
    buildpack-packager build
    ```

1. Use in Cloud Foundry

    Upload the buildpack to your Cloud Foundry and optionally specify it by name

    ```bash
    cf create-buildpack [BUILDPACK_NAME] [BUILDPACK_ZIP_FILE_PATH] 1
    cf push my_app [-b BUILDPACK_NAME]
    ```

### Testing

Buildpacks use the [Cutlass](https://github.com/cloudfoundry/libbuildpack/cutlass) framework for running integration tests.

To test this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

    ```bash
    source .envrc
    ```

    To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Run unit tests

    ```bash
    ./scripts/unit.sh
    ```

1. Run integration tests

    ```bash
    ./scripts/integration.sh
    ```

    More information can be found on Github [cutlass](https://github.com/cloudfoundry/libbuildpack/cutlass).

### Reporting Issues

Open an issue on this project

## Disclaimer

This buildpack is experimental and has not been heavily tested in production. 
