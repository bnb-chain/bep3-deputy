# Run a local bnb chain with docker

Change into this folder, then:

Initialize a new chain

    docker run -it --rm -v $(pwd):/root/bin -v $(pwd)/.bnbchaind:/root/.bnbchaind -v $(pwd)/.bnbcli:/root/.bnbcli kava/bnbchaintesting:0.6.3 bash /root/bin/start-new-bnbchain.sh

Start the local chain

    docker run -it --rm -v $(pwd)/.bnbchaind:/root/.bnbchaind -v $(pwd)/.bnbcli:/root/.bnbcli -p 26657:36657 kava/bnbchaintesting:0.6.3 bnbchaind start
