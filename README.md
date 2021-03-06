## OSGEO Buildpack

this buildpack can be used as an additional buildpack in a multi-buildpack scenario to install the osgeo packages into the droplet. This will install the low level OSGEO packages but not any language specific bindings, those can be installed with a lanuguage specific buildpack. see the python example below. this is a minimal buildpack and should be used in conjunction with another buildpack. most of the heavy lifting is done in the [build script](/scripts/osgeo-build.sh). 

inspiration taken from this [buildpack](https://github.com/planetfederal/osgeolib-buildpack)

### Included Packages 

#### GDAL 2.2.1
* proj-6.3.1
* gdal-2.2.1

#### GDAL 3.0.4
* proj-6.3.1
* gdal-2.2.1

to switch between versions use the `OSGEO_VERSION` env variable

### Added ENV Vars

the following vars are added to the runtime environment of the droplet

`GDAL_DATA`
`PROJ_LIB`
`LDFLAGS`
`CPLUS_INCLUDE_PATH`
`C_INCLUDE_PATH`


### Usage

this buildpack will work with both `cflinuxfs2` or `cflinuxfs3`. you will need to use a different build for each though.you can also specify the version of the compiled osgeo libs using an environment var `OSGEO_VERSION`. see the releases page for details

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
  #Optional version can be specified
  env:
    OSGEO_VERSION: 3.0.4

```

4. push your application, depending on if your buildpack is cached or not it may go to the internet to download the vendored osgeo libs.

### Python offline Example

assuming the cached version of the buildpack has been uploaded to the foundation here is a quick python example.

1. cd into the `examples/python` dir in this repo

`cd examples/python`

2. vendor the dependencies

```
pip download -r requirements.txt --no-binary=:none: -d vendor
```

3. push the app to your foundation. this will trigger the use of the osgeo and python buildpacks. by having the `gdal` pip packaged vendored this will prevent it from going to the internet.

```
cf push -f manifest.yml
```


4. you will see the gdal buildpack run and then python. when python runs it will install the gdal python bindings. it may take a minute since it is compiling against the installed gdal packages.

5. your app should be running now.




### Building the vendored osgeo libs
**only needed to be done when updating versions of osgeo libs**

becuase osgeo needs to be installed in an a user space and/or offline inside of the droplet we need to vendor the packages.

1. execute the build script from the root of this repo

```bash
docker run -v ${PWD}/scripts:/app -it --env arch=cflinuxfs3 cloudfoundry/cflinuxfs3 /app/build-osgeo.sh
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
