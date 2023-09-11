package main

import (
	"fmt"
	"sync"

	"github.com/IBM/sarama"
)

var (
	wg sync.WaitGroup
)

func main() {
	consumer, err := sarama.NewConsumer([]string{"62.65.17.194:9092"}, nil)

	if err != nil {
		panic(err)
	}

	partitionList, err := consumer.Partitions("vp_car_event")
	fmt.Println("partitionList", partitionList)
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("vp_car_event", int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		defer pc.AsyncClose()
		go func(sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, Key:%s, Value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	wg.Wait()
}
