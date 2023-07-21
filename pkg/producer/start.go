package producer

import (
	"github.com/the-engineers-site/data-producer/pkg/handlers"
	"github.com/the-engineers-site/data-producer/pkg/logger"
	"github.com/the-engineers-site/data-producer/pkg/store"
	"go.uber.org/zap"
	"net"
	"strconv"
	"sync"
)

type Producer struct {
	Records    []string
	Connection *net.Conn
	EPS        int
}

func FormProducers(mapping map[string][]*handlers.Configuration) []*Producer {
	var producers []*Producer
	for host, conf := range mapping {
		producers = append(producers, GetProducer(host, conf)...)
	}
	return producers
}

func StartAll(producers []*Producer) {
	wg := sync.WaitGroup{}
	for _, p := range producers {
		wg.Add(1)
		go p.Start(&wg)
	}
	wg.Wait()
}
func (p *Producer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.GetLogger().Info("starting producer", zap.Int("eps", p.EPS))
	guard := make(chan struct{}, p.EPS)
	for {
		counter := 0
		for _, line := range p.Records {
			guard <- struct{}{}
			go func(line string) {
				store.Send(&counter, p.Connection, line)
				<-guard
			}(line)
		}
		logger.GetLogger().Info("Batch completed", zap.Int("count", counter))
	}
}

func GetProducer(host string, c []*handlers.Configuration) []*Producer {
	var producers []*Producer
	for _, conf := range c {
		conn, err := net.Dial("tcp", host)
		if err != nil {
			logger.GetLogger().Error("error while connecting to server", zap.Error(err), zap.String("Server", host))
			continue
		}
		eps, err := strconv.Atoi(conf.Eps)
		if err != nil {
			logger.GetLogger().Error("invalid eps", zap.Error(err))
			continue
		}
		producers = append(producers, &Producer{
			EPS:        eps,
			Records:    conf.Records,
			Connection: &conn,
		})
	}
	return producers
}
