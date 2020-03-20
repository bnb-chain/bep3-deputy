FROM kava/bnbchain:0.6.3

# Install sponge and jq for convenience
RUN apt-get update
RUN apt-get install -y moreutils jq

# The node binaries are in /node-binary, ie /node-binary/fullnode/prod/0.6.3/linux/bnbchaind
# Add the latest version to the PATH for ease of use:
RUN ln -sf /node-binary/fullnode/prod/0.6.3/linux/bnbchaind /usr/local/bin/bnbchaind
RUN ln -sf /node-binary/cli/prod/0.6.3/linux/bnbcli /usr/local/bin/bnbcli
#chmod +x /usr/local/bin/bnbchaind ??

# These are the default data directories for the chain and cli.
# You can mount local directories to these locations if you want access to the data outside of the container.
# eg docker run -it --rm -v /Users/someone/.bnbchaind:/root/.bnbchaind -v /Users/someone/.bnbcli:/root/.bnbcli kava.bnbchaintesting:0.6.3 bnbchaind init
VOLUME "/root/.bnbchaind"
VOLUME "/root/.bnbcli"

# Remove binance's custom entry point
ENTRYPOINT []

# Set the default command
CMD bnbchaind

