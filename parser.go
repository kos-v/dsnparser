package dsnparser

// DSN container for data after dsn parsing.
type DSN struct {
	raw       string
	scheme    string
	source    string // dsn source : dsn without schema
	user      string
	password  string
	host      string
	port      string
	hostport  string
	path      string
	params    map[string]string
	transport string
}

// GetHost returns a host as the string.
func (d *DSN) GetHost() string {
	return d.host
}

// GetHostPort returns a hostport as the string.
func (d *DSN) GetHostPort() string {
	return d.hostport
}

// GetParam returns an additional parameter by key as the string.
func (d *DSN) GetParam(key string) string {
	if !d.HasParam(key) {
		return ""
	}
	return d.params[key]
}

// GetParams returns additional parameters as key-value map.
func (d *DSN) GetParams() map[string]string {
	return d.params
}

// HasParam checks for the existence of an additional parameter.
func (d *DSN) HasParam(key string) bool {
	if _, ok := d.params[key]; ok {
		return true
	}
	return false
}

// GetPassword returns a credential password as the string.
func (d *DSN) GetPassword() string {
	return d.password
}

// GetPath returns a path as the string.
func (d *DSN) GetPath() string {
	return d.path
}

// GetPort returns a port as the string.
func (d *DSN) GetPort() string {
	return d.port
}

// GetRaw returns the dsn in its raw form, as it was passed to the Parse function.
func (d *DSN) GetRaw() string {
	return d.raw
}

// GetScheme returns a scheme as the string.
func (d *DSN) GetScheme() string {
	return d.scheme
}

// GetTransport returns a transport as the string.
func (d *DSN) GetTransport() string {
	return d.transport
}

// GetUser returns a credential user as the string.
func (d *DSN) GetUser() string {
	return d.user
}

// GetSource returns a dsn source  user as the string.
func (d *DSN) GetSource() string {
	return d.source
}

// Parse receives a raw dsn as argument, parses it and returns it in the DSN structure.
func Parse(raw string) *DSN {
	d := DSN{
		raw:    raw,
		source: raw,
		params: map[string]string{},
	}
	dsn := []rune(d.raw)

	// Parsing the scheme
	for pos, symbol := range dsn {
		// Found end of the scheme name
		if symbol == '/' && pos > 2 && string(dsn[pos-2:pos+1]) == "://" {
			d.scheme = string(dsn[0 : pos-2])
			dsn = dsn[pos+1:]
			d.source = string(dsn)
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
					d.user = string(unEscape([]rune{':', '@'}, credentials[0:credPos]))
					d.password = string(unEscape([]rune{':', '@'}, credentials[credPos+1:]))
					break
				}
			}
			if !hasSeparator {
				d.user = string(unEscape([]rune{':', '@'}, credentials))
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

		d.transport = string(dsn[:hpExtractBeginPos-1])
		dsn = append(dsn[hpExtractBeginPos:hpExtractEndPos+1], dsn[hpExtractEndPos+2:]...)
		break
	}

	// multihost := false
	// Host and port parsing
	for dsnPos, dsnSymbol := range dsn {
		endPos := -1
		if dsnSymbol == '/' {
			endPos = dsnPos
		} else if dsnPos == len(dsn)-1 {
			endPos = len(dsn)
		}
		if endPos == -1 {
			continue
		}
		hostPort := dsn[0:endPos]

		hasSeparator := false
		for hpPos, hpSymbol := range hostPort {
			if hpSymbol == ':' {
				hasSeparator = true
				d.host = string(hostPort[0:hpPos])
				d.port = string(hostPort[hpPos+1:])
				break
			}
		}
		if !hasSeparator {
			d.host = string(hostPort)
		}
		d.hostport = string(hostPort)

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
			d.path = string(dsn[0:endPos])
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
				d.params[string(unEscape([]rune{'=', '&'}, paramKey))] = string(unEscape([]rune{'=', '&'}, paramVal))
			}
		}
	}

	return &d
}
