**This repo is out of maintenance and decommissioned.**
# Deputy

**Deputy is responsible for handling swap requests from users automatically**. It will monitor swap requests from both 
chain and react to these requests. If swap request from user is legal, deputy will send corresponding HTLT tx to other chain.
When user claims tokens on other chain, deputy will claim tokens on the original chain.

Deputy will manage lifecycle of each swap request. It is supposed to handle normal cases or exceptional cases correctly 
and send corresponding txs timely util lifecycle of swap request is end. 

But there are still some cases that deputy can not handle:
+ Insufficient tokens: either gas token or swap token
+ Errors on blockchain nodes: you may not be able to get blocks or send txs to blockchain node.

You may need to resend tx when you have deposit enough tokens to deputy address or fix your node issue. You can refer to 
[admin doc](./spec/admin.md) for more detail.

## Installation

### Build

```bash
make build
```
### Config

There is a config template in `config` directory, you should create your own config to run your deputy correctly. 
You can refer to [config doc](./spec/config.md) for more details. **Please make sure you understand each config option 
in config file, and testing before open to public is recommended.**

### Start deputy

```bash
cd build
./deputy --bnb-network [0 for testnet, 1 for mainnet] --config-type [local or aws] --config-path config_file_path --aws-region [aws region or omit] --aws-secret-key [aws secret key for config or omit]
```

## Docker 

## Docker build

```bash
$ docker build --tag deputy .
```

## Docker run

For there is log file and database file(maybe if you are using sqlite), so you should mount a directory storing data files
into container, so even if container disappeared, you will not lose your data.

```bash
$ docker run -it -v /your/data/path:/deputy -e BNB_NETWORK={0 or 1} -e CONFIG_TYPE="local" -e CONFIG_FILE_PATH=/your/config/file/path/in/container -d deputy
```

## Contributing

You can refer to [Developer guide](./spec/develop_guide.md) for more details.

## License

Distributed under the [GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html). See [LICENSE](LICENSE) for more information.

## Links

+ [Deputy Overview](./spec/deputy.md)
+ [Executor](./spec/executor.md)
+ [Admin](./spec/admin.md)
+ [State machine](./spec/state_machine.md)
+ [Configuration](./spec/config.md)
