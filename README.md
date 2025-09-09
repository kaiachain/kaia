# Kaia
Helloworld
Official golang implementation of the Kaia blockchain. Please visit our [Docs](https://docs.kaia.io/) for more details on Kaia design, node operation guides and application development resources.

## Building from Sources

Building the Kaia node binaries as well as utility tools, such as `kcn`, `kpn`, `ken`, `kbn`, `kscn`, `kspn`, `ksen`, `kgen`, `homi` and `abigen` requires
both a Go and a C compiler. You can install them using your favorite package manager. Once the dependencies are installed, run

    make all (or make {kcn, kpn, ken, kbn, kscn, kspn, ksen, kgen, homi, abigen})

## Executables

After successful build, executable binaries are installed at `build/bin/`.

| Command    | Description |
|:----------:|-------------|
| `kcn` | The CLI client for Kaia Consensus Node. Run `kcn --help` for command-line flags. |
| `kpn` | The CLI client for Kaia Proxy Node. Run `kpn --help` for command-line flags. |
| `ken` | The CLI client for Kaia Endpoint Node, which is the entry point into the Kaia network (main-, test- or private net).  It can be used by other processes as a gateway into the Kaia network via JSON RPC endpoints exposed on top of HTTP, WebSocket, gRPC, and/or IPC transports. Run `ken --help` for command-line flags. |
| `kscn` | The CLI client for Kaia ServiceChain Consensus Node.  Run `kscn --help` for command-line flags. |
| `kspn` | The CLI client for Kaia ServiceChain Proxy Node.  Run `kspn --help` for command-line flags. |
| `ksen` | The CLI client for Kaia ServiceChain Endpoint Node.  Run `ksen --help` for command-line flags. |
| `kbn` | The CLI client for Kaia Bootnode. Run `kbn --help` for command-line flags. |
| `kgen` | The CLI client for Kaia Nodekey Generation Tool. Run `kgen --help` for command-line flags. |
| `homi` | The CLI client for Kaia Helper Tool to generate initialization files. Run `homi --help` for command-line flags. |
| `abigen` | Source code generator to convert Kaia contract definitions into easy to use, compile-time type-safe Go packages. |

Both `kcn` and `ken` are capable of running as a full node (default) or an archive
node (retaining all historical state).

## Running a Core Cell

Core Cell (CC) is a set of one consensus node (CN) and one or more proxy nodes
(PNs). Core Cell plays a role of generating blocks in the Kaia network. We recommend to visit
the [CC Operation Guide](https://docs.kaia.io/docs/nodes/core-cell/)
for the details of CC bootstrapping process.

## Running an Endpoint Node

Endpoint Node (EN) is an entry point from the outside of the network in order to interact with the Kaia network. Currently, two official networks are available: Kairos and Mainnet (legacy names: Baobab and Cypress). Please visit the official
[EN Operation Guide](https://docs.kaia.io/docs/nodes/endpoint-node/).

## Running a Service Chain Node

Service chain node is a node for Service Chain which is an auxiliary blockchain independent from the main chain tailored for individual service requiring special node configurations, customized security levels, and scalability for high throughput services. Service Chain expands and augments Kaia by providing a data integrity mechanism and supporting token transfers between different chains.
Although the service chain feature is under development, we provide the operation guide for testing purposes. Please visit the official document [Service Chain Operation Guide](https://docs.kaia.io/docs/nodes/service-chain/).
Furthermore, for those who are interested in the Kaia Service Chain, please check out [Scaling Solutions](https://docs.kaia.io/docs/learn/scaling-solutions/).

## Import as library

Please use git branch (main) or commit hash. Do not use tags (v2.0.0)
```
go get github.com/kaiachain/kaia@main
go mod tidy
```

## License

The Kaia library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The Kaia binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.

## Contributing

As an open source project, Kaia always welcomes your contribution. Please read our [CONTRIBUTING.md](./CONTRIBUTING.md) for a walk-through of the contribution process.
