package main

import (
	"context"
	"encoding/json"
	"log"
	"os/signal"
	"strconv"
	"syscall"

	"fmt"
	"os"
	"time"

	"encoding/base64"
	mv2ray "gin/miniv2ray"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

var (
	MAINVER = "0.0.0-src"
)

const (
	accessKey    = ""
	secretKey    = ""
	region       = "ap-southeast-1"
	instanceName = "CentOS-1-V2ray"
)

func main() {
	r := gin.Default()
	r.LoadHTMLFiles("index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{})
	})

	r.GET("/vmess_ping", func(c *gin.Context) {

		verbose := false
		showNode := false
		usemux := false
		desturl := "http://www.google.com/gen_204"
		count := uint(3)
		timeout := uint(3)
		inteval := uint(1)
		quit := uint(0)

		vmess := "vmess://" + c.Query("vmess")
		fmt.Printf("vmess : %s/\n", vmess)
		// vmess := "vmess://eyJhZGQiOiIxOC4xNDEuMTEuMTk4IiwicGF0aCI6IiIsInBzIjoiYXdzLWt0Y3AiLCJwb3J0IjoiMzMwNiIsInYiOiIyIiwiaG9zdCI6IiIsInRscyI6IiIsImlkIjoiNWQ0ODkzYTAtMThkNS0xMWViLWE1MDEtMDI5NDA1YmI5MjBlIiwibmV0Ijoia2NwIiwidHlwZSI6Im5vbmUiLCJhaWQiOiIyIiwic25pIjoiIn0="

		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)

		ps, err := Ping(vmess, count, desturl, timeout, inteval, quit, osSignals, showNode, verbose, usemux)

		// vmessping.PrintVersion(MAINVER)
		// ps, err := vmessping.Ping(vmess, count, desturl, timeout, inteval, quit, osSignals, showNode, verbose, usemux)
		if err != nil {
			os.Exit(1)
		}
		ps.PrintStats()
		// if ps.IsErr() {
		// os.Exit(1)
		// }

		c.JSON(200, gin.H{"vmess": vmess, "counter": strconv.Itoa(int(ps.ReqCounter)), "success": strconv.Itoa(len(ps.Delays))})
	})

	r.GET("/instance", func(c *gin.Context) {

		creds := credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		}
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region), config.WithCredentialsProvider(creds))
		if err != nil {
			log.Fatalf("unable to load SDK config, %v", err)
		}
		// Using the Config value, create the DynamoDB client
		svc := lightsail.NewFromConfig(cfg)
		resp, err := svc.GetInstance(context.TODO(), &lightsail.GetInstanceInput{
			InstanceName: aws.String(instanceName),
		})

		if err != nil {
			log.Fatalf("failed to list tables, %v", err)
		}
		instance := resp.Instance

		fmt.Println(*instance.Name)
		// for _, tableName := range resp.TableNames {
		// 	fmt.Println(tableName)
		// }

		vmess, err := json.Marshal(map[string]string{
			"v":    "2",
			"ps":   "aws-tcp",
			"add":  *instance.PublicIpAddress,
			"port": "3306",
			"id":   "5d4893a0-18d5-11eb-a501-029405bb920e",
			"aid":  "0",
			"net":  "tcp",
			"type": "none",
			"host": "",
			"path": "",
			"tls":  "",
		})

		result := gin.H{"Arn": *instance.Arn, "BlueprintId": *instance.BlueprintId,
			"BlueprintName":   *instance.BlueprintName,
			"PublicIpAddress": *instance.PublicIpAddress,
			"State":           *instance.State.Name, "instanceName": *instance.Name, "vmess": base64.StdEncoding.EncodeToString(vmess)}

		c.JSON(200, result)
	})

	r.Use(favicon.New("./favicon.ico"))
	r.Run(":8080")
}

func PrintVersion(mv string) {
	fmt.Fprintf(os.Stderr,
		"VMessPing ver[%s], A prober for v2ray (v2ray-core: %s)\n", mv, mv2ray.CoreVersion())
}

type PingStat struct {
	StartTime  time.Time
	SumMs      uint
	MaxMs      uint
	MinMs      uint
	AvgMs      uint
	Delays     []int64
	ReqCounter uint
	ErrCounter uint
}

func (p *PingStat) CalStats() {
	for _, v := range p.Delays {
		p.SumMs += uint(v)
		if p.MaxMs == 0 || p.MinMs == 0 {
			p.MaxMs = uint(v)
			p.MinMs = uint(v)
		}
		if uv := uint(v); uv > p.MaxMs {
			p.MaxMs = uv
		}
		if uv := uint(v); uv < p.MinMs {
			p.MinMs = uv
		}
	}
	if len(p.Delays) > 0 {
		p.AvgMs = uint(float64(p.SumMs) / float64(len(p.Delays)))
	}
}

func (p PingStat) PrintStats() {
	fmt.Println("\n--- vmess ping statistics ---")
	fmt.Printf("%d requests made, %d success, total time %v\n", p.ReqCounter, len(p.Delays), time.Since(p.StartTime))
	fmt.Printf("rtt min/avg/max = %d/%d/%d ms\n", p.MinMs, p.AvgMs, p.MaxMs)
}

func (p PingStat) IsErr() bool {
	return len(p.Delays) == 0
}

func Ping(vmess string, count uint, dest string, timeoutsec, inteval, quit uint, stopCh <-chan os.Signal, showNode, verbose, usemux bool) (*PingStat, error) {
	server, err := mv2ray.StartV2Ray(vmess, verbose, usemux)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if err := server.Start(); err != nil {
		fmt.Println("Failed to start", err)
		return nil, err
	}
	defer server.Close()

	if showNode {
		go func() {
			info, err := mv2ray.GetNodeInfo(server, time.Second*10)
			if err != nil {
				return
			}

			fmt.Printf("Node Outbound: %s/%s\n", info["loc"], info["ip"])
		}()
	}

	ps := &PingStat{}
	ps.StartTime = time.Now()
	round := count
L:
	for round > 0 {
		seq := count - round + 1
		ps.ReqCounter++

		chDelay := make(chan int64)
		go func() {
			delay, err := mv2ray.MeasureDelay(server, time.Second*time.Duration(timeoutsec), dest)
			if err != nil {
				ps.ErrCounter++
				fmt.Printf("Ping %s: seq=%d err %v\n", dest, seq, err)
			}
			chDelay <- delay
		}()

		select {
		case delay := <-chDelay:
			if delay > 0 {
				ps.Delays = append(ps.Delays, delay)
				fmt.Printf("Ping %s: seq=%d time=%d ms\n", dest, seq, delay)
			}
		case <-stopCh:
			break L
		}

		if quit > 0 && ps.ErrCounter >= quit {
			break
		}

		if round--; round > 0 {
			select {
			case <-time.After(time.Second * time.Duration(inteval)):
				continue
			case <-stopCh:
				break L
			}
		}
	}

	ps.CalStats()
	return ps, nil
}
