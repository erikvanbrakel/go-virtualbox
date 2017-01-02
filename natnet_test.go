package virtualbox

import "testing"

func TestNATNets(t *testing.T) {
	m, err := NATNets()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", m)
}

func TestCreateAndDeleteNATNet(t *testing.T) {
    nat, err := CreateNATNet("testing123", "10.0.123.0/24", false, true)
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("%+v", nat)

    nat.DHCPEnabled = true

    nat.Update()

    t.Logf("%+v", nat)

    err = DeleteNATNet("testing123")
    if err != nil {
        t.Fatal(err)
    }

    t.Logf("%+v", nat)
}
