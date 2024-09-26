package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"sort"

	"github.com/redis/go-redis/v9"
)

type Student struct {
	Id    int
	Name  string
	Total int
}

func main() {
	addrs := []string{"0.0.0.0:6370", "0.0.0.0:6371", "0.0.0.0:6372"}
	// pusher := New("0.0.0.0:6379", "")
	pusher := New(addrs, "")
	studentScores := map[string]Student{}
	students := []string{}
	for i := 0; i < 1000; i++ {
		students = append(students, randomString())
	}

	// fmt.Println(students)

	for range 100000 {
		randStudent := rand.Intn(len(students))
		name := students[randStudent]

		score := rand.Intn(101)
		var s Student
		s, ok := studentScores[name]
		if !ok {
			s = Student{randStudent, name, 0}
		}
		s.Total += score
		studentScores[s.Name] = s

		pusher.Push(studentScores[s.Name])
	}

	sort.Slice(students, func(i, j int) bool {
		return studentScores[students[i]].Total > studentScores[students[j]].Total
	})

	redisResult := pusher.Top()

	for i := 0; i < 10; i++ {
		log.Println(studentScores[students[i]], redisResult[i])
	}

	//

}

func (st Student) String() string {
	return fmt.Sprintf("name:%s,score:%d", st.Name, st.Total)
}

func randomString() string {
	min := int('a')
	max := int('z')

	name := bytes.Buffer{}

	for i := 0; i < 10; i++ {
		randChar := min + rand.Intn(max-min+1)
		name.Write([]byte{byte(randChar)})
	}
	return name.String()
}

type ScorePublisher struct {
	rClient *redis.ClusterClient
	ctx     context.Context
}

// func New(addr, password string) *ScorePublisher {
func New(addrs []string, password string) *ScorePublisher {
	// rClient := redis.NewClient(&redis.Options{
	// 	Addr: addr,
	// })
	rClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})

	err := rClient.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("Error pinging Redis: %v", err)
	} else {
		log.Println("Successfully connected to Redis!")
	}

	v := &ScorePublisher{rClient: rClient, ctx: context.Background()}
	v.DelSets()
	return v
}

func (p *ScorePublisher) DelSets() error {
	var err error
	for i := 0; i < 10; i++ {
		_, err = p.rClient.Del(p.ctx, fmt.Sprintf("sorted_list_%d", i)).Result()
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (p *ScorePublisher) Push(s Student) error {

	// fmt.Printf("pushing result for %s\n", s)

	sortedSetName := fmt.Sprintf("sorted_list_%d", s.Id%10)

	_, err := p.rClient.ZRem(p.ctx, sortedSetName, s.Name).Result()

	if err != nil {
		fmt.Printf("unable to remove key due to %s\n", err)
	}

	_, err = p.rClient.ZAdd(p.ctx, sortedSetName, redis.Z{Member: s.Name, Score: float64(s.Total)}).Result()
	if err != nil {
		fmt.Printf("unable to push due to %s\n", err)
		return err
	}

	_, err = p.rClient.ZRemRangeByRank(p.ctx, sortedSetName, 0, -11).Result()

	if err != nil {
		fmt.Printf("unable to truncate due to %s\n", err)
		return err
	}
	return nil
}

func (p *ScorePublisher) Top() []Student {
	students := []Student{}
	for i := 0; i < 10; i++ {
		students = append(students, p.TopFromSet(fmt.Sprintf("sorted_list_%d", i))...)
	}

	sort.Slice(students, func(i, j int) bool {
		return students[i].Total > students[j].Total
	})
	students = students[0:10]

	return students

}

func (p *ScorePublisher) TopFromSet(set string) []Student {
	top10, err := p.rClient.ZRevRangeWithScores(p.ctx, set, 0, 9).Result()
	if err != nil {
		log.Fatalf("Error fetching top 10 entries: %v", err)
	}

	// Print the top 10 members
	fmt.Println("Top 10 leaderboard entries:")
	students := []Student{}
	for _, z := range top10 {
		students = append(students, Student{-1, z.Member.(string), int(z.Score)})
		fmt.Printf("Member: %s, Score: %f\n", z.Member, z.Score)
	}
	return students
}
