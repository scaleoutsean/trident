## Install w/o building it yourself

It's recommended to first compare this repository vs. NetApp's and build it from source, but if you don't have the resources you are welcome to try my Docker Hub builds.

tldr;

```sh
git clone https://github.com/scaleoutsean/trident -b v22.10.0-arm64
cd trident; mkdir bin # to save tridentctl
wget https://github.com/scaleoutsean/trident/releases/download/v22.10.0-arm64/tridentctl -O ./bin/tridentctl
chmod +x ./bin/tridentctl
# use scaleoutsean's container from Docker Hub
./bin/tridentctl install -n trident --use-custom-yaml --trident-image docker.io/scaleoutsean/trident-arm64:v22.10.0
```

## Build it yourself on ARM64 system with Docker-CE

See the official Trident documentation (including BUILD.md) for more. My source has ARM64 hard-coded so GOARCH=arm64 is probably not necessary.

```sh
git clone https://github.com/scaleoutsean/trident -b v22.10.0-arm64
cd trident
sudo GOOS=linux GOARCH=arm64 make trident_build
```

See BUILD.md for other options. Then view your work:

```sh
docker images
```

Verify that Trident container has been built - you should see your Trident build, golang and an ARM64 base container:

```sh
REPOSITORY                    TAG              IMAGE ID       CREATED          SIZE
trident                       22.10.0-custom   759ced7fb159   45 minutes ago   175MB
golang                        1.18             425d66fc5b6c   4 days ago       822MB
gcr.io/distroless/static      latest-arm64     1fa3b6b7eabc   52 years ago     2.34MB
```

Upload Trident container to private registry or elsewhere and RTFM to find out how to refer to your image in YAML files. 

See the Docker and other documentation for the details.

## Deploying Trident with `tridentctl` 

There are several ways to do it:

- use the `tridentctl` binary (for ARM64!) from Releases, which is the method for all-ARM64 clusters (explained at the top)
- build your own container and extract `tridentctl` from the container image, then use tridentctl with self-built or my image (at the top)
- when deploying from x86_64 client: use the official `tridentctl` (for x86_64) to deploy this Trident build to ARM64 systems. This is probably the easiest choice for those who don't want to build from source. You'd need to download the official Trident Installer for x86_64 to get tridentctl (x86_64) or (if you want to build everything from the source) build Trident twice (once for this ARM64 source, to create an ARM64 container, and once for AMD64, using the NetApp source for tridentctl (x86_64))
  - For mixed clusters there's another variant, where Master nodes may be x86_64 and Worker nodes are ARM64. If you want to deploy for this situation, it's probaly easiest to deploy twice - first to x86_64 nodes using the official image, and another time using ARM64 image to deploy to Workers. Files in `setup` are meant for ARM64 so I don't think you should attempt to use them for x86_64 deployments

As mentioned in Release Notes, [custom deployment](https://docs.netapp.com/us-en/trident-2204/trident-get-started/kubernetes-customize-deploy-tridentctl.html) lets you customize installation - for example, ASUP (autosupport) can be stripped. That's already been done for sample files in `setup` because there's no autosupport image for ARM64.

### Details on installing Trident v22.10.0 (ARM64) with `tridentctl`

If you wish to install to the namespace `trident`:

- Create the namespace `trident`
- If you want to build from the source or private registry, use `--trident-image` to override the image location. You can also hard-code it into:
  - setup/trident-daemonset.yaml
  - setup/trident-deployment.yaml
  - Image locations in daemonset and deployment YAML were changed from the NetApp Trident (x86_64) defaults to `scaleoutsean/trident-arm64:v22.10.0` (Docker Hub) for people who don't want to build their own or don't want to RTFM
  - Autosupport (ASUP) was removed as mentioned earlier
- Run `tridentctl install -n trident --use-custom-yaml` to deploy Trident to the Trident namespace. Add `sudo` in front and `--trident-image ${LOCATION}` at the end if you need that

Users who use Helm, Trident Operator, etc. should check the official docs. I don't use that stuff. This time (v22.10.0) I built an image and posted it [here](https://hub.docker.com/repository/docker/scaleoutsean/trident-operator-arm64). No idea if it works. 

