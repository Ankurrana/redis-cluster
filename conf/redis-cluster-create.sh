# sleep 3
# node_0_ip=$(getent hosts redis-cluster-node-0 | awk '{ print $1 }')
# node_1_ip=$(getent hosts redis-cluster-node-1 | awk '{ print $1 }')
# node_2_ip=$(getent hosts redis-cluster-node-2 | awk '{ print $1 }')


# redis-cli --cluster create $node_0_ip:6370 $node_1_ip:6371 $node_2_ip:6372 --cluster-replicas 0 --cluster-yes
redis-cli --cluster create 127.0.0.1:6370 127.0.0.1:6371 127.0.0.1:6372 --cluster-replicas 0 --cluster-yes
