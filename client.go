package sshkeymanager

import (
	"golang.org/x/crypto/ssh"
)

type IClient struct {
	Cl *ssh.Client
	Ses *ssh.Session
	User string
	Host string
	Port string
}

func (c *IClient) NewConnection(user string, host string, port string)  error {
	c.User = user
	c.Host = host
	c.Port = port
	var err error
	c.Cl, err = ConfigSSH(user, host, port)
	if err != nil {
		return err
	}
	return nil
}

func (c *IClient) NewSession()  error {
	var err error
	c.Ses, err = c.Cl.NewSession()
	if err != nil {
		return err
	}
	return nil
}

func (c *IClient) CloseConnection() error {
	err := c.Cl.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *IClient) CloseSession() error {
	err := c.Ses.Close()
	if err != nil {
		return err
	}
	return nil
}