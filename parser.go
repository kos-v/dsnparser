package dsnparser

// DSN container for data after dsn parsing.
type DSN struct {
	Raw    string
	Scheme string
	// dsn source : dsn without schema
	Source    string
	User      string
	Password  string
	Host      string
	Port      string
	Path      string
	Params    map[string]string
	Transport string
}

// GetHost returns a host as the string.
func (d *DSN) GetHost() string {
	return d.Host
}

// GetParam returns an additional parameter by key as the string.
func (d *DSN) GetParam(key string) string {
	if !d.HasParam(key) {
		return ""
	}
	return d.Params[key]
}

// GetParams returns additional parameters as key-value map.
func (d *DSN) GetParams() map[string]string {
	return d.Params
}

// HasParam checks for the existence of an additional parameter.
func (d *DSN) HasParam(key string) bool {
	if _, ok := d.Params[key]; ok {
		return true
	}
	return false
}

// GetPassword returns a credential password as the string.
func (d *DSN) GetPassword() string {
	return d.Password
}

// GetPath returns a path as the string.
func (d *DSN) GetPath() string {
	return d.Path
}

// GetPort returns a port as the string.
func (d *DSN) GetPort() string {
	return d.Port
}

// GetRaw returns the dsn in its raw form, as it was passed to the Parse function.
func (d *DSN) GetRaw() string {
	return d.Raw
}

// GetScheme returns a scheme as the string.
func (d *DSN) GetScheme() string {
	return d.Scheme
}

// GetTransport returns a transport as the string.
func (d *DSN) GetTransport() string {
	return d.Transport
}

// GetUser returns a credential user as the string.
func (d *DSN) GetUser() string {
	return d.User
}

// GetSource returns a dsn source  user as the string.
func (d *DSN) GetSource() string {
	return d.User
}

// Parse receives a raw dsn as argument, parses it and returns it in the DSN structure.
func Parse(raw string) *DSN {
	d := DSN{
		Raw:    raw,
		Source: raw,
		Params: map[string]string{},
	}
	dsn := []rune(d.Raw)

	// Parsing the scheme
	for pos, symbol := range dsn {
		// Found end of the scheme name
		if symbol == '/' && pos > 2 && string(dsn[pos-2:pos+1]) == "://" {
			d.Scheme = string(dsn[0 : pos-2])
			dsn = dsn[pos+1:]
			d.Source = string(dsn)
			break
		}
	}

	// Parsing the credentials
	for dsnPos, dsnSymbol := range dsn {
		// Found end of the credentials
		if dsnSymbol == '@' && !isEscaped(dsnPos, dsn) {
			credentials := dsn[0:dsnPos]

			// Separating username and password
			hasSeparator := false
			for credPos, credChar := range credentials {
				if credChar == ':' && !isEscaped(credPos, credentials) {
					hasSeparator = true
					d.User = string(unEscape([]rune{':', '@'}, credentials[0:credPos]))
					d.Password = string(unEscape([]rune{':', '@'}, credentials[credPos+1:]))
					break
				}
			}
			if !hasSeparator {
				d.User = string(unEscape([]rune{':', '@'}, credentials))
			}

			dsn = dsn[dsnPos+1:]
			break
		}
	}

	// Transport parsing
	for dsnPos, dsnSymbol := range dsn {
		if dsnSymbol != '(' {
			continue
		}

		hpExtractBeginPos := dsnPos + 1
		hpExtractEndPos := -1
		for hpPos, hpSymbol := range dsn[hpExtractBeginPos:] {
			if hpSymbol == ')' {
				hpExtractEndPos = dsnPos + hpPos
			}
		}
		if hpExtractEndPos == -1 {
			continue
		}

		d.Transport = string(dsn[:hpExtractBeginPos-1])
		dsn = append(dsn[hpExtractBeginPos:hpExtractEndPos+1], dsn[hpExtractEndPos+2:]...)
		break
	}

	multihost := false
	// Host and port parsing
	for dsnPos, dsnSymbol := range dsn {
		endPos := -1
		if dsnSymbol == '/' {
			endPos = dsnPos
		} else if dsnSymbol == ',' {
			multihost = true
		} else if dsnPos == len(dsn)-1 {
			endPos = len(dsn)
		}
		if endPos == -1 {
			continue
		}
		hostPort := dsn[0:endPos]

		if multihost {
			d.Host = string(hostPort)
		} else {

			hasSeparator := false
			for hpPos, hpSymbol := range hostPort {
				if hpSymbol == ':' {
					hasSeparator = true
					d.Host = string(hostPort[0:hpPos])
					d.Port = string(hostPort[hpPos+1:])
					break
				}
			}
			if !hasSeparator {
				d.Host = string(hostPort)
			}
		}

		dsn = dsn[dsnPos+1:]
		break

	}

	// Path parsing
	for pos, symbol := range dsn {
		endPos := -1
		if symbol == '?' {
			endPos = pos
		} else if pos == len(dsn)-1 {
			endPos = len(dsn)
		}

		if endPos > -1 {
			d.Path = string(dsn[0:endPos])
			dsn = dsn[pos+1:]
			break
		}
	}

	// Params parsing
	beginPosParam := 0
	for symbolPos, symbol := range dsn {
		param := []rune{}
		if symbol == '&' && !isEscaped(symbolPos, dsn) {
			param = dsn[beginPosParam:symbolPos]
			beginPosParam = symbolPos + 1
		} else if symbolPos == len(dsn)-1 {
			param = dsn[beginPosParam:]
		}

		// Separating key and value
		if len(param) > 0 {
			paramKey := []rune{}
			paramVal := []rune{}

			hasSeparator := false
			for paramSymbolPos, paramSymbol := range param {
				if paramSymbol == '=' && !isEscaped(paramSymbolPos, param) {
					hasSeparator = true
					paramKey = param[0:paramSymbolPos]
					paramVal = param[paramSymbolPos+1:]
					break
				}
			}
			if !hasSeparator {
				paramKey = param
			}

			if len(paramKey) > 0 {
				d.Params[string(unEscape([]rune{'=', '&'}, paramKey))] = string(unEscape([]rune{'=', '&'}, paramVal))
			}
		}
	}

	return &d
}
