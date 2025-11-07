package polity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const input_1 = `
-----BEGIN ORACLE PRIVATE KEY-----
addr: memnet://autumn-brook
nick: autumn-brook

RvsfZTHebtuEc7zi1mvT8cTaG1wczJg7akzqz+9pD2eoe4qZvYimlobJlKgLYp2B
VIcaopVUCDTAHXn0+RZQLy0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0t
LS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0=
-----END ORACLE PRIVATE KEY-----

-----BEGIN ORACLE PEER-----
nick: falling-dawn

pOCSkrZRwni5dyxWn1+puxPZBrRqtoyd+dwrRAn4ogmKiOPddAnxlf1S2y08ul1y
ymcJvx2UEhvzdIgBtA9vXA==
-----END ORACLE PEER-----
`

const input_2 = `
-----BEGIN ORACLE PRIVATE KEY-----
addr: memnet://autumn-brook
nick: autumn-brook

RvsfZTHebtuEc7zi1mvT8cTaG1wczJg7akzqz+9pD2eoe4qZvYimlobJlKgLYp2B
VIcaopVUCDTAHXn0+RZQLy0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0t
LS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0=
-----END ORACLE PRIVATE KEY-----
`

const input_3 = `
	i am not pem
`

const input_4 = `
-----BEGIN ORACLE PRIVATE KEY-----
addr: memnet://autumn-brook
nick: autumn-brook

RvsfZTHebtuEc7zi1mvT8cTaG1wczJg7akzqz+9pD2eoe4qZvYimlobJlKgLYp2B
VIcaopVUCDTAHXn0+RZQLy0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0t
LS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0=
-----END ORACLE PRIVATE KEY-----


i am not pem

`

func TestPemBag_Write_2_pems(t *testing.T) {
	pb := make(PemBag)
	n, err := pb.Write([]byte(input_1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == 0 {
		t.Fatalf("expected bytes to be written, got %d", n)
	}
	if got := len(pb); got != 2 {
		t.Fatalf("expected 2 entries in PemBag, got %d: %#v", got, pb)
	}
	if blocks, ok := pb["ORACLE PRIVATE KEY"]; !ok {
		t.Errorf("missing 'ORACLE PRIVATE KEY' entry")
	} else if len(blocks) != 1 {
		t.Errorf("expected 1 block for 'ORACLE PRIVATE KEY', got %d", len(blocks))
	}
	if blocks, ok := pb["ORACLE PEER"]; !ok {
		t.Errorf("missing 'ORACLE PEER' entry")
	} else if len(blocks) != 1 {
		t.Errorf("expected 1 block for 'ORACLE PEER', got %d", len(blocks))
	}
}

func TestPemBag_Write_not_pem(t *testing.T) {
	pb := make(PemBag)
	n, err := pb.Write([]byte(input_3))
	assert.Equal(t, 0, n)
	assert.NoError(t, err)
}

func TestPemBag_Write_1_pem(t *testing.T) {

	t.Run("exactcly one pem", func(t *testing.T) {
		pb := make(PemBag)
		n, err := pb.Write([]byte(input_2))
		assert.Equal(t, 1, pb.Size())
		assert.NoError(t, err)
		assert.Greater(t, n, 0)
	})

	t.Run("one pem and some garbage", func(t *testing.T) {
		pb := make(PemBag)
		n, err := pb.Write([]byte(input_4))
		assert.Equal(t, 1, pb.Size())
		assert.NoError(t, err)
		assert.Greater(t, n, 0)
	})

}
