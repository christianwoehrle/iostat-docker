package main

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type IO struct {
	ContainerName       string
	Rchar               int64
	Wchar               int64
	Syscr               int64
	Syscw               int64
	ReadBytes           int64
	WriteBytes          int64
	CancelledWriteBytes int64
}

func (a IO) Stringer() string {
	//fmt.Println(a.containerName + " ==> " + strconv.Itoa(int(a.deltaIO)))
	return a.ContainerName + " ==> " + strconv.Itoa(int(a.ReadBytes)) + "  " + strconv.Itoa(int(a.WriteBytes))
}

var IOIndizes = map[string]int{
	"proc.io.rchar":                 0,
	"proc.io.wchar":                 1,
	"proc.io.syscr":                 2,
	"proc.io.syscw":                 3,
	"proc.io.read_bytes":            4,
	"proc.io.write_bytes":           5,
	"proc.io.cancelled_write_bytes": 6,
}

func readIO(filename string) (string, error) {
	pStatFile, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(pStatFile), nil
}

func ReadIOArray(filename string) ([]string, error) {
	pIORaw, err := readIO(filename)

	if err != nil {
		return nil, err
	}

	var io = strings.Split(pIORaw, "\n")
	io = io[:len(io)-1]

	for i := range io {
		io[i] = io[i][strings.Index(io[i], " "):]
	}

	return io, nil
}

func ReadIO(filename string) (*IO, error) {
	pIORaw, err := readIO(filename)

	//fmt.Println(pIORaw)
	if err != nil {
		return nil, err
	}

	p, err := ParseIO(pIORaw)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func ParseIO(io string) (*IO, error) {
	var pIO = new(IO)

	var ioRaw = strings.Split(io, "\n")

	ioRaw2 := ioRaw[0 : len(ioRaw)-1]
	var r = make([]int64, len(ioRaw2))

	for i := range r {
		//fmt.Println(ioRaw[i])
		//fmt.Println(">" + ioRaw[i][strings.Index(ioRaw[i], " ")+1:] + "<")
		r[i], _ = strconv.ParseInt(ioRaw[i][strings.Index(ioRaw[i], " ")+1:], 10, 32)
		//fmt.Println(r[i])
	}

	pIO.Rchar = r[0]
	pIO.Wchar = r[1]
	pIO.Syscr = r[2]
	pIO.Syscw = r[3]
	pIO.ReadBytes = r[4]
	pIO.WriteBytes = r[5]
	pIO.CancelledWriteBytes = r[6]

	return pIO, nil
}
