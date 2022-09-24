## Install required software

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

You may use the binary (it's for ARM64!) from Releases or - on x86_64 - use the official `tridentctl`.

As mentioned in Release Notes, it is recommended to use [custom deployment](https://docs.netapp.com/us-en/trident-2204/trident-get-started/kubernetes-customize-deploy-tridentctl.html) and remove ASUP (autosupport) from it.

If you wish to install to the namespace `trident`, create the namespace and use YAML files from setup directory:

- Go to `setup` directory and change Trident image location in lines with `image` or override them `--trident-image` in `tridentctl` (below)
- From parent directory run (sudo) `tridentctl install -n trident --use-custom-yaml` (--trident-image) to deploy. You may need to use sudo depending on your environment


