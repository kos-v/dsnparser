package dsnparser

import "testing"

func TestDSN_HasParam(t *testing.T) {
	result := Parse("mysql://user:password@example.com:3306/dbname?foo=foo val")
	if !result.HasParam("foo") {
		t.Errorf("Unexpected value. Must be true")
	}
	if result.HasParam("bar") {
		t.Errorf("Unexpected value. Must be false")
	}
}

func TestParse_Scheme(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"", ""},
		{"mysql", ""},
		{"://", ""},
		{"mysql://", "mysql"},
		{"mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_", "mysql"},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetScheme() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetScheme())
		}
	}
}

func TestParse_User(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"", ""},
		{"@", ""},
		{"user", ""},
		{"user@", "user"},
		{"mysql://user", ""},
		{"mysql://user@", "user"},
		{"user:password", ""},
		{"user:password@", "user"},
		{"mysql://user:password", ""},
		{"mysql://user:password@", "user"},
		{"mysql://@", ""},
		{"mysql://:password@", ""},
		{"user:@", "user"},
		{":password@", ""},
		{"u@", "u"},
		{":@", ""},
		{"mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_", "user"},
		{"mysql://user@example.com:3306/dbname?tblsprefix=fs_", "user"},
		{"mysql://example.com:3306/dbname?tblsprefix=fs_", ""},
		{"\\@@", "@"},
		{"us\\@er@", "us@er"},
		{"mysql://use\\@r:password@", "use@r"},
		{"mysql://us\\@e\\:r@", "us@e:r"},
		{"mysql://user\\:password\\@example.com\\:3306/dbname?tblsprefix=fs_:password@", "user:password@example.com:3306/dbname?tblsprefix=fs_"},
		{"mysql://пользователь:пароль", ""},
		{"mysql://пользователь:пароль@", "пользователь"},
		{"mysql://пользователь1:пароль@", "пользователь1"},
		{"mysql://пользователь", ""},
		{"mysql://пользователь@", "пользователь"},
		{"mysql://пользователь1@", "пользователь1"},
		{"mysql://п#о\\@льзователь\\::\\@пар\\:оль@", "п#о@льзователь:"},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetUser() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetUser())
		}
	}
}

func TestParse_Password(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"", ""},
		{"@", ""},
		{"user", ""},
		{"user@", ""},
		{"mysql://user", ""},
		{"mysql://user@", ""},
		{"user:password", ""},
		{"user:password@", "password"},
		{"mysql://user:password", ""},
		{"mysql://user:password@", "password"},
		{"mysql://@", ""},
		{"mysql://:password@", "password"},
		{"user:@", ""},
		{":password@", "password"},
		{":p@", "p"},
		{":@", ""},
		{"mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_", "password"},
		{"mysql://user@example.com:3306/dbname?tblsprefix=fs_", ""},
		{"mysql://example.com:3306/dbname?tblsprefix=fs_", ""},
		{"\\@@", ""},
		{"us\\@er@", ""},
		{"mysql://user:p\\@ssw\\:ord@", "p@ssw:ord"},
		{"mysql://:p\\@ssw\\:ord@", "p@ssw:ord"},
		{"mysql://user:password\\@example.com\\:3306/dbname?tblsprefix=fs_@", "password@example.com:3306/dbname?tblsprefix=fs_"},
		{"mysql://пользователь:пароль", ""},
		{"mysql://пользователь:пароль1", ""},
		{"mysql://пользователь:пароль@", "пароль"},
		{"mysql://пользователь:пароль1@", "пароль1"},
		{"mysql://:пароль", ""},
		{"mysql://:пароль@", "пароль"},
		{"mysql://:1пароль_1@", "1пароль_1"},
		{"mysql://п#о\\@льзователь\\::\\@пар\\:оль@", "@пар:оль"},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetPassword() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetPassword())
		}
	}
}

func TestParse_Transport(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"example.com", ""},
		{"example.com:3306", ""},
		{"(example.com)", ""},
		{"tcp()", "tcp"},
		{"x(example.com)", "x"},
		{"tcp(example.com)", "tcp"},
		{"tcp(example.com:3306)", "tcp"},
		{"mysql://user:password@tcp(example.com)", "tcp"},
		{"mysql://user:password@tcp(example.com:3306)", "tcp"},
		{"mysql://user:password@tcp(example.com)/dbname?tblsprefix=fs_", "tcp"},
		{"mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_", "tcp"},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetTransport() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetTransport())
		}
	}
}

func TestParse_Host(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"mysql://user:password@example.com:3306/dbname?tblsprefix=fs_", "example.com"},
		{"mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_", "example.com"},
		{"mysql://user:password@example.com/dbname?tblsprefix=fs_", "example.com"},
		{"mysql://user:password@tcp(example.com)/dbname?tblsprefix=fs_", "example.com"},
		{"mysql://user:password@localhost/dbname?tblsprefix=fs_", "localhost"},
		{"mysql://user:password@127.0.0.1/dbname?tblsprefix=fs_", "127.0.0.1"},
		{"mysql://user:password@example.com", "example.com"},
		{"mysql://user:password@example.com:3306", "example.com"},
		{"mysql://user:password@/dbname?tblsprefix=fs_", ""},
		{"mysql://user:password@:3306", ""},
		{"mysql://user:password@", ""},
		{"mysql://example.loc", "example.loc"},
		{"example.loc", "example.loc"},
		{"example.loc:3306", "example.loc"},
		{"example.loc/path", "example.loc"},
		{"example.loc:3306/path", "example.loc"},
		{"example.loc:/", "example.loc"},
		{"example", "example"},
		{"example:3306", "example"},
		{"example/path", "example"},
		{"example:/", "example"},
		{"mysql://user:password@not$valid@host/dbname?tblsprefix=fs_", "not$valid@host"},
		{"mysql://user:password@not$valid@hostdbname?tblsprefix=fs_", "not$valid@hostdbname?tblsprefix=fs_"},
		{"mysql://user:password@хост.лок:3306/dbname?tblsprefix=fs_", "хост.лок"},
		{"mysql://user:password@хост.лок", "хост.лок"},
		{"mysql://user:password@хост.лок:3306/путь", "хост.лок"},
		{"mysql://user:password@хост.лок:3306", "хост.лок"},
		{"mysql://user:password@хост.лок1", "хост.лок1"},
		{"mysql://user:password@хост.ло1к", "хост.ло1к"},
		{"kafka://username:pasword@tcp(ip1:9093,ip2:9093,ip3:9093)/?topic=vsulblog", "ip1"},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetHost() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetHost())
		}
	}
}

func TestParse_HostPort(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"mysql://user:password@example.com:3306/dbname?tblsprefix=fs_", "example.com:3306"},
		{"mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_", "example.com:3306"},
		{"mysql://user:password@example.com/dbname?tblsprefix=fs_", "example.com"},
		{"mysql://user:password@tcp(example.com)/dbname?tblsprefix=fs_", "example.com"},
		{"mysql://user:password@localhost/dbname?tblsprefix=fs_", "localhost"},
		{"mysql://user:password@127.0.0.1/dbname?tblsprefix=fs_", "127.0.0.1"},
		{"mysql://user:password@example.com", "example.com"},
		{"mysql://user:password@example.com:3306", "example.com:3306"},
		{"mysql://user:password@/dbname?tblsprefix=fs_", ""},
		{"mysql://user:password@:3306", ":3306"},
		{"mysql://user:password@", ""},
		{"mysql://example.loc", "example.loc"},
		{"example.loc", "example.loc"},
		{"example.loc:3306", "example.loc:3306"},
		{"example.loc/path", "example.loc"},
		{"example.loc:3306/path", "example.loc:3306"},
		{"example.loc:/", "example.loc:"},
		{"example", "example"},
		{"example:3306", "example:3306"},
		{"example/path", "example"},
		{"example:/", "example:"},
		{"mysql://user:password@not$valid@host/dbname?tblsprefix=fs_", "not$valid@host"},
		{"mysql://user:password@not$valid@hostdbname?tblsprefix=fs_", "not$valid@hostdbname?tblsprefix=fs_"},
		{"mysql://user:password@хост.лок:3306/dbname?tblsprefix=fs_", "хост.лок:3306"},
		{"mysql://user:password@хост.лок", "хост.лок"},
		{"mysql://user:password@хост.лок:3306/путь", "хост.лок:3306"},
		{"mysql://user:password@хост.лок:3306", "хост.лок:3306"},
		{"mysql://user:password@хост.лок1", "хост.лок1"},
		{"mysql://user:password@хост.ло1к", "хост.ло1к"},
		{"kafka://username:pasword@tcp(ip1:9093,ip2:9093,ip3:9093)/?topic=vsulblog", "ip1:9093,ip2:9093,ip3:9093"},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetHostPort() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetHostPort())
		}
	}
}

func TestParse_Port(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"mysql://user:password@example.com:3306/dbname?tblsprefix=fs_", "3306"},
		{"mysql://user:password@tcp(example.com:3306)/dbname?tblsprefix=fs_", "3306"},
		{"mysql://user:password@example.com:3306", "3306"},
		{"mysql://user:password@tcp(example.com:3306)", "3306"},
		{"mysql://user:password@example.com/dbname?tblsprefix=fs_", ""},
		{"mysql://user:password@example.com:/dbname?tblsprefix=fs_", ""},
		{"mysql://user:password@example.com:bad, but working/dbname?tblsprefix=fs_", "bad, but working"},
		{"example.com:3306", "3306"},
		{"tcp(example.com:3306)", "3306"},
		{"example.com:", ""},
		{"example.com", ""},
		{"хост.лок:3306", "3306"},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetPort() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetPort())
		}
	}
}

func TestParse_Path(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"mysql://user:password@example.com:3306/foo?tblsprefix=fs_", "foo"},
		{"mysql://user:password@tcp(example.com:3306)/foo?tblsprefix=fs_", "foo"},
		{"mysql://user:password@example.com:3306/foo/bar/baz?tblsprefix=fs_", "foo/bar/baz"},
		{"mysql://user:password@example.com:3306//?tblsprefix=fs_", "/"},
		{"mysql://user:password@example.com:3306/?tblsprefix=fs_", ""},
		{"mysql://user:password@example.com:3306/foo", "foo"},
		{"mysql://user:password@example.com:3306/foo/bar/baz", "foo/bar/baz"},
		{"mysql://user:password@example.com:3306", ""},
		{"mysql://user:password@example.com:3306/", ""},
		{"mysql://user:password@example.com:3306//", "/"},
		{"example.com/foo?tblsprefix=fs_", "foo"},
		{"example.com/foo/bar/baz?tblsprefix=fs_", "foo/bar/baz"},
		{"example.com/foo", "foo"},
		{"example.com/foo/bar/baz", "foo/bar/baz"},
		{"socket:///foo/bar.sock", "foo/bar.sock"},
		{"mysql://user:password@example.com:3306/фу/бар/баз?tblsprefix=fs_", "фу/бар/баз"},
		{"mysql://user:password@example.com:3306/фу/бар/баз", "фу/бар/баз"},
		{"kafka://username:pasword@tcp(ip1:9093,ip2:9093,ip3:9093)/?topic=vsulblog", ""},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetPath() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetPath())
		}
	}
}
func TestParse_Source(t *testing.T) {
	tests := []struct {
		dsn      string
		expected string
	}{
		{"mysql://user:password@example.com:3306/foo?tblsprefix=fs_", "user:password@example.com:3306/foo?tblsprefix=fs_"},
		{"mysql://user:password@tcp(example.com:3306)/foo?tblsprefix=fs_", "user:password@tcp(example.com:3306)/foo?tblsprefix=fs_"},
		{"mysql://user:password@example.com:3306/foo/bar/baz?tblsprefix=fs_", "user:password@example.com:3306/foo/bar/baz?tblsprefix=fs_"},
		{"mysql://user:password@example.com:3306//?tblsprefix=fs_", "user:password@example.com:3306//?tblsprefix=fs_"},
		{"mysql://user:password@example.com:3306/?tblsprefix=fs_", "user:password@example.com:3306/?tblsprefix=fs_"},
		{"mysql://user:password@example.com:3306/foo", "user:password@example.com:3306/foo"},
		{"mysql://user:password@example.com:3306/foo/bar/baz", "user:password@example.com:3306/foo/bar/baz"},
		{"mysql://user:password@example.com:3306", "user:password@example.com:3306"},
		{"mysql://user:password@example.com:3306/", "user:password@example.com:3306/"},
		{"mysql://user:password@example.com:3306//", "user:password@example.com:3306//"},
		{"example.com/foo?tblsprefix=fs_", "example.com/foo?tblsprefix=fs_"},
		{"example.com/foo/bar/baz?tblsprefix=fs_", "example.com/foo/bar/baz?tblsprefix=fs_"},
		{"example.com/foo", "example.com/foo"},
		{"example.com/foo/bar/baz", "example.com/foo/bar/baz"},
		{"socket:///foo/bar.sock", "/foo/bar.sock"},
		{"mysql://user:password@example.com:3306/фу/бар/баз?tblsprefix=fs_", "user:password@example.com:3306/фу/бар/баз?tblsprefix=fs_"},
		{"mysql://user:password@example.com:3306/фу/бар/баз", "user:password@example.com:3306/фу/бар/баз"},
		{"kafka://username:pasword@tcp(ip1:9093,ip2:9093,ip3:9093)/?topic=vsulblog", "username:pasword@tcp(ip1:9093,ip2:9093,ip3:9093)/?topic=vsulblog"},
	}

	for i, test := range tests {
		dsn := Parse(test.dsn)
		if dsn.GetSource() != test.expected {
			t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", i+1, test.expected, dsn.GetSource())
		}
	}
}

func TestParse_Params(t *testing.T) {
	type ExpectedItem struct {
		key   string
		value string
	}
	type ExpectedList []ExpectedItem

	tests := []struct {
		dsn      string
		expected ExpectedList
	}{
		{"mysql://user:password@example.com:3306/dbname?fooKey=foo val&barKey=bar val", ExpectedList{
			ExpectedItem{"fooKey", "foo val"},
			ExpectedItem{"barKey", "bar val"},
		}},
		{"mysql://user:password@example.com:3306/dbname?fooKey&barKey=&&", ExpectedList{
			ExpectedItem{"fooKey", ""},
			ExpectedItem{"barKey", ""},
		}},
		{"mysql://user:password@example.com:3306/dbname?foo\\&Key=foo\\&val&bar\\&Key=bar\\&val", ExpectedList{
			ExpectedItem{"foo&Key", "foo&val"},
			ExpectedItem{"bar&Key", "bar&val"},
		}},
		{"mysql://user:password@example.com:3306/dbname?foo\\=Key=foo\\=val&bar\\=Key=bar\\=val", ExpectedList{
			ExpectedItem{"foo=Key", "foo=val"},
			ExpectedItem{"bar=Key", "bar=val"},
		}},
		{"mysql://user:password@example.com:3306/dbname?foo\\=Key\\&=\\&foo\\=val&bar\\=Key\\&=\\&bar\\=val", ExpectedList{
			ExpectedItem{"foo=Key&", "&foo=val"},
			ExpectedItem{"bar=Key&", "&bar=val"},
		}},
		{"mysql://user:password@example.com:3306/?fooKey=foo val&barKey=bar val", ExpectedList{
			ExpectedItem{"fooKey", "foo val"},
			ExpectedItem{"barKey", "bar val"},
		}},
		{"mysql://user:password@example.com", ExpectedList{}},
		{"kafka://username:pasword@tcp(ip1:9093,ip2:9093,ip3:9093)/?topic=vsulblog", ExpectedList{
			ExpectedItem{"topic", "vsulblog"},
		}},
	}

	for testId, test := range tests {
		dsn := Parse(test.dsn)

		if len(test.expected) != len(dsn.GetParams()) {
			t.Errorf("The number of results obtained is different from the expected. Test: %v.", testId+1)
		}

		for _, expected := range test.expected {
			if dsn.GetParam(expected.key) != expected.value {
				t.Errorf("Unexpected value in test \"%v\". Expected: \"%s\". Result: \"%s\"", testId+1, expected.value, dsn.GetParam(expected.key))
			}
		}
	}
}
