package simpleSftp

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SimpleSFTP struct {
	*sftp.Client
	sshClient *ssh.Client
}

func ConnectWithPassword(username string, password string, host string, port string) (*SimpleSFTP, error) {
	host = fmt.Sprintf("%v:%v", host, port)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Replace with real host key validation in production
	}

	// Connect to SSH server
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial SSH: %v", err)
	}

	// Start SFTP session
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("Failed to start SFTP session: %v", err)
	}

	return &SimpleSFTP{
			Client:    client,
			sshClient: conn,
		},
		nil
}

func (ssftp *SimpleSFTP) Close() error {
	sshCloseErr := ssftp.sshClient.Close()
	sftpCloseErr := ssftp.Client.Close()
	if sshCloseErr != nil {
		return sshCloseErr
	}
	return sftpCloseErr
}

func (ssftp *SimpleSFTP) PutProgress(local string, remote string, progress func(currentFile string, currentBytes int, totalBytes int)) error {
	localInfo, err := os.Stat(local)
	if err != nil {
		return err
	}
	if localInfo.IsDir() {
		contents, err := os.ReadDir(local)
		if err != nil {
			return err
		}
		for _, dirContent := range contents {
			err = ssftp.PutProgress(path.Join(local, dirContent.Name()), path.Join(remote, dirContent.Name()), progress)

			if err != nil {
				return err
			}
		}
		return nil
	} else {
		localF, err := os.Open(local)
		if err != nil {
			return fmt.Errorf("problem opening local file %v: %v", local, err)
		}
		defer localF.Close()

		err = ssftp.MkdirAll(path.Dir(remote))
		if err != nil {
			return fmt.Errorf("problem creating remote file path %v: %v", path.Dir(remote), err)
		}

		remoteF, err := ssftp.Create(remote)
		if err != nil {
			return fmt.Errorf("problem creating remote file %v: %v", local, err)
		}
		defer remoteF.Close()

		localProgressReader := progressReader{
			Reader: localF,
			Reporter: func(r int64) {
				if progress != nil {
					progress(localInfo.Name(), int(r), int(localInfo.Size()))
				}
			},
		}

		_, err = io.Copy(remoteF, &localProgressReader)
		if err != nil {
			return fmt.Errorf("problem copying file %v to %v: %v", local, remote, err)
		}
		return nil
	}
}

func (ssftp *SimpleSFTP) Put(local string, remote string) error {
	return ssftp.PutProgress(local, remote, nil)
}
