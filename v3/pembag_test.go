package polity

import "testing"

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

func TestPemBag_Write(t *testing.T) {
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
