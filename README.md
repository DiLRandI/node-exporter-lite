# node-exporter-lite

[![Go](https://github.com/DiLRandI/node-exporter-lite/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/DiLRandI/node-exporter-lite/actions/workflows/go.yml)
[![CodeQL](https://github.com/DiLRandI/node-exporter-lite/actions/workflows/codeql.yml/badge.svg?branch=main)](https://github.com/DiLRandI/node-exporter-lite/actions/workflows/codeql.yml)

A lightweight version of the [node-exporter](https://github.com/prometheus/node_exporter), specifically optimized for Raspberry Pi and Orange Pi devices. It aims to provide a simple and efficient method to monitor system resources on these low-power platforms. Unlike the original node-exporter, this lite version focuses on essential collectors, ensuring minimal resource usage while still offering valuable insights into system performance.

## Install

To install you can run the following command. Replace `${{env.RELEASE_VERSION}}` with the version you want to install.

```bash
curl -sSfL https://raw.githubusercontent.com/DiLRandI/node-exporter-lite/main/scripts/installer.sh | sh -s ${{env.RELEASE_VERSION}}
```

## Building

To build application, go 1.22 is required.

```bash
make build
```

## Running tests

 ```bash
make test
 ```
