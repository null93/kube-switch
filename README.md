# Kube-Switch
> Switch between Kubernetes context & namespace using an interactive menu.

![MIT License](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)
![Version 1.0.1](https://img.shields.io/badge/Version-1.0.1-yellow.svg?style=for-the-badge)
![MacOS](https://img.shields.io/badge/OS-MacOS-yellow.svg?style=for-the-badge)

### About

Similar to [kubectx](https://github.com/ahmetb/kubectx) & [kubens](https://github.com/ahmetb/kubectx), but with the ability to pick using arrow keys.

### Install (Homebrew)

```shell
$ brew install null93/kube-switch/kube-switch
$ bind -x '"\C-k":"kube-switch"'
```

### Package

```shell
$ make package
```
