package dsnparser

type dsn struct {
	raw       string
	scheme    string
	user      string
	password  string
	host      string
	port      string
	path      string
	params    map[string]string
	transport string
}

func (d *dsn) GetHost() string {
	return d.host
}

func (d *dsn) GetPassword() string {
	return d.password
}

func (d *dsn) GetParam(key string) string {
	if !d.HasParam(key) {
		return ""
	}
	return d.params[key]
}

func (d *dsn) GetParams() map[string]string {
	return d.params
}

func (d *dsn) GetPath() string {
	return d.path
}

func (d *dsn) GetPort() string {
	return d.port
}

func (d *dsn) GetRaw() string {
	return d.raw
}

func (d *dsn) GetScheme() string {
	return d.scheme
}

func (d *dsn) GetTransport() string {
	return d.transport
}

func (d *dsn) GetUser() string {
	return d.user
}

func (d *dsn) HasParam(key string) bool {
	if _, ok := d.params[key]; ok {
		return true
	}
	return false
}

func Parse(raw string) *dsn {
	d := dsn{
		raw:    raw,
		params: map[string]string{},
	}
	dsn := []rune(d.raw)

	// Parsing the scheme
	for pos, symbol := range dsn {
		// Found end of the scheme name
		if symbol == '/' && pos > 2 && string(dsn[pos-2:pos+1]) == "://" {
			d.scheme = string(dsn[0 : pos-2])
			dsn = dsn[pos+1:]
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

	// Host and port parsing
	for dsnPos, dsnSymbol := range dsn {
		endPos := -1
		if dsnSymbol == '/' {
			endPos = dsnPos
		} else if dsnPos == len(dsn)-1 {
			endPos = len(dsn)
		}

		if endPos > -1 {
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

			dsn = dsn[dsnPos+1:]
			break
		}
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
