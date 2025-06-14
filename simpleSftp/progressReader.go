package simpleSftp

import "io"

type progressReader struct {
	io.Reader
	Reporter func(r int64)
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	if pr.Reporter != nil {
		pr.Reporter(int64(n))
	}
	return
}
