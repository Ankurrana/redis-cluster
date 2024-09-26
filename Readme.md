Redis Cluster On Local

# Steps
1. Install redis
2. Start atleast 3 redis nodes (with 0 replicas)
3. Start redis server using the commands 
    `redis-server redis-0.conf`
    `redis-server redis-1.conf`
    `redis-server redis-2.conf`
4. Form the cluster using the command
    `redis-cli --cluster create 127.0.0.1:6370 127.0.0.1:6371 127.0.0.1:6372 --cluster-replicas 0 --cluster-yes`
5. Run the sample go code which pushes student score to redis cluster and then gets the top 10 students


# Challenges
1. Quite tricky to start redis cluster using docker in a macbook. Will figure this out someday






