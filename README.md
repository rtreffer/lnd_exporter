# `lnd_exporter` - lightning node metrics

`lnd_exporter` collects and exports metrics from a lightning node.

All metrics have a `nodekey` label with the pubkey current node.

The exporter is split into multiple collectors
| flag                           | default | description |
| ------------------------------ | ------- | ----------- |
| `--collector.basics`           | true    | `lnd_info`/`lnd_node_info` and generic sync information of the node |
| `--collector.network`          | true    | `lnd_graph_*` network based metrics |
| `--collector.channels`         | true    | `lnd_channel_*` metrics |
| `--collector.channels.closed`  | true    | `lnd_channel_*` metrics for closed channels |
| `--collector.channels.pending` | true    | `lnd_channel_*` metrics for pending channels |
| `--collector.peers`            | true    | `lnd_peer_*` metrics |
| `--collector.wallet`           | true    | `lnd_wallet_*` metrics |
| `--collector.peer.nodes`       | true    | `lnd_remote_node_*` metrics |
| `--collector.earnings`         | false   | `lnd_fwd_*` metrics |

## exporter metrics

| metric   | label   | meaning |
| -------- | ------- | ------- |
| `lnd_up` | `value` | `1` if lnd can be reached |

## collector.basics

This collector is based on `lncli getinfo`. It gives a quick overview
of the local node.

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_info` | `value` | `1` |
|            | `version` | lnd software version |
|            | `commit` | lnd software git commit |
| `lnd_node_info` | `value` | `1` |
|                 | `alias` | node alias |
|                 | `color` | node color |
|                 | `chain` | `bitcoin` |
|                 | `network` | `mainnet` / `testnet` |
| `lnd_chain_synced` | `value` | `1` if the blockchain is synced up |
| `lnd_graph_synced` | `value` | `1` if the lightning network graph is synced up | 
| `lnd_latest_block_time_seconds` | `value` | `unix time` of the last seen bitcoin block |
| `lnd_block_height_total` | `value` | the blockchain height |
| `lnd_peers_total`        | `value` | number of node peers |
| `lnd_channels_total`     | `value` | number of channels with a given set of properties |
|                  | `state` | `pending`, `active`, `inactive` |

## collector.network

This collector is based on `lncli getnetworkinfo`. It gives a quick overview
of the lightning network graph.

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_graph_nodes` | `value` | total number of lightning node |
| `lnd_graph_channels` | `value` | total number of lightning node |
| `lnd_graph_max_degree` | `value` | maximum nbumber of channels on one node |
| `lnd_graph_capacity_satoshis` | `value` | total number of satoshies in the lightning network |
| `lnd_graph_median_channel_size_satoshis` | `value` | median channel size in the lightning network |
| `lnd_graph_min_channel_size_satoshis` | `value` | minimum channel size in the lightning network |
| `lnd_graph_max_channel_size_satoshis` | `value` | maximum channel size in the lightning network |

## collector.channels

Thid collector is based on `lncli listchannels` and `lncli feereport`.
It gives an overview of all existing channels.

All metrics have a `channel_id` label to group channel related data.

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_channel_active` | `value` | `1` if the channel is active |
| `lnd_channel_info` | `value` | `1` |
| | `remotekey` | public key of the remote end |
| | `channel_point` | the channel point `txid:index` |
| | `public` | `true`/`false` if the channel is public/private |
| | `opened` | `local`/`remote` if the channel was opened by the local/remote node |
| | `state` | `open` |
| | `commit_type` | channel commitment type |
| | `close_address` | bitcoin address used on close |
| `lnd_channel_capacity_satoshis` | `value` | the channel capacity |
| `lnd_channel_local_balace_satoshis` | `value` | the local channel capacity |
| `lnd_channel_remote_balace_satoshis` | `value` | the remote channel capacity |
| `lnd_channel_unsettled_satoshis` | `value` | the sum of open HTLC transactions |
| | `direction` | `incoming`/`outgoing` |
| `lnd_channel_sent_satoshis_total` | `value` | the sum of sent satoshis |
| `lnd_channel_received_satoshis_total` | `value` | the sum of received satoshis |
| `lnd_channel_lifetime_seconds` | `value` | the lifetime of the channel |
| `lnd_channel_uptime_seconds` | `value` | the uptime of the channel |
| `lnd_channel_updates_total` | `value` | the number of channel updates |
| `lnd_channel_base_fee_satoshis` | `value` | the local base fee of the channel |
| `lnd_channel_fee_rate` | `value` | the fee rate of the channel as a ratio of sat per sat, `[0..1]` |

## collector.channels.closed

This collector is based on `lncli closedchannels`.

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_channel_info` | `value` | `1` |
| | `remotekey` | public key of the remote end |
| | `channel_id` | the channel id |
| | `channel_point` | the channel point `txid:index` |
| | `state` | `closed` |
| | `opened` | `local`/`remote` if the channel was opened by the local/remote node |
| | `closed` | `local`/`remote` if the channel was closed by the local/remote node |
| | `close_type` | `cooperative`/`local_force`/`remote_force`/`breach`/`cancelled`/`abandoned` |

**TODO**: this collector does not yet handle balances or resolutions.

## collector.channels.pending

This collector is based on `lncli pendingchannels`. Channels are keyed by `remotekey` and `channel_point`.

**NOTE**: the balance/capacity is namespaces as `lnd_pending_channel_*` due to a lack of `channel_id`.

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_channel_info` | `value` | `1` |
| | `remotekey` | public key of the remote end |
| | `channel_point` | the channel point `txid:index` |
| | `state` | `pending` |
| | `opened` | `local`/`remote` if the channel was opened by the local/remote node |
| | `commit_type` | channel commitment type |
| `lnd_pending_channel_capacity_satoshis` | `value` | the channel capacity |
| `lnd_pending_channel_local_balace_satoshis` | `value` | the local channel capacity |
| `lnd_pending_channel_remote_balace_satoshis` | `value` | the remote channel capacity |

## collector.peers

This collector is based on `lncli listpeers`.
All metrics are keyed by `remotekey`.

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_peer_info` | `value` | `1` |
| | `open` | `local`/`remote` depending on who opened the connection |
| | `graph_sync` | `unknown`/`active`/`passive`/`pinned` graph sync mode |
| `lnd_peer_address` | `value` | `1` |
| | `address` | public key of the remote end |
| `lnd_peer_sent_bytes_total` | `value` | Total bytes sent over the network to this peer |
| `lnd_peer_received_bytes_total` | `value` | Total bytes received over the network from this peer |
| `lnd_peer_sent_satoshis_total` | `value` | Total satoshis sent to this peer |
| `lnd_peer_received_satoshis_total` | `value` | Total satoshis received from this peer |
| `lnd_peer_reconnects_total` | `value` | Total number of reconnects to a peer |
| `lnd_peer_last_reconnect_time_seconds` | `value` | Unix time of the last reconnect |

## collector.peer.nodes

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_peer_node_info` | `value` | `1` |
| | `remotekey` | node alias |
| | `alias` | node alias |
| | `color` | node color |
| `lnd_peer_node_address` | `value` | `1` |
| | `address` | public key of the remote end |
| `lnd_channel_incoming_base_fee_satoshis` | `value` | the local base fee of the channel |
| `lnd_channel_incoming_fee_rate` | `value` | the fee rate of the channel as a ratio of sat per sat, `[0..1]` |

## collector.wallet

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_wallet_balance_satoshis` | `value` | the wallet balance |
| | `account` | `default` usually |
| | `confirmed` | `true`/`false` if the transaction is confirmed/unconfirmed |

## collector.earnings

**WARNING**: This collector fetches all historical forwardings. This will take a while on startup.

| metric     | label   | meaning |
| -----------| ------- | ------- |
| `lnd_fwd_fee_satoshis_total` | `value` | the total amount of fees earned |
| | `channel_id` | channel id |
| `lnd_fwd_sent_satoshis_total` | `value` | the total amount of satoshis sent |
| | `channel_id` | channel id |
| `lnd_fwd_received_satoshis_total` | `value` | the total amount of satoshis received |
| | `channel_id` | channel id |

# Under investigation

The following command reveal interesting informaiton but might be too expensive to add

- `lncli listinvoices`: list of invoices. Patches welcome.
- `lncli describegraph`: full graph, prometheus can handle the cardinality but export time is too long.
- `lncli getnodemetrics`: computes betweenness metric on a all nodes. Export time is too high.