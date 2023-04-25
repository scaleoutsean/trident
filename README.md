## Official ARM64 builds

Starting with v23.04, NetApp Trident supports ARM64, so please use upstream repository for v23.04 and future releases.

Below you can find instructions for v23.01 and some earlier releases.

## Install w/o building it yourself

It's recommended to first compare this repository vs. NetApp's and build it from source, but if you don't have the resources you are welcome to try my Docker Hub builds.

tldr;

```sh
# clone the repo with custom YAML setup files
git clone https://github.com/scaleoutsean/trident -b v23.01-arm64
cd trident; mkdir bin # to save tridentctl
wget https://github.com/scaleoutsean/trident/releases/download/v23.01/tridentctl -O ./bin/tridentctl
chmod +x ./bin/tridentctl
# use scaleoutsean's container from Docker Hub
./bin/tridentctl install -n trident --use-custom-yaml --trident-image docker.io/scaleoutsean/trident-arm64:v23.01
# if you get a permission error, try to prefix that with "sudo "
```

## Build it yourself on ARM64 system with Docker-CE

See the official Trident documentation (including BUILD.md) for more. My source has ARM64 hard-coded so GOARCH=arm64 is probably not necessary if you build on ARM64.

```sh
git clone https://github.com/scaleoutsean/trident -b v23.01-arm64
cd trident
GOOS=linux GOARCH=arm64 make trident_build
# prefix with "sudo " if necessary
```

See BUILD.md for other options. Then view your work:

```sh
docker images
ls ./bin
# should have trident_operator in addition to tridentctl downloaded earlier
```

Verify that Trident container has been built - you should see your Trident build, golang and an ARM64 base container:

```sh
$ sudo docker images
REPOSITORY                   TAG               IMAGE ID       CREATED          SIZE
trident-arm64                v23.01            11e85254d47a   47 minutes ago   180MB
trident                      23.01.0-custom    11e85254d47a   47 minutes ago   180MB
gcr.io/distroless/static     latest-arm64      69ffe45fb9e6   11 hours ago     2.34MB
```

Upload Trident container to private registry or elsewhere and RTFM to find out how to refer to your image in YAML files. 

See the Docker and other documentation for the details.

## Deploying Trident with `tridentctl` 

Create a namespace for Trident:

```sj
kubectl create namespace trident
```

There are several ways to do it:

- use the `tridentctl` binary (for ARM64!) from this repo's Releases, which is the method for all-ARM64 clusters (explained at the top)
- build your own container and extract `tridentctl` from the container image, then use tridentctl with self-built or my image (at the top)
- when deploying from x86_64 client: use the official `tridentctl` (for x86_64) to deploy this Trident build to ARM64 systems. This is probably the easiest choice for those who don't want to build from source. You'd need to download the official Trident Installer for x86_64 to get tridentctl (x86_64) or (if you want to build everything from the source) build Trident twice (once for this ARM64 source, to create an ARM64 container, and once for AMD64, using the NetApp source for tridentctl (x86_64))
  - For mixed clusters there's another variant, where Master nodes may be x86_64 and Worker nodes are ARM64. If you want to deploy for this situation, it's probaly easiest to deploy twice - first to x86_64 nodes using the official image, and another time using ARM64 image to deploy to Workers. Files in `setup` are meant for ARM64 so I don't think you should attempt to use them for x86_64 deployments

As mentioned in the official Release Notes, [custom deployment](https://docs.netapp.com/us-en/trident/trident-get-started/kubernetes-customize-deploy-tridentctl.html) lets you customize Trident installation - for example, ASUP (autosupport) can be stripped and that's already been done for sample files in `setup` because there's no autosupport image for ARM64.

To use own custom files, generate them and then edit them by yourself:

```sh
./tridentctl install --generate-custom-yaml
```

### Details on installing Trident v23.01 (ARM64) with `tridentctl`

If you wish to install to the namespace `trident`:

- Create the namespace `trident`
- If you want to build from the source or private registry, use `--trident-image` to override the image location. You can also hard-code it into:
  - setup/trident-daemonset.yaml
  - setup/trident-deployment.yaml
  - Image locations in daemonset and deployment YAML were changed from the NetApp Trident (x86_64) defaults to `scaleoutsean/trident-arm64:v23.01` (Docker Hub) for people who don't want to build their own or don't want to RTFM
  - Autosupport (ASUP) was removed as mentioned earlier
  - trident-operator string is replaced with trident-
- Run `tridentctl install -n trident --use-custom-yaml` to deploy Trident to the Trident namespace. Add `sudo` in front and `--trident-image ${LOCATION}` at the end if you need that. If you're using my Docker Hub image, those will install by default. Otherwise try `tridentctl install --help` to see how to use own images (local or from private container registry).

```sh
$ sudo kubectl get deploy -n trident
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
trident-controller   1/1     1            1           55s

$ sudo kubectl get pods -n trident
NAME                                  READY   STATUS    RESTARTS   AGE
trident-node-linux-hvpbt              1/2     Running   0          57s
trident-controller-766f68f4d7-f6gtr   5/5     Running   0          57s


$ sudo ./tridentctl version -n trident
+------------------------+------------------------+
|     SERVER VERSION     |     CLIENT VERSION     |
+------------------------+------------------------+
| 23.01.0-custom+unknown | 23.01.0-custom+unknown |
+------------------------+------------------------+

$ sudo ./tridentctl version -n trident --client
+------------------------+
|     CLIENT VERSION     |
+------------------------+
| 23.01.0-custom+unknown |
+------------------------+
```

Users who use Helm, Trident Operator, etc. should check the official docs. I don't use that stuff.

Version v23.01 contains some Windows stuff. I haven't tried to use that as I don't have Windows on ARM64, so unless you want to try to install Trident on Windows ARM64 clients or something like that, better use the official repo for Windows-related experimentation.

## Next steps

Just follow the official docs at this point - configure your client, pick a protocol/back-end, etc.

In my case I have a Debian client and a SolidFire VM running on ESXi v7:

```sh
$ cat /etc/debian_version 
11.6

$ uname -a
Linux k1 6.1.7-meson64 #22.11.4 SMP PREEMPT Mon Jan 23 21:25:00 UTC 2023 aarch64 GNU/Linux

$ sudo kubectl get nodes
NAME   STATUS   ROLES                  AGE    VERSION
k1     Ready    control-plane,master   129d   v1.24.4+k3s1
```

Modified ./trident-installer/sample-input/backends-samples/solidfire/backend-solidfire.json configured for a SolidFire storage account `k3s`:

```json
{
    "version": 1,
    "storageDriverName": "solidfire-san",
    "Endpoint": "https://admin:****@192.168.105.32/json-rpc/11.0",
    "SVIP": "192.168.1.32:3260",
    "TenantName": "k3s",
    "Types": [{"Type": "Bronze", "Qos": {"minIOPS": 100, "maxIOPS": 200, "burstIOPS": 400}},
              {"Type": "Silver", "Qos": {"minIOPS": 400, "maxIOPS": 600, "burstIOPS": 800}},
              {"Type": "Gold", "Qos": {"minIOPS": 600, "maxIOPS": 800, "burstIOPS": 1000}}]
}
```

Create a back-end in the Trident namespace.

```sh
$ ./bin/tridentctl -n trident create backend -f trident-installer/sample-input/backends-samples/solidfire/backend-solidfire.json
+------------------------+----------------+--------------------------------------+--------+---------+
|          NAME          | STORAGE DRIVER |                 UUID                 | STATE  | VOLUMES |
+------------------------+----------------+--------------------------------------+--------+---------+
| solidfire_192.168.1.32 | solidfire-san  | 62eda961-bdac-4298-a2c6-27282e627427 | online |       0 |
+------------------------+----------------+--------------------------------------+--------+---------+
```
