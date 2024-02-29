package zookeeper

import (
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

type ZkManager struct {
	hosts []string `json:"hosts`
	conn *zk.Conn
	pathPrefix string
}

func (z *ZkManager) RegisterServerPath(nodePath, host string) error {
 ex, _, err := z.conn.Exists(nodePath)

 if err != nil {
	log.Println(err)
	return err
 }

 if !ex {
	_, err := z.conn.Create(nodePath, nil , 0, zk.WorldACL(zk.PermAll))

	if err != nil {
		log.Println(err)
		return err
	 }
 }

 subNodePath := nodePath + "/" + host
 ex, _, err = z.conn.Exists(subNodePath)
 if err != nil {
	log.Println(err)
	return err
 }

 if !ex {
	_, err := z.conn.Create(subNodePath, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Println(err)
		return err
	 }
 }
 return nil
}

func (z *ZkManager) GetConnection() error {
	conn, event, err := zk.Connect(z.hosts, 5 * time.Second)
	if err != nil {
		log.Println(err)
	}
	z.conn = conn
	e := <- event
	log.Printf("event: %+v", e)
	return err
}

func (z *ZkManager) GetServerListPath(path string) (list []string, err error) {
	list, _, err = z.conn.Children(path)
	return
}

func (z *ZkManager) WatchServerListByPath(path string) (chan []string ,chan error) {
	conn := z.conn
	snapshots := make(chan []string)
	errors := make(chan error)

	go func ()  {
		for {
			snapshot, _, events, err := conn.ChildrenW(path)

			if err != nil {
				errors <- err
			}

			snapshots <- snapshot
			evt := <-events
			if evt.Err != nil {
					errors <- evt.Err
			}

			log.Printf("ChildrenW Event: %v+", evt)
		}
	}()
	return snapshots, errors
}


func (z *ZkManager) Close()  {
	z.conn.Close()
}

func NewZkManager(hosts []string) *ZkManager {
	return &ZkManager{
		hosts: hosts,
		pathPrefix: "/gateway_servers_",
	}
}

