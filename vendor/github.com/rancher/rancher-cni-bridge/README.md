rancher-cni-bridge
========

This repo has the source code for the Rancher CNI plugin which is
based on the upstream CNI bridge plugin with customizations on top.

## Building

`make`


## Running

The `rancher-cni-bridge` plugin binary needs to be placed
inside `/opt/cni/bin` directory and the necessary configuration
file needs to be placed inside `/etc/cni/net.d` directory.


## License
Copyright (c) 2014-2016 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
