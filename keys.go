package sshkeymanager

import (
	"fmt"
	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"log"
	"os"
	"path"
	"strings"
)

type SSHKey struct {
	Num   int
	Key   string
	Email string
}

var sshKeys []SSHKey

func GetKeys(uid string, rootUser string, host string, port string) []SSHKey {

	usrs := GetUsers(rootUser, host, port)

	for _, u := range usrs {
		if u.UID == uid {
			client := ConfigSSH(rootUser, host, port)
			defer client.Close()
			session, err := client.NewSession()
			if err != nil {
				log.Fatal("Unable to create session ", err)
			}
			defer session.Close()
			raw, err := session.CombinedOutput("cat " + u.Home + "/.ssh/authorized_keys")
			if err != nil {
				log.Fatal("Unable to run command ", err)
			}
			rawToString := string(raw)

			strs := strings.Split(rawToString, "\n")
			for i, s := range strs {
				k := strings.Fields(s)
				if len(k) > 1 {
					var sshKey SSHKey
					sshKey.Num = i + 1
					sshKey.Key = k[0] + " " + k[1]
					if len(k) > 2 {
						sshKey.Email = k[2]
					}
					sshKeys = append(sshKeys, sshKey)
				}
			}

		}
	}
	return sshKeys
}

func DeleteKey(key string, uid string, rootUser string, host string, port string) {
	var (
		newKeys []SSHKey
		newKey SSHKey
	)

	keys := GetKeys(uid, rootUser, host, port)

	fields := strings.Fields(key)
	newKey.Key = fields[0] + " " + fields[1]
	if len(fields) > 2 {
		newKey.Email = fields[2]
	}
	for _, k := range keys {
		if k.Key != newKey.Key {
			newKeys = append(newKeys, k)
		}
	}
	sync(newKeys)
}

func AddKey(key string, uid string, rootUser string, host string, port string) {

	var k SSHKey

	keys := GetKeys(uid, rootUser, host, port)
	fields := strings.Fields(key)
	k.Num = len(keys) + 1
	k.Key = fields[0] + " " + fields[1]
	if len(fields) > 2 {
		k.Email = fields[2]
	}

	keys = append(keys, k)
	sync(keys)
}

func sync(keys []SSHKey, rootUser string, host string, port string)  {
	f, err := os.Create("authorized_keys")
	if err != nil {
		log.Fatal("Cannot create file ", err)
		f.Close()
		return
	}

	for _, k := range keys {
		fmt.Fprintln(f, k.Key + " " + k.Email)
	}
	err = f.Close()
	if err != nil {
		log.Fatal("Cannot write to file", err)
		return
	}

	//Using SCP for copy authorized_keys to server (maybe replace in future)
	clientConfig, _ := auth.PrivateKey(rootUser, path.Join(Home, ".ssh/id_rsa"), HostKeyCallback)

	client := scp.NewClient(host + ":" + port, &clientConfig)

	errConn := client.Connect()
	if errConn != nil {
		log.Fatal("Couldn't establish a connection to the remote server ", err)
		return
	}

	f, errFile := os.Open("authorized_keys")
	if errFile != nil {
		log.Fatal("Couldn't open file ", errFile)
	}

	defer client.Close()

	defer f.Close()

	//TODO remote path and set owner of file
	//
	err = client.CopyFile(f, "remote/user/path", "0600")

	if err != nil {
		log.Fatal("Error while copying file ", err)
	}


}