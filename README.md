# Toy Overlay

A Kubernetes-first Pod overlay network with a UDP backend, perfect for debugging.

## Docs
For a full explanation of what this is, how it works and how it was was built, you can find full documentation of this project at:
[samxsmith.com/toyoverlay](samxsmith.com/toyoverlay)

## To Use
- Use my docker image or build your own using the Makefile
- Set the image name and tag in the daemon-set.yml
- Copy the yml file to the master node
- On the master node run: `kubectl create -f daemon-set.yml`

### Environment Variables
- HOSTNAME: name of the kubernetes host
