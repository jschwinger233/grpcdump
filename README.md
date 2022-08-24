# grpcdump

A grpcdump that really works.

## Installation

For Linux user, please download from [releaess](https://github.com/jschwinger233/grpcdump/releases).

## Requirements

[libpcap](https://www.tcpdump.org/) is required for Linux, for Ubuntu:

```bash
apt install libpcap-dev
```

## Usage

Let's grpcdump [Etcd](https://etcd.io/)!

### 1. Preparation

I'm using Etcd v3.4.9, so I've prepared the proto files as follows:

```
.
|-- rpc.proto
|-- etcd
|   |-- auth
|   |   `-- authpb
|   |       `-- auth.proto
|   `-- mvcc
|       `-- mvccpb
|           `-- kv.proto
|-- gogoproto
|   `-- gogo.proto
`-- google
    `-- api
        |-- annotations.proto
        `-- http.proto
```

- `rpc.proto` is copied from `etcdserver/etcdserverpb/rpc.proto`, where the `Watch` service is defined;
- the others are dependency of `rpc.proto`, you can see them from `import` definitions;
- `auth.proto` and `kv.proto` can be found in the etcd repo, and the rest protos must be secured by downloading online: [gogo.proto](https://raw.githubusercontent.com/gogo/protobuf/master/gogoproto/gogo.proto), [annotations.proto](https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto), [http.proto](https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto), you are welcome;

### 2. Sniffing

I already have a local etcd service, so let's sniff its traffic!

```
$ grpcdump -i lo -p 2379 -f rpc.proto
Jan 15 22:19:59.336063	127.0.0.1:42452->127.0.0.1:2379	packetno:12	streamid:1	header:map[:authority:127.0.0.1:2379 :method:POST :path:/etcdserverpb.KV/Range :scheme:http content-type:application/grpc te:trailers user-agent:grpc-go/1.7.5]
Jan 15 22:19:59.336063	127.0.0.1:42452->127.0.0.1:2379	packetno:12	streamid:1	data:key:"/calico/resources/v3/projectcalico.org/clusterinformations/default"
Jan 15 22:19:59.336312	127.0.0.1:2379->127.0.0.1:42452	packetno:20	streamid:1	header:map[:status:200 content-type:application/grpc]
Jan 15 22:19:59.336312	127.0.0.1:2379->127.0.0.1:42452	packetno:20	streamid:1	data:header:<cluster_id:14841639068965178418 member_id:10276657743932975437 revision:108218 raft_term:203> kvs:<key:"/calico/resources/v3/projectcalico.org/clusterinformations/default" create_revision:270 mod_revision:274 version:2 value:"{\"kind\":\"ClusterInformation\",\"apiVersion\":\"projectcalico.org/v3\",\"metadata\":{\"name\":\"default\",\"uid\":\"91b3d4cf-96bd-11eb-a00f-cc483a63a267\",\"creationTimestamp\":\"2021-04-06T09:50:50Z\"},\"spec\":{\"clusterGUID\":\"0a35a7ca25b04e45863fbe4bbdc1d34b\",\"calicoVersion\":\"v3.4.4-2-g1f083c2\",\"datastoreReady\":true}}"> count:1
Jan 15 22:19:59.336312	127.0.0.1:2379->127.0.0.1:42452	packetno:20	streamid:1	header:map[grpc-message: grpc-status:0]
Jan 15 22:19:59.336703	127.0.0.1:42452->127.0.0.1:2379	packetno:26	streamid:3	header:map[:authority:127.0.0.1:2379 :method:POST :path:/etcdserverpb.KV/Range :scheme:http content-type:application/grpc te:trailers user-agent:grpc-go/1.7.5]
Jan 15 22:19:59.336703	127.0.0.1:42452->127.0.0.1:2379	packetno:26	streamid:3	data:key:"/calico/resources/v3/projectcalico.org/felixconfigurations/default"
Jan 15 22:19:59.337918	127.0.0.1:2379->127.0.0.1:42452	packetno:32	streamid:3	header:map[:status:200 content-type:application/grpc]
Jan 15 22:19:59.337918	127.0.0.1:2379->127.0.0.1:42452	packetno:32	streamid:3	data:header:<cluster_id:14841639068965178418 member_id:10276657743932975437 revision:108218 raft_term:203> kvs:<key:"/calico/resources/v3/projectcalico.org/felixconfigurations/default" create_revision:275 mod_revision:275 version:1 value:"{\"kind\":\"FelixConfiguration\",\"apiVersion\":\"projectcalico.org/v3\",\"metadata\":{\"name\":\"default\",\"uid\":\"91b4e67e-96bd-11eb-a00f-cc483a63a267\",\"creationTimestamp\":\"2021-04-06T09:50:50Z\"},\"spec\":{\"logSeverityScreen\":\"Info\",\"reportingInterval\":\"0s\"}}"> count:1
Jan 15 22:19:59.337918	127.0.0.1:2379->127.0.0.1:42452	packetno:32	streamid:3	header:map[grpc-message: grpc-status:0]
Jan 15 22:19:59.338955	127.0.0.1:42452->127.0.0.1:2379	packetno:38	streamid:5	header:map[:authority:127.0.0.1:2379 :method:POST :path:/etcdserverpb.KV/Range :scheme:http content-type:application/grpc te:trailers user-agent:grpc-go/1.7.5]
Jan 15 22:19:59.338955	127.0.0.1:42452->127.0.0.1:2379	packetno:38	streamid:5	data:key:"/calico/resources/v3/projectcalico.org/felixconfigurations/node.gray-latitude-5410"
Jan 15 22:19:59.339154	127.0.0.1:2379->127.0.0.1:42452	packetno:44	streamid:5	header:map[:status:200 content-type:application/grpc]
Jan 15 22:19:59.339154	127.0.0.1:2379->127.0.0.1:42452	packetno:44	streamid:5	data:header:<cluster_id:14841639068965178418 member_id:10276657743932975437 revision:108218 raft_term:203> kvs:<key:"/calico/resources/v3/projectcalico.org/felixconfigurations/node.gray-latitude-5410" create_revision:276 mod_revision:276 version:1 value:"{\"kind\":\"FelixConfiguration\",\"apiVersion\":\"projectcalico.org/v3\",\"metadata\":{\"name\":\"node.gray-latitude-5410\",\"uid\":\"91b51ba2-96bd-11eb-a00f-cc483a63a267\",\"creationTimestamp\":\"2021-04-06T09:50:50Z\"},\"spec\":{\"defaultEndpointToHostAction\":\"Return\"}}"> count:1
Jan 15 22:19:59.339154	127.0.0.1:2379->127.0.0.1:42452	packetno:44	streamid:5	header:map[grpc-message: grpc-status:0]
^C
```

The output covers capture time, connection info, packet number, stream id, and payload of header or data.

We can change the output format to JSON, which allows us to filter the frames much more easier:

```
$ grpcdump -i lo -p 2379 -f rpc.proto -o json | jq
{
  "time": "2022-01-15T22:26:39.832096442+08:00",
  "packet_number": 118,
  "src": "127.0.0.1",
  "dst": "127.0.0.1",
  "sport": 2379,
  "dport": 44082,
  "stream_id": 7,
  "type": "Header",
  "payload": {
    ":status": "200",
    "content-type": "application/grpc"
  },
  "ext": {}
}
{
  "time": "2022-01-15T22:26:39.832096442+08:00",
  "packet_number": 118,
  "src": "127.0.0.1",
  "dst": "127.0.0.1",
  "sport": 2379,
  "dport": 44082,
  "stream_id": 7,
  "type": "Data",
  "payload": {
    "count": "1",
    "header": {
      "clusterId": "14841639068965178418",
      "memberId": "10276657743932975437",
      "raftTerm": "203",
      "revision": "108218"
    },
    "kvs": [
      {
        "createRevision": "271",
        "key": "L2NhbGljby9yZXNvdXJjZXMvdjMvcHJvamVjdGNhbGljby5vcmcvbm9kZXMvZ3JheS1sYXRpdHVkZS01NDEw",
        "modRevision": "103996",
        "value": "eyJraW5kIjoiTm9kZSIsImFwaVZlcnNpb24iOiJwcm9qZWN0Y2FsaWNvLm9yZy92MyIsIm1ldGFkYXRhIjp7Im5hbWUiOiJncmF5LWxhdGl0dWRlLTU0MTAiLCJ1aWQiOiI5MWI0MWJiYi05NmJkLTExZWItYTAwZi1jYzQ4M2E2M2EyNjciLCJjcmVhdGlvblRpbWVzdGFtcCI6IjIwMjEtMDQtMDZUMDk6NTA6NTBaIn0sInNwZWMiOnsiYmdwIjp7ImlwdjRBZGRyZXNzIjoiMTAuMjIuNzEuMTg5LzIxIn19fQ==",
        "version": "185"
      }
    ]
  },
  "ext": {
    "data_direction": "service_to_client",
    "data_path": "/etcdserverpb.KV/Range"
  }
}
{
  "time": "2022-01-15T22:26:39.832096442+08:00",
  "packet_number": 118,
  "src": "127.0.0.1",
  "dst": "127.0.0.1",
  "sport": 2379,
  "dport": 44082,
  "stream_id": 7,
  "type": "Header",
  "payload": {
    "grpc-message": "",
    "grpc-status": "0"
  },
  "ext": {}
}
^C
```

JSON output makes convenience to filter data as long as you are familar with [jq](https://stedolan.github.io/jq/), such as `grpcdump -i lo -p 2379 -f rpc.proto -o json | jq '. | select(.dport==44082)'` to pick out the frames from etcd to client with port 44082.

### 3. Pcap Parsing

The ability to parse gRPC from pcap is also important!

Let's dump a pcap:

```
$ tcpdump -i lo port 2379 -w etcd.pcap
tcpdump: listening on lo, link-type EN10MB (Ethernet), capture size 262144 bytes
^C391 packets captured
782 packets received by filter
0 packets dropped by kernel
```

Then use `-r etcd.pcap` instead of previous `-i lo`:

```
$ grpcdump -p 2379 -r etcd.pcap -f rpc.proto
```

Parsing from pcap is important, because the grpcdump doesn't need to implement BPF filter (see `man 7 pcap-filter`) like `gateway snup and ip[2:2] > 576`; all you need to do, is to use `tcpdump(8)` to generate a pcap file, then the grpcdump will do the remaining.

### 4. Path Guessing

GRPC on HTTP2 has an amazing feature: hpack header compression, which also causes troubles because we can't always capture the complete traffic from the beginning of connection, leaving us unable to hpack-decode the headers thereafter.

Beside hpack issues, the missing request frames also lead to the similar consequences.

To demonstrate this easily, let's make an etcd watch:

```
etcdctl watch /zc --prefix
```

Then we inspect the client port of etcdctl:

```
$ lsof -p $(pidof etcdctl)
etcdctl 1803971 root    5u     IPv4 10672558      0t0      TCP localhost:48576->localhost:2379 (ESTABLISHED)
```

Then use grpcdump starts to sniff this connection:

```
grpcdump -i lo -p 2379 -f rpc.proto | grep 48576
```

Then we put a key to trigger data push:

```
etcdctl put /zc/a a
```

And this time the grpcdump gives us this:

```
$ s grpcdump -i lo -p 2379 -f rpc.proto | grep 48576
Jan 15 22:51:20.720648	127.0.0.1:2379->127.0.0.1:48576	packetno:638	streamid:1	data:(unknown)
```

The data frame fails to parse because the watch request has been sent ahead of the grpcdump's sniffing, and the missing request header makes it impossible to parse the data.

Unless we guess!

`grpcdump --guest-path/-m` is designed just for this scenario:

```
$ grpcdump -i lo -p 2379 -f rpc.proto -m /etcdserverpb.Watch/Watch | grep 48576
Jan 15 22:57:31.046106	127.0.0.1:2379->127.0.0.1:48576	packetno:379	streamid:1	data:(guess)header:<cluster_id:14841639068965178418 member_id:10276657743932975437 revision:108220 raft_term:203> events:<kv:<key:"/zc/a/" create_revision:108219 mod_revision:108220 version:2 value:"a">>
```

See!

If the data is parsed basing on guess, there is a `(guess)` indicator after data field.

And the grpcdump can even guess the missing `:path` automatically!

```
$ grpcdump -i lo -p 2379 -f rpc.proto -m AUTO -o json | jq '. | select(.dport==48576)'
{
  "time": "2022-01-15T23:03:47.356580804+08:00",
  "packet_number": 559,
  "src": "127.0.0.1",
  "dst": "127.0.0.1",
  "sport": 2379,
  "dport": 48576,
  "stream_id": 1,
  "type": "Data",
  "payload": {
    "events": [
      {
        "kv": {
          "createRevision": "108219",
          "key": "L3pjL2Ev",
          "modRevision": "108222",
          "value": "YQ==",
          "version": "4"
        }
      }
    ],
    "header": {
      "clusterId": "14841639068965178418",
      "memberId": "10276657743932975437",
      "raftTerm": "203",
      "revision": "108222"
    }
  },
  "ext": {
    "data_direction": "service_to_client",
    "data_guessed": "yes",
    "data_path": "/etcdserverpb.Watch/Watch"
  }
}
```

Awesome!

### 5. Copy As Grpcurl

One of my favorite feature in Chrome's developer tools is `Copy As cURL`, so I was wondering if I can do something similar to that.

`grpcdump --grpcurl` will add an additional [grpcurl](https://github.com/fullstorydev/grpcurl) command after every request frame:

```
$ grpcdump -i lo -p 2379 -f rpc.proto --grpcurl
Jan 15 23:09:35.449540	127.0.0.1:54944->127.0.0.1:2379	packetno:14	streamid:1	header:map[:authority:127.0.0.1:2379 :method:POST :path:/etcdserverpb.KV/Range :scheme:http content-type:application/grpc te:trailers user-agent:grpc-go/1.7.5]
Jan 15 23:09:35.449540	127.0.0.1:54944->127.0.0.1:2379	packetno:14	streamid:1	data:key:"/calico/resources/v3/projectcalico.org/clusterinformations/default"
grpcurl -plaintext -proto rpc.proto -d '{"key":"L2NhbGljby9yZXNvdXJjZXMvdjMvcHJvamVjdGNhbGljby5vcmcvY2x1c3RlcmluZm9ybWF0aW9ucy9kZWZhdWx0"}' 127.0.0.1:2379 etcdserverpb.KV/Range
^C
```

See the last line above, you can simply copy that command and run it from your terminal, to re-send an exactly same request.

This is extremely useful to reproduce issues and investigate the causes.
