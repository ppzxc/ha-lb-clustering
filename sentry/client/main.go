package main

import (
	"flag"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	sys "gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func main() {
	dsn := flag.String("dsn", "", "sentry DSN")
	ip := flag.String("syslog.ip", "0.0.0.0", "syslog received server ip")
	port := flag.String("syslog.port", "5141", "syslog received server ip")
	flag.Parse()
	if err := raven.SetDSN(*dsn); err != nil {
		panic(err)
	}

	channel := make(sys.LogPartsChannel)
	handler := sys.NewChannelHandler(channel)

	server := sys.NewServer()
	server.SetFormat(sys.Automatic)
	server.SetHandler(handler)
	server.SetFormat(sys.Automatic)
	server.ListenUDP(*ip + ":" + *port)
	server.Boot()

	go func(channel sys.LogPartsChannel) {
		for logParts := range channel {
			// for key, value := range logParts {
			// 	fmt.Printf("key=%s, value=%s\r\n", key, value)
			// }
			// fmt.Println("===========================================")
			send(parse(logParts))
		}
	}(channel)

	server.Wait()
}

func send(syslog Syslog) {
	if syslog.Severity == 4 {
		packet := &raven.Packet{
			Message:    syslog.Content,
			Level:      raven.WARNING,
			ServerName: syslog.Hostname,
			Timestamp:  raven.Timestamp(syslog.Timestamp),
			Tags: raven.Tags{
				raven.Tag{
					Key:   "process.name",
					Value: syslog.Tag,
				},
				raven.Tag{
					Key:   "hostname",
					Value: syslog.Hostname,
				},
				raven.Tag{
					Key:   "client",
					Value: syslog.Client,
				},
			},
		}
		logrus.Warnf("%+#v\r\n", syslog)
		_, err := raven.Capture(packet, nil)
		select {
		case <-err:
			//if err != nil {
			//	logrus.Error(eid, e)
			//}
		}
	}
	if syslog.Severity >= 5 {
		packet := &raven.Packet{
			Message:    syslog.Content,
			Level:      raven.ERROR,
			ServerName: syslog.Hostname,
			Timestamp:  raven.Timestamp(syslog.Timestamp),
			Tags: raven.Tags{
				raven.Tag{
					Key:   "process.name",
					Value: syslog.Tag,
				},
				raven.Tag{
					Key:   "hostname",
					Value: syslog.Hostname,
				},
				raven.Tag{
					Key:   "client",
					Value: syslog.Client,
				},
			},
		}
		logrus.Errorf("%+#v\r\n", syslog)
		_, err := raven.Capture(packet, nil)
		select {
		case <-err:
			//if err != nil {
				//logrus.Error(eid, e)
			//
		}
	}
}

type Syslog struct {
	Severity  int
	Priority  int
	Facility  int
	Hostname  string
	Timestamp time.Time
	Client    string
	Content   string
	Tls_peer  string
	Tag       string
}

func parse(logParts format.LogParts) (syslog Syslog) {
	if severity, ok := logParts["severity"]; ok && severity != nil {
		syslog.Severity = severity.(int)
	}
	if priority, ok := logParts["priority"]; ok && priority != nil {
		syslog.Priority = priority.(int)
	}
	if facility, ok := logParts["facility"]; ok && facility != nil {
		syslog.Facility = facility.(int)
	}
	if hostname, ok := logParts["hostname"]; ok && hostname != nil {
		syslog.Hostname = hostname.(string)
	}
	if timestamp, ok := logParts["timestamp"]; ok && timestamp != nil {
		syslog.Timestamp = timestamp.(time.Time)
		// loc, _ := time.LoadLocation("Asia/Seoul")
		// syslog.Timestamp = syslog.Timestamp.In(loc)
	}
	if client, ok := logParts["client"]; ok && client != nil {
		syslog.Client = client.(string)
	}
	if content, ok := logParts["content"]; ok && content != nil {
		syslog.Content = content.(string)
	}
	if tls_peer, ok := logParts["tls_peer"]; ok && tls_peer != nil {
		syslog.Tls_peer = tls_peer.(string)
	}
	if tag, ok := logParts["tag"]; ok && tag != nil {
		syslog.Tag = tag.(string)
	}
	return
}
