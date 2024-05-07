package ssh

import "testing"

func TestMd5File(t *testing.T) {
	file, err := Md5File("/opt/work/go/server/server.tgz")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(file)
}
