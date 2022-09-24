## Install w/o building it yourself

It's recommended to compare this repository vs. NetApp's and build from source code, but if you don't have the resources you are welcome to try my Docker Hub builds.

- Download `tridentctl` from releases
- Clone this repository

```sh
git clone https://github.com/scaleoutsean/trident -b v22.07.0-arm64

```

- Deploy (skip to Deploy with `tridentctl` below)


## Install required software to build

Install Docker-CE, go and the rest. See BUILD.md for more.

## Build on ARM64 system with Docker-CE

```sh
git clone https://github.com/scaleoutsean/trident -b v22.07.0-arm64
sudo GOOS=linux GOARCH=arm64 make trident_build
```

## View

```sh
sudo docker images | grep trident
```

Verify container has been built:

```sh
trident                    22.07.0-custom   298db8d80b51   31 seconds ago   163MB
```

Upload it to private registry or elsewhere and RTFM to find out how to refer to your image in YAML files.

See the Docker and other documentation for the details.

## Deploy with `tridentctl` 

You may use the `tridentctl` binary (for ARM64!) from Releases, you can build your own and extract `tridentctl` from the container image, or - when deploying from x86_64 client - you can use the official `tridentctl` (x86_64) to deploy to ARM64 systems.

As mentioned in Release Notes, it is recommended to use [custom deployment](https://docs.netapp.com/us-en/trident-2204/trident-get-started/kubernetes-customize-deploy-tridentctl.html) and remove ASUP (autosupport) from it. That's already been done for sample files in the `setup` directory.

If you wish to install to the namespace `trident`:

- Create a new namespace such as `trident`
- Edit two files in `setup` directory to change Trident image location (`image`) or use `--trident-image` to override the location from `tridentctl` (below)
  - setup/trident-daemonset.yaml
  - setup/trident-deployment.yaml
  - Image locations in daemonset and deployment YAML were changed from the NetApp (x86_64) default to `scaleoutsean/trident-arm64:v22.07.0` (Docker Hub) for people who don't want to build their own
  - Autosupport (ASUP) has been removed as mentioned earlier
- Run (sudo) `tridentctl install -n trident --use-custom-yaml` (--trident-image) to deploy Trident to the Trident namespace 
  - You may need to use sudo depending on your environment
  - If you hacked the YAML files or want to use my Docker Hub containers there's no need to specify custom image location with `--trident-image`

Users who use Helm, Trident Operator, etc. should check the official docs. I don't use that stuff.

