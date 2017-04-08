package host

func (d driver) FlagsForRun(svcName string) ([]string, error) {
	return []string{`--network=host`, `-e=`}, nil
}
