## WTF is this?

- Uses `scaleoutsean/trident-arm64:experimental` built from patched NetApp Trident code as of [Sep 26, 2022](https://github.com/NetApp/trident/tree/2451037d57cad6a5146fc75b615505daa1c1b57b)
  - Contains post-v22.07.0 commits which may be of interest to some
  - Introduces node selection (set to `amd64` by upstream, and changed to `arm64` in setup-experimental) and has some additional commits (see commit history at the link above)
  - Unlike `setup`, this configuration disables custom priorityClassName (`# priorityClassName: system-node-critical`) to make it easier to use. You may restore the setting in `setup-experimental/trident-daemonset.yaml` before deployment
- My Docker Hub image can be found [here](https://hub.docker.com/layers/scaleoutsean/trident-arm64/experimental/images/sha256-1c8b00f288a04d3bbb34f9b3375056fde3378416e11f7d3fbe11c380ff9dda5e?context=explore). As in main README.md, feel free to override with own image

## How to use it

Follow the main README.md, but use setup files from `setup-experimental` instead of `setup`.

```sh
mv setup setup-original
mv setup-experimental setup
./tridentctl -n trident install --use-custom-yaml 
```

Installation attempted and succeeded with K3s v1.24.4+k3s1 (ARM64).

