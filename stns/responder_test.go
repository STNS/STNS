package stns

import "testing"

func TestUserBuildResource(t *testing.T) {
	g := v3Users{}
	attr := Attribute{
		ID: 1,
	}

	r := g.buildResource("test", &attr)

	switch v := r.(type) {
	case *v3User:
		if v.ID != 1 {
			t.Errorf("ummatch id got %d", v.ID)
		}
	default:
		t.Error("ummatch type")
	}

	attr = Attribute{
		ID:   1,
		User: &User{GroupID: 1},
	}

	r = g.buildResource("test", &attr)

	switch v := r.(type) {
	case *v3User:
		if v.ID != 1 {
			t.Errorf("ummatch id got %d", v.ID)
		}
		if v.GroupID != 1 {
			t.Errorf("ummatch id got %d", v.GroupID)
		}
	default:
		t.Error("ummatch type")
	}
}

func TestGroupBuildResource(t *testing.T) {
	g := v3Groups{}
	attr := Attribute{
		ID: 1,
	}

	r := g.buildResource("test", &attr)

	switch v := r.(type) {
	case *v3Group:
		if v.ID != 1 {
			t.Errorf("ummatch id got %d", v.ID)
		}

	default:
		t.Error("ummatch type")
	}

	attr = Attribute{
		ID:    1,
		Group: &Group{Users: []string{"test"}},
	}

	r = g.buildResource("test", &attr)
	switch v := r.(type) {
	case *v3Group:
		if v.ID != 1 {
			t.Errorf("ummatch id got %d", v.ID)
		}
		if v.Users[0] != "test" {
			t.Errorf("ummatch users got %s", v.Users[0])
		}
	default:
		t.Error("ummatch type")
	}
}

func TestSudoBuildResource(t *testing.T) {
	s := v3Sudoers{}
	attr := Attribute{
		User: &User{Password: "password"},
	}

	r := s.buildResource("test", &attr)

	switch v := r.(type) {
	case *v3Sudo:
		if v.Password != "password" {
			t.Errorf("ummatch pasword got %s", v.Password)
		}
	default:
		t.Error("ummatch type")
	}
}
