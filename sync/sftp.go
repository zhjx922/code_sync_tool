package sync

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type SFTP struct {
	Session    *ssh.Session
	client     *ssh.Client
	sftpClient *sftp.Client
	dirCache   map[string]bool
}

func CreateSFTP(user, password, host string, port int) (*SFTP, error) {
	s := &SFTP{
		dirCache: make(map[string]bool),
	}

	if err := s.connect(user, password, host, port); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SFTP) wait(user, password, host string, port int) {
	err := s.client.Wait()
	if err != nil {
		log.Print("连接异常，开始重新连接(", err, ")")
		if err := s.connect(user, password, host, port); err != nil {
			panic(err)
		}
	}
}

// Connect 连接
func (s *SFTP) connect(user, password, host string, port int) error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		err          error
	)

	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr = fmt.Sprintf("%s:%d", host, port)
	retry := 999999
	for i := 0; i < retry; i++ {
		if s.client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
			if i == (retry - 1) {
				return err
			} else {
				log.Println("3秒后重新尝试连接")
				time.Sleep(3 * time.Second)
			}
		} else {
			break
		}
	}

	// 异常监控
	go s.wait(user, password, host, port)

	// create session
	if s.Session, err = s.client.NewSession(); err != nil {
		return err
	}

	s.sftpClient, err = sftp.NewClient(s.client)
	if err != nil {
		return err
	}

	return nil
}

func (s *SFTP) Upload(localFile string, deploymentFile string) {
	log.Println("upload", localFile)

	//打开本地文件流
	srcFile, err := os.Open(localFile)
	if err != nil {
		//fmt.Println("os.Open error : ", localFile)
		log.Println("sftp->open:", err)
		return
	}
	defer srcFile.Close()

	// 判断当前目录是否存在
	i := strings.LastIndex(deploymentFile, "/")
	deploymentPath := deploymentFile[:i]
	// 缓存目录状态
	if _, ok := s.dirCache[deploymentPath]; !ok {
		log.Println("dp:", deploymentPath)
		if _, err := s.sftpClient.Stat(deploymentPath); err != nil {
			s.sftpClient.Mkdir(deploymentPath)
		}
		s.dirCache[deploymentPath] = true
	}

	//打开远程文件,如果不存在就创建一个
	dstFile, err := s.sftpClient.Create(deploymentFile)
	if err != nil {
		log.Println("sftp->create:", deploymentFile, err)
		return
	}
	//关闭远程文件
	defer dstFile.Close()

	//本地文件流拷贝到上传文件流
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Println("sftp->copy:", err.Error())
	}
}

func (s *SFTP) Download(localFile string, deploymentFile string) {

}

func (s *SFTP) Delete(localFile string, deploymentFile string) {
	// log.Println("删除远程文件：", deploymentFile)
}
