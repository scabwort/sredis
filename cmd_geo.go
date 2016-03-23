package sredis

const (
	CMD_geoadd            = "GEOADD"
	CMD_geodist           = "GEODIST"
	CMD_geohash           = "GEOHASH"
	CMD_geopos            = "GEOPOS"
	CMD_georadius         = "GEORADIUS"
	CMD_georadiusbymember = "GEORADIUSBYMEMBER"
)

var (
	Arg_WITHDIST  = "WITHDIST"
	Arg_WITHCOORD = "WITHCOORD"
)

func (cmd *Commond) GeoAdd(key string, latitude, longitude float64, name string) *Result {
	cmd.Cmd = CMD_geoadd
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_geoadd, 5)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteFloat64(&latitude)
	cmd.buf.WriteFloat64(&longitude)
	cmd.buf.WriteString(&name)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) GeoDist(key, name1, name2 string) *Result {
	cmd.dostr3(CMD_geodist, key, name1, name2)
	return &cmd.result
}

func (cmd *Commond) GeoHash(key string, names ...string) *Result {
	cmd.Cmd = CMD_geohash
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_geohash, len(names)+1)
	for idx, _ := range names {
		cmd.buf.WriteString(&names[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) GeoPos(key string, names ...string) *Result {
	cmd.Cmd = CMD_geopos
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_geopos, len(names)+1)
	for idx, _ := range names {
		cmd.buf.WriteString(&names[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) GeoRadius(key string, latitude, longitude float64, unit string) *Result {
	cmd.Cmd = CMD_georadius
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_georadius, 7)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteFloat64(&latitude)
	cmd.buf.WriteFloat64(&longitude)
	cmd.buf.WriteString(&unit)
	cmd.buf.WriteString(&Arg_WITHDIST)
	cmd.buf.WriteString(&Arg_WITHCOORD)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) GeoRadiusWithDist(key string, latitude, longitude float64, unit string) *Result {
	cmd.Cmd = CMD_georadius
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_georadius, 6)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteFloat64(&latitude)
	cmd.buf.WriteFloat64(&longitude)
	cmd.buf.WriteString(&unit)
	cmd.buf.WriteString(&Arg_WITHDIST)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) GeoRadiusWithCoord(key string, latitude, longitude float64, unit string) *Result {
	cmd.Cmd = CMD_georadius
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_georadius, 6)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteFloat64(&latitude)
	cmd.buf.WriteFloat64(&longitude)
	cmd.buf.WriteString(&unit)
	cmd.buf.WriteString(&Arg_WITHCOORD)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) GeoRadiusByMember(key string, name string, dist float64, unit string) *Result {
	cmd.Cmd = CMD_georadiusbymember
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_georadiusbymember, 5)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteString(&name)
	cmd.buf.WriteFloat64(&dist)
	cmd.buf.WriteString(&unit)
	cmd.waitConn()
	return &cmd.result
}
