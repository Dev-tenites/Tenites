#!/bin/bash
set -e

VALIDATOR_NUM=${1:-1}
MONIKER=${2:-"tenites-validator-$VALIDATOR_NUM"}
CHAIN_ID="tenites-testnet-1"
BINARY_URL="https://github.com/Tenites-Inc/Tenites-Connecting-Africa/releases/download/v1.0.0/tenitesd"

echo "Setting up Tenites validator $VALIDATOR_NUM: $MONIKER"

useradd -m -s /bin/bash tenites 2>/dev/null || true

EXPECTED_SHA256="REPLACE_WITH_RELEASE_SHA256"

wget -O /usr/local/bin/tenitesd $BINARY_URL
ACTUAL_SHA256=$(sha256sum /usr/local/bin/tenitesd | awk '{print $1}')
if [ "$EXPECTED_SHA256" != "REPLACE_WITH_RELEASE_SHA256" ] && [ "$ACTUAL_SHA256" != "$EXPECTED_SHA256" ]; then
    echo "FATAL: Binary checksum mismatch!"
    echo "Expected: $EXPECTED_SHA256"
    echo "Got:      $ACTUAL_SHA256"
    rm -f /usr/local/bin/tenitesd
    exit 1
fi
chmod +x /usr/local/bin/tenitesd

sudo -u tenites tenitesd init $MONIKER \
  --chain-id $CHAIN_ID \
  --home /home/tenites/.tenitesd

GENESIS_URL="https://raw.githubusercontent.com/Tenites-Inc/Tenites-Connecting-Africa/main/tenites-chain/testnet/genesis.json"
wget -O /home/tenites/.tenitesd/config/genesis.json $GENESIS_URL
chown tenites:tenites /home/tenites/.tenitesd/config/genesis.json

cp "$(dirname "$0")/tenitesd.service" /etc/systemd/system/
systemctl daemon-reload
systemctl enable tenitesd

echo "Validator $VALIDATOR_NUM setup complete"
echo "Next: Add validator key and configure peers"
echo "Then: systemctl start tenitesd"
