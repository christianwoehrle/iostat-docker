package main

import (
	"context"
	"flag"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DeltaIO struct {
	deltaIO       int64
	containerName string
}

func (a DeltaIO) Stringer() string {
	//fmt.Println(a.containerName + " ==> " + strconv.Itoa(int(a.deltaIO)))
	return a.containerName + ";" + strconv.Itoa(int(a.deltaIO))
}

type DeltaIOs []DeltaIO

func (a DeltaIOs) Len() int           { return len(a) }
func (a DeltaIOs) Less(i, j int) bool { return a[i].deltaIO < a[j].deltaIO }
func (a DeltaIOs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (a DeltaIOs) Stringer() string {
	s := ""
	for i := len(a) - 1; i > len(a)-10 && i >= 0; i-- {
		s = s + a[i].Stringer() + "\n"
	}

	return s
}

func main() {

	duration := flag.Int("d", 10, "duration in seconds")

	flag.Parse()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	currentContainers := make(map[int]*IO, 100)

	for {
		oldContainers := make(map[int]*IO, 100)
		for i, v := range currentContainers {
			oldContainers[i] = v
		}

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			panic(err)
		}

		for _, container := range containers {
			//fmt.Printf("%s %s %s\n", container.ID[:10], container.Image, container.Names[0])
			containerJson, err := cli.ContainerInspect(context.Background(), container.ID[:10])
			if err == nil {
				//fmt.Printf("Pid %d\n", containerJson.State.Pid)

				io, err := ReadIO("/proc/" + strconv.Itoa(containerJson.State.Pid) + "/io")
				if err == nil {
					io.ContainerName = containerJson.Name
					if err != nil {
						fmt.Println(err.Error())
					} else {
						currentContainers[containerJson.State.Pid] = io
						//fmt.Printf("IO: %v\n", io)
					}
				} else {
					fmt.Println("Error reading File " + "/proc/" + strconv.Itoa(containerJson.State.Pid) + "/io: " + err.Error())
				}

			} else {
				fmt.Println("err: " + err.Error())
			}

		}

		ios := make(DeltaIOs, 0)
		for i, oldIoContainer := range oldContainers {
			//fmt.Println(i, oldIoContainer)

			newIOContainer, ok := currentContainers[i]
			if ok {
				diffW := newIOContainer.WriteBytes - oldIoContainer.WriteBytes
				diffR := newIOContainer.ReadBytes - oldIoContainer.ReadBytes
				//diff := newIOContainer.Wchar - oldIoContainer.Wchar
				//fmt.Printf("Pid: %d %d %d \n", i, newIOContainer.WriteBytes, newIOContainer.ReadBytes)
				//fmt.Printf("Pid: %d %d %d \n", i, oldIoContainer.WriteBytes, newIOContainer.ReadBytes)
				//fmt.Println("Pid: " + strconv.Itoa(i) + " " + oldIoContainer.ContainerName + " ==> " + strconv.Itoa(int(diffW+diffR+diff)))
				deltaIo := DeltaIO{
					containerName: newIOContainer.ContainerName,
					deltaIO:       diffW + diffR,
				}
				ios = append(ios, deltaIo)
			}
		}
		sort.Sort(DeltaIOs(ios))

		fmt.Println(time.Now().Format(time.RFC1123))
		fmt.Println(ios.Stringer())

		time.Sleep(time.Duration(*duration) * time.Second)
	}

}
